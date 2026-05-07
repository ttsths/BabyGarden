package middleware

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"strings"
	"time"
	"yuanzi-backend/config"
	"yuanzi-backend/model"
	"yuanzi-backend/pkg/gredis"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func JWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		code := model.SUCCESS
		token := c.GetHeader("Authorization")
		if token == "" || !strings.HasPrefix(token, "Bearer ") {
			code = model.ERROR_AUTH_CHECK_TOKEN_FAIL
		} else {
			token = strings.TrimPrefix(token, "Bearer ")
			claims, err := ParseToken(token)
			if err != nil {
				code = model.ERROR_AUTH_CHECK_TOKEN_FAIL
			} else {
				c.Set("userId", claims.UserID)
				c.Set("phone", claims.Phone)
				c.Set("claims", claims)
			}
		}
		if code != model.SUCCESS {
			c.JSON(http.StatusUnauthorized, gin.H{"code": code, "msg": model.GetMsg(code), "data": nil})
			c.Abort()
			return
		}
		c.Next()
	}
}

type Claims struct {
	UserID string `json:"user_id"`
	Phone  string `json:"phone"`
	JTI    string `json:"jti"`
	Type   string `json:"type"`
	jwt.RegisteredClaims
}

func GenerateTokenPair(userID string, phone string) (string, string, error) {
	jti, err := generateJTI()
	if err != nil {
		return "", "", err
	}
	accessClaims := Claims{UserID: userID, Phone: phone, JTI: jti, Type: "access", RegisteredClaims: jwt.RegisteredClaims{ExpiresAt: jwt.NewNumericDate(time.Now().Add(2 * time.Hour)), Issuer: "yuanzi"}}
	accessTokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessToken, err := accessTokenClaims.SignedString([]byte(config.GlobalConfig.JWT.Secret))
	if err != nil {
		return "", "", err
	}
	refreshClaims := Claims{UserID: userID, Phone: phone, JTI: jti, Type: "refresh", RegisteredClaims: jwt.RegisteredClaims{ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)), Issuer: "yuanzi"}}
	refreshTokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshToken, err := refreshTokenClaims.SignedString([]byte(config.GlobalConfig.JWT.Secret))
	if err != nil {
		return "", "", err
	}
	return accessToken, refreshToken, nil
}

func ParseToken(token string) (*Claims, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.GlobalConfig.JWT.Secret), nil
	})
	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
			return claims, nil
		}
	}
	return nil, err
}

func generateJTI() (string, error) {
	buf := make([]byte, 6)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	return fmt.Sprintf("jti_%d_%s", time.Now().UnixNano(), hex.EncodeToString(buf)), nil
}
func AddToBlacklist(jti string, expiration time.Duration) error {
	return gredis.SetEx("token:blacklist:"+jti, "1", int(expiration.Seconds()))
}
func IsBlacklisted(jti string) (bool, error) { return gredis.Exists("token:blacklist:" + jti) }
func GetUserIDOrZero(c *gin.Context) string {
	userID, exists := c.Get("userId")
	if !exists {
		return ""
	}
	if id, ok := userID.(string); ok {
		return id
	}
	return ""
}
