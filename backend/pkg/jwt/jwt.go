package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"gomall/backend/internal/config"
)

var (
	// ErrTokenExpired token已过期
	ErrTokenExpired = errors.New("token已过期")
	// ErrInvalidToken token无效
	ErrInvalidToken = errors.New("token无效")
)

// Claims JWT载荷
type Claims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	jwt.RegisteredClaims
}

// JWT JWT工具类
type JWT struct {
	secretKey       []byte
	expireHours     int
	refreshHours    int  // refresh token 过期时间（小时）
}

// NewJWT 创建JWT实例
func NewJWT() *JWT {
	jwtConfig := config.GetJWT()
	expireHours := jwtConfig.GetInt("expire_hours")
	refreshHours := jwtConfig.GetInt("refresh_hours")

	// 设置默认值
	if expireHours <= 0 {
		expireHours = 24 // 默认24小时
	}
	if refreshHours <= 0 {
		refreshHours = 168 // 默认7天
	}

	return &JWT{
		secretKey:    []byte(jwtConfig.GetString("secret")),
		expireHours:  expireHours,
		refreshHours: refreshHours,
	}
}

// GenerateToken 生成访问Token（短有效期）
func (j *JWT) GenerateToken(userID uint, username, email string) (string, error) {
	nowTime := time.Now()
	expireTime := nowTime.Add(time.Duration(j.expireHours) * time.Hour)

	claims := Claims{
		UserID:   userID,
		Username: username,
		Email:    email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expireTime),
			IssuedAt:  jwt.NewNumericDate(nowTime),
			Issuer:    "gomall",
			Subject:   "access",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.secretKey)
}

// GenerateRefreshToken 生成刷新Token（长有效期）
func (j *JWT) GenerateRefreshToken(userID uint, username, email string) (string, error) {
	nowTime := time.Now()
	expireTime := nowTime.Add(time.Duration(j.refreshHours) * time.Hour)

	claims := Claims{
		UserID:   userID,
		Username: username,
		Email:    email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expireTime),
			IssuedAt:  jwt.NewNumericDate(nowTime),
			Issuer:    "gomall",
			Subject:   "refresh",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.secretKey)
}

// GenerateTokenPair 生成Token对（access + refresh）
type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn   int    `json:"expires_in"` // 过期时间（秒）
}

func (j *JWT) GenerateTokenPair(userID uint, username, email string) (*TokenPair, error) {
	accessToken, err := j.GenerateToken(userID, username, email)
	if err != nil {
		return nil, err
	}

	refreshToken, err := j.GenerateRefreshToken(userID, username, email)
	if err != nil {
		return nil, err
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:   j.expireHours * 3600,
	}, nil
}

// ParseToken 解析Token
func (j *JWT) ParseToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return j.secretKey, nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrTokenExpired
		}
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, ErrInvalidToken
	}

	return claims, nil
}

// RefreshToken 刷新Token
func (j *JWT) RefreshToken(tokenString string) (string, error) {
	claims, err := j.ParseToken(tokenString)
	if err != nil {
		return "", err
	}

	return j.GenerateToken(claims.UserID, claims.Username, claims.Email)
}
