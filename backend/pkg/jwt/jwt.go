package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrTokenExpired     = errors.New("token已过期")
	ErrTokenInvalid     = errors.New("token无效")
	ErrTokenMalformed   = errors.New("token格式错误")
	ErrTokenNotValidYet = errors.New("token尚未生效")
)

// Claims JWT声明
type Claims struct {
	UserID     uint   `json:"user_id"`
	Email      string `json:"email"`
	Username   string `json:"username"`
	Role       int    `json:"role"`
	TrustLevel int    `json:"trust_level"`
	LinuxDoID  string `json:"linuxdo_id"`
	jwt.RegisteredClaims
}

// JWTManager JWT管理器
type JWTManager struct {
	secret      []byte
	expireHours int
}

// NewJWTManager 创建JWT管理器
func NewJWTManager(secret string, expireHours int) *JWTManager {
	return &JWTManager{
		secret:      []byte(secret),
		expireHours: expireHours,
	}
}

// GenerateToken 生成JWT Token
func (m *JWTManager) GenerateToken(userID uint, email, username string, role int, trustLevel int, linuxDoID string) (string, error) {
	claims := Claims{
		UserID:     userID,
		Email:      email,
		Username:   username,
		Role:       role,
		TrustLevel: trustLevel,
		LinuxDoID:  linuxDoID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(m.expireHours) * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "linuxdo-review",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(m.secret)
}

// ParseToken 解析JWT Token
func (m *JWTManager) ParseToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return m.secret, nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenMalformed) {
			return nil, ErrTokenMalformed
		} else if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrTokenExpired
		} else if errors.Is(err, jwt.ErrTokenNotValidYet) {
			return nil, ErrTokenNotValidYet
		}
		return nil, ErrTokenInvalid
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, ErrTokenInvalid
}
