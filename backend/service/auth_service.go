package service

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"linuxdo-review/config"
	"linuxdo-review/dto"
	"linuxdo-review/models"
	"linuxdo-review/pkg/jwt"
	"linuxdo-review/repository"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

const (
	// Linux.do OAuth endpoints
	linuxDoAuthorizeURL = "https://connect.linux.do/oauth2/authorize"
	linuxDoTokenURL     = "https://connect.linux.do/oauth2/token"
	linuxDoUserInfoURL  = "https://connect.linux.do/api/user"
)

// OAuthState OAuth状态管理(防止CSRF攻击)
type OAuthState struct {
	State     string
	ExpiresAt time.Time
}

// oauthStateStore 简单的内存状态存储
var oauthStateStore = struct {
	sync.RWMutex
	states map[string]OAuthState
}{states: make(map[string]OAuthState)}

// AuthService 认证服务
type AuthService struct {
	userRepo    *repository.UserRepository
	jwtManager  *jwt.JWTManager
	oauthConfig *config.LinuxDoOAuthConfig
}

// NewAuthService 创建认证服务
func NewAuthService(userRepo *repository.UserRepository, cfg *config.Config) *AuthService {
	return &AuthService{
		userRepo:    userRepo,
		jwtManager:  jwt.NewJWTManager(cfg.JWT.Secret, cfg.JWT.ExpireHours),
		oauthConfig: &cfg.OAuth.LinuxDo,
	}
}

// Register 用户注册
func (s *AuthService) Register(req *dto.RegisterRequest) (*models.User, error) {
	// 检查邮箱是否已存在
	existing, err := s.userRepo.FindByEmail(req.Email)
	if err == nil && existing != nil {
		return nil, errors.New("邮箱已被注册")
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	// 密码加密
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("密码加密失败")
	}

	// 创建用户
	user := &models.User{
		Email:    req.Email,
		Password: string(hashedPassword),
		Username: req.Username,
		Role:     models.RoleNormal,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, errors.New("创建用户失败")
	}

	return user, nil
}

// Login 用户登录
func (s *AuthService) Login(req *dto.LoginRequest) (*dto.LoginResponse, error) {
	// 查找用户
	user, err := s.userRepo.FindByEmail(req.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("用户不存在")
		}
		return nil, err
	}

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, errors.New("密码错误")
	}

	// 生成JWT Token
	token, err := s.jwtManager.GenerateToken(user.ID, user.Email, user.Username, int(user.Role))
	if err != nil {
		return nil, errors.New("生成Token失败")
	}

	return &dto.LoginResponse{
		Token: token,
		User:  dto.ToUserResponse(user),
	}, nil
}

// GetUserByID 根据ID获取用户
func (s *AuthService) GetUserByID(id uint) (*models.User, error) {
	return s.userRepo.FindByID(id)
}

// IsSystemInitialized 检查系统是否已初始化（是否存在管理员）
func (s *AuthService) IsSystemInitialized() (bool, error) {
	return s.userRepo.HasAdmin()
}

// SetupAdmin 初始化管理员账号
func (s *AuthService) SetupAdmin(req *dto.SetupAdminRequest) (*models.User, error) {
	// 检查是否已有管理员
	hasAdmin, err := s.userRepo.HasAdmin()
	if err != nil {
		return nil, errors.New("检查系统状态失败")
	}
	if hasAdmin {
		return nil, errors.New("系统已初始化，不能重复创建管理员")
	}

	// 检查邮箱是否已存在
	existing, err := s.userRepo.FindByEmail(req.Email)
	if err == nil && existing != nil {
		return nil, errors.New("邮箱已被注册")
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	// 密码加密
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("密码加密失败")
	}

	// 创建管理员用户
	user := &models.User{
		Email:    req.Email,
		Password: string(hashedPassword),
		Username: req.Username,
		Role:     models.RoleAdmin,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, errors.New("创建管理员失败")
	}

	return user, nil
}

