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
	// 绑定模式专用字段
	BindMode bool   // 是否是绑定模式
	UserID   uint   // 绑定模式下的用户ID
	Token    string // 绑定模式下的用户Token（用于前端恢复登录状态）
}

// oauthStateStore 简单的内存状态存储
var oauthStateStore = struct {
	sync.RWMutex
	states map[string]OAuthState
}{states: make(map[string]OAuthState)}

// AuthService 认证服务
type AuthService struct {
	userRepo     *repository.UserRepository
	jwtManager   *jwt.JWTManager
	oauthConfig  *config.LinuxDoOAuthConfig
	emailService *EmailService
}

// NewAuthService 创建认证服务
func NewAuthService(userRepo *repository.UserRepository, cfg *config.Config, emailService *EmailService) *AuthService {
	return &AuthService{
		userRepo:     userRepo,
		jwtManager:   jwt.NewJWTManager(cfg.JWT.Secret, cfg.JWT.ExpireHours),
		oauthConfig:  &cfg.OAuth.LinuxDo,
		emailService: emailService,
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
	token, err := s.jwtManager.GenerateToken(user.ID, user.Email, user.Username, int(user.Role), user.TrustLevel, user.LinuxDoID)
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
	token, err := s.jwtManager.GenerateToken(user.ID, user.Email, user.Username, int(user.Role), user.TrustLevel, user.LinuxDoID)
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
	params.Set("scope", "user")
	params.Set("state", state)

	authURL := fmt.Sprintf("%s?%s", linuxDoAuthorizeURL, params.Encode())

	return authURL, state, nil
}

// OAuthCallbackResult OAuth回调结果
type OAuthCallbackResult struct {
	LoginResponse *dto.LoginResponse // 登录模式返回
	BindUser      *models.User       // 绑定模式返回绑定后的用户
	IsBindMode    bool               // 是否是绑定模式
	UserToken     string             // 绑定模式下原用户的token
}

// HandleOAuthCallback 处理OAuth回调（支持登录和绑定两种模式）
func (s *AuthService) HandleOAuthCallback(code, state string) (*OAuthCallbackResult, error) {
	// 获取并验证state
	stateData, valid := getAndValidateState(state)
	if !valid {
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

	// 根据模式处理
	if stateData.BindMode {
		// 绑定模式：将LinuxDO账号绑定到现有用户
		user, err := s.bindLinuxDoToUser(stateData.UserID, userInfo)
		if err != nil {
			return nil, err
		}
		return &OAuthCallbackResult{
			BindUser:   user,
			IsBindMode: true,
			UserToken:  stateData.Token,
		}, nil
	}

	// 登录模式：登录或注册
	loginResp, err := s.oauthLoginOrRegister(userInfo)
	if err != nil {
		return nil, err
	}
	return &OAuthCallbackResult{
		LoginResponse: loginResp,
		IsBindMode:    false,
	}, nil
}

// bindLinuxDoToUser 内部方法：将LinuxDO账号绑定到现有用户
func (s *AuthService) bindLinuxDoToUser(userID uint, linuxDoInfo *LinuxDoUserInfo) (*models.User, error) {
	linuxDoID := fmt.Sprintf("%d", linuxDoInfo.ID)

	// 检查该LinuxDo账号是否已被其他用户绑定
	existingUser, err := s.userRepo.FindByLinuxDoID(linuxDoID)
	if err == nil && existingUser != nil && existingUser.ID != userID {
		return nil, errors.New("该LinuxDo账号已被其他用户绑定")
	}

	// 获取当前用户
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, errors.New("用户不存在")
	}

	// 更新用户的LinuxDo信息
	user.LinuxDoID = linuxDoID
	user.LinuxDoUsername = linuxDoInfo.Username
	user.AvatarURL = linuxDoInfo.AvatarURL
	user.TrustLevel = linuxDoInfo.TrustLevel

	// 如果信任等级>=2，提升为认证用户
	if linuxDoInfo.TrustLevel >= 2 && user.Role == models.RoleNormal {
		user.Role = models.RoleCertified
	}

	if err := s.userRepo.Update(user); err != nil {
		return nil, errors.New("绑定失败")
	}

	return user, nil
}

// GetBindLinuxDoURL 获取绑定LinuxDO的OAuth URL (绑定模式，使用标准回调地址)
func (s *AuthService) GetBindLinuxDoURL(userID uint, userToken string) (string, string, error) {
	if s.oauthConfig.ClientID == "" {
		return "", "", errors.New("OAuth未配置")
	}

	// 生成随机state防止CSRF攻击
	state, err := generateRandomState()
	if err != nil {
		return "", "", errors.New("生成state失败")
	}

	// 存储state(5分钟过期)，包含绑定模式信息
	oauthStateStore.Lock()
	oauthStateStore.states[state] = OAuthState{
		State:     state,
		ExpiresAt: time.Now().Add(5 * time.Minute),
		BindMode:  true,
		UserID:    userID,
		Token:     userToken,
	}
	oauthStateStore.Unlock()

	// 清理过期的state
	go cleanExpiredStates()

	// 使用配置的标准回调地址
	callbackURL := s.oauthConfig.RedirectURI

	// 构建授权URL
	params := url.Values{}
	params.Set("client_id", s.oauthConfig.ClientID)
	params.Set("redirect_uri", callbackURL)
	params.Set("response_type", "code")
	params.Set("scope", "user")
	params.Set("state", state)

	authURL := fmt.Sprintf("%s?%s", linuxDoAuthorizeURL, params.Encode())

	return authURL, state, nil
}

// BindEmail 绑定邮箱（适用于LinuxDO登录用户）
func (s *AuthService) BindEmail(userID uint, req *dto.BindEmailRequest) (*models.User, error) {
	// 获取当前用户
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, errors.New("用户不存在")
	}

	// 检查用户是否已经有真实邮箱（不是占位邮箱）
	if user.Password != "" {
		return nil, errors.New("您已绑定邮箱，无法重复绑定")
	}

	// 检查邮箱是否已被其他用户使用
	existingUser, err := s.userRepo.FindByEmail(req.Email)
	if err == nil && existingUser != nil && existingUser.ID != userID {
		return nil, errors.New("该邮箱已被其他用户使用")
	}

	// 密码加密
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("密码加密失败")
	}

	// 更新用户邮箱和密码
	user.Email = req.Email
	user.Password = string(hashedPassword)

	if err := s.userRepo.Update(user); err != nil {
		return nil, errors.New("绑定邮箱失败")
	}

	return user, nil
}

