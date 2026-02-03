package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"gomall/internal/config"
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
	secretKey   []byte
	expireHours int
}

// NewJWT 创建JWT实例
func NewJWT() *JWT {
	jwtConfig := config.GetJWT()
	return &JWT{
		secretKey:   []byte(jwtConfig.GetString("secret")),
		expireHours: jwtConfig.GetInt("expire_hours"),
	}
}

// GenerateToken 生成Token
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
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.secretKey)
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