// OAuthLoginOrRegister OAuth登录或注册(内部方法)
func (s *AuthService) oauthLoginOrRegister(userInfo *LinuxDoUserInfo) (*dto.LoginResponse, error) {
	linuxDoID := fmt.Sprintf("%d", userInfo.ID)

	// 先尝试通过LinuxDo ID查找用户
	user, err := s.userRepo.FindByLinuxDoID(linuxDoID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	if user == nil {
		// 用户不存在,创建新用户
		// 生成一个唯一的邮箱(如果用户没有提供)
		email := userInfo.Email
		if email == "" {
			email = fmt.Sprintf("%s@linuxdo.user", userInfo.Username)
		}

		// 根据信任等级确定用户角色
		role := models.RoleNormal
		if userInfo.TrustLevel >= 2 {
			role = models.RoleCertified // 信任等级2及以上为认证用户
		}

		user = &models.User{
			Email:           email,
			Username:        userInfo.Username,
			LinuxDoID:       linuxDoID,
			LinuxDoUsername: userInfo.Username,
			AvatarURL:       userInfo.AvatarURL,
			TrustLevel:      userInfo.TrustLevel,
			Role:            role,
		}

		if err := s.userRepo.Create(user); err != nil {
			return nil, errors.New("创建用户失败")
		}
	} else {
		// 更新用户信息
		user.LinuxDoUsername = userInfo.Username
		user.AvatarURL = userInfo.AvatarURL
		user.TrustLevel = userInfo.TrustLevel

		// 如果信任等级提升,更新角色
		if userInfo.TrustLevel >= 2 && user.Role == models.RoleNormal {
			user.Role = models.RoleCertified
		}

		if err := s.userRepo.Update(user); err != nil {
			return nil, errors.New("更新用户信息失败")
		}
	}

	// 生成JWT Token
	token, err := s.jwtManager.GenerateToken(user.ID, user.Email, user.Username, int(user.Role))
	if err != nil {
		return nil, errors.New("生成Token失败")
	}

	return &dto.LoginResponse{
		Token: token,
		User:  dto.ToUserResponse(user),
	}, nil
}

// GetOAuthURL 获取OAuth授权URL
func (s *AuthService) GetOAuthURL() (string, string, error) {
	if s.oauthConfig.ClientID == "" {
		return "", "", errors.New("OAuth未配置")
	}

	// 生成随机state防止CSRF攻击
	state, err := generateRandomState()
	if err != nil {
		return "", "", errors.New("生成state失败")
	}

	// 存储state(5分钟过期)
	oauthStateStore.Lock()
	oauthStateStore.states[state] = OAuthState{
		State:     state,
		ExpiresAt: time.Now().Add(5 * time.Minute),
	}
	oauthStateStore.Unlock()

	// 清理过期的state
	go cleanExpiredStates()

	// 构建授权URL
	params := url.Values{}
	params.Set("client_id", s.oauthConfig.ClientID)
	params.Set("redirect_uri", s.oauthConfig.RedirectURI)
	params.Set("response_type", "code")
	params.Set("state", state)

	authURL := fmt.Sprintf("%s?%s", linuxDoAuthorizeURL, params.Encode())

	return authURL, state, nil
}

// HandleOAuthCallback 处理OAuth回调
func (s *AuthService) HandleOAuthCallback(code, state string) (*dto.LoginResponse, error) {
	// 验证state
	if !validateState(state) {
		return nil, errors.New("无效的state参数")
	}

	// 用code换取access_token
	accessToken, err := s.exchangeCodeForToken(code)
	if err != nil {
		return nil, fmt.Errorf("获取access_token失败: %w", err)
	}

	// 获取用户信息
	userInfo, err := s.fetchUserInfo(accessToken)
	if err != nil {
		return nil, fmt.Errorf("获取用户信息失败: %w", err)
	}

	// 登录或注册
	return s.oauthLoginOrRegister(userInfo)
}

// LinuxDoUserInfo Linux.do用户信息
type LinuxDoUserInfo struct {
	ID         int    `json:"id"`
	Username   string `json:"username"`
	Name       string `json:"name"`
	Email      string `json:"email"`
	AvatarURL  string `json:"avatar_url"`
	TrustLevel int    `json:"trust_level"`
	Active     bool   `json:"active"`
	Silenced   bool   `json:"silenced"`
}

// LinuxDoTokenResponse Linux.do Token响应
type LinuxDoTokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
}

// exchangeCodeForToken 用授权码换取access_token
func (s *AuthService) exchangeCodeForToken(code string) (string, error) {
	data := url.Values{}
	data.Set("grant_type", "authorization_code")
	data.Set("code", code)
	data.Set("redirect_uri", s.oauthConfig.RedirectURI)
	data.Set("client_id", s.oauthConfig.ClientID)
	data.Set("client_secret", s.oauthConfig.ClientSecret)

	req, err := http.NewRequest("POST", linuxDoTokenURL, strings.NewReader(data.Encode()))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("token请求失败: %s", string(body))
	}

	var tokenResp LinuxDoTokenResponse
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return "", err
	}

	if tokenResp.AccessToken == "" {
		return "", errors.New("未获取到access_token")
	}

	return tokenResp.AccessToken, nil
}

// fetchUserInfo 获取用户信息
func (s *AuthService) fetchUserInfo(accessToken string) (*LinuxDoUserInfo, error) {
	req, err := http.NewRequest("GET", linuxDoUserInfoURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Accept", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("获取用户信息失败: %s", string(body))
	}

	var userInfo LinuxDoUserInfo
	if err := json.Unmarshal(body, &userInfo); err != nil {
		return nil, err
	}

	return &userInfo, nil
}

// generateRandomState 生成随机state
func generateRandomState() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

// validateState 验证state
func validateState(state string) bool {
	oauthStateStore.Lock()
	defer oauthStateStore.Unlock()

	storedState, exists := oauthStateStore.states[state]
	if !exists {
		return false
	}

	// 删除已使用的state
	delete(oauthStateStore.states, state)

	// 检查是否过期
	return time.Now().Before(storedState.ExpiresAt)
}

// cleanExpiredStates 清理过期的state
func cleanExpiredStates() {
	oauthStateStore.Lock()
	defer oauthStateStore.Unlock()

	now := time.Now()
	for state, data := range oauthStateStore.states {
		if now.After(data.ExpiresAt) {
			delete(oauthStateStore.states, state)
		}
	}
}
