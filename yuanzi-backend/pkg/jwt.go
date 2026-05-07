package pkg

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

type JWTClaims struct {
	UserID   string `json:"user_id"`
	Phone    string `json:"phone"`
	Nickname string `json:"nickname"`
	IsAdmin  bool   `json:"is_admin"`
	jwt.RegisteredClaims
}

type JWTManager struct {
	secretKey     []byte
	tokenDuration time.Duration
}

func NewJWTManager(secretKey string, tokenDuration time.Duration) *JWTManager {
	return &JWTManager{secretKey: []byte(secretKey), tokenDuration: tokenDuration}
}

func (m *JWTManager) GenerateToken(userID string, phone, nickname string) (string, error) {
	claims := JWTClaims{UserID: userID, Phone: phone, Nickname: nickname, RegisteredClaims: jwt.RegisteredClaims{ExpiresAt: jwt.NewNumericDate(time.Now().Add(m.tokenDuration)), IssuedAt: jwt.NewNumericDate(time.Now()), NotBefore: jwt.NewNumericDate(time.Now())}}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(m.secretKey)
}

func (m *JWTManager) VerifyToken(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrTokenInvalid
		}
		return m.secretKey, nil
	})
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrTokenExpired
		}
		if errors.Is(err, jwt.ErrTokenMalformed) {
			return nil, ErrTokenMalformed
		}
		if errors.Is(err, jwt.ErrTokenNotValidYet) {
			return nil, ErrTokenNotValidYet
		}
		return nil, ErrTokenInvalid
	}
	claims, ok := token.Claims.(*JWTClaims)
	if !ok || !token.Valid {
		return nil, ErrTokenInvalid
	}
	return claims, nil
}

func (m *JWTManager) RefreshToken(tokenString string) (string, error) {
	claims, err := m.VerifyToken(tokenString)
	if err != nil {
		return "", err
	}
	return m.GenerateToken(claims.UserID, claims.Phone, claims.Nickname)
}