// UnbindLinuxDo 解绑LinuxDO账号
func (s *AuthService) UnbindLinuxDo(userID uint) (*models.User, error) {
	// 获取当前用户
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, errors.New("用户不存在")
	}

	// 检查用户是否有密码（邮箱注册用户）
	if user.Password == "" {
		return nil, errors.New("OAuth登录用户不能解绑LinuxDo账号")
	}

	// 清除LinuxDo绑定信息
	user.LinuxDoID = ""
	user.LinuxDoUsername = ""
	user.TrustLevel = 0

	// 如果用户是因为LinuxDo而获得的认证用户权限，降级为普通用户
	// 注意：管理员不会被降级
	if user.Role == models.RoleCertified {
		user.Role = models.RoleNormal
	}

	if err := s.userRepo.Update(user); err != nil {
		return nil, errors.New("解绑失败")
	}

	return user, nil
}

// UpdateProfile 更新用户资料
func (s *AuthService) UpdateProfile(userID uint, req *dto.UpdateProfileRequest) (*models.User, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, errors.New("用户不存在")
	}

	if req.Username != "" {
		user.Username = req.Username
	}

	if err := s.userRepo.Update(user); err != nil {
		return nil, errors.New("更新资料失败")
	}

	return user, nil
}

// SendEmailVerificationCode 发送邮箱验证码
func (s *AuthService) SendEmailVerificationCode(userID uint, newEmail string) error {
	// 检查用户是否存在
	_, err := s.userRepo.FindByID(userID)
	if err != nil {
		return errors.New("用户不存在")
	}

	// 检查新邮箱是否已被其他用户使用
	existingUser, err := s.userRepo.FindByEmail(newEmail)
	if err == nil && existingUser != nil && existingUser.ID != userID {
		return errors.New("该邮箱已被其他用户使用")
	}

	// 发送验证码
	if s.emailService == nil {
		return errors.New("邮件服务未配置")
	}

	_, err = s.emailService.SendEmailVerificationCode(newEmail, userID)
	if err != nil {
		return errors.New("发送验证码失败")
	}

	return nil
}

// ChangeEmail 修改邮箱
func (s *AuthService) ChangeEmail(userID uint, req *dto.ChangeEmailRequest) (*models.User, error) {
	// 检查用户是否存在
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, errors.New("用户不存在")
	}

	// 检查新邮箱是否已被其他用户使用
	existingUser, err := s.userRepo.FindByEmail(req.NewEmail)
	if err == nil && existingUser != nil && existingUser.ID != userID {
		return nil, errors.New("该邮箱已被其他用户使用")
	}

	// 验证验证码
	if s.emailService == nil {
		return nil, errors.New("邮件服务未配置")
	}

	if !s.emailService.VerifyEmailCode(userID, req.NewEmail, req.Code) {
		return nil, errors.New("验证码无效或已过期")
	}

	// 更新邮箱
	user.Email = req.NewEmail

	if err := s.userRepo.Update(user); err != nil {
		return nil, errors.New("修改邮箱失败")
	}

	return user, nil
}

// UpdateAvatar 更新头像
func (s *AuthService) UpdateAvatar(userID uint, req *dto.UpdateAvatarRequest) (*models.User, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, errors.New("用户不存在")
	}

	user.AvatarURL = req.AvatarURL

	if err := s.userRepo.Update(user); err != nil {
		return nil, errors.New("更新头像失败")
	}

	return user, nil
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
	return s.exchangeCodeForTokenWithRedirectURI(code, s.oauthConfig.RedirectURI)
}

// exchangeCodeForTokenWithRedirectURI 用授权码换取access_token（自定义redirect_uri）
func (s *AuthService) exchangeCodeForTokenWithRedirectURI(code, redirectURI string) (string, error) {
	if redirectURI == "" {
		redirectURI = s.oauthConfig.RedirectURI
	}

	data := url.Values{}
	data.Set("grant_type", "authorization_code")
	data.Set("code", code)
	data.Set("redirect_uri", redirectURI)
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

// validateState 验证state（简单版本，仅返回是否有效）
func validateState(state string) bool {
	_, valid := getAndValidateState(state)
	return valid
}

// getAndValidateState 获取并验证state，返回state数据和是否有效
func getAndValidateState(state string) (OAuthState, bool) {
	oauthStateStore.Lock()
	defer oauthStateStore.Unlock()

	storedState, exists := oauthStateStore.states[state]
	if !exists {
		return OAuthState{}, false
	}

	// 删除已使用的state
	delete(oauthStateStore.states, state)

	// 检查是否过期
	if time.Now().After(storedState.ExpiresAt) {
		return OAuthState{}, false
	}

	return storedState, true
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
