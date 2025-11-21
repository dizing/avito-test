package handler

import (
	"avito-test/internal/domain"
	"fmt"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/samber/lo"
)

const (
	CLAIMS_USERNAME     = "Username"
	CLAIMS_EXP          = "exp"
	SECRET              = "secret"
	CONTEXT_CREDENTIALS = "credentials"
)

func NewAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			SetUnauthorizedError(c, fmt.Errorf("'Authorization' header required"))
			return
		}
		tokenString = strings.TrimPrefix(tokenString, "Bearer ")

		username, err := extractClaims(tokenString)
		if err != nil {
			SetUnauthorizedError(c, err)
			return
		}

		c.Set(CONTEXT_CREDENTIALS, username)

		c.Next()
	}
}

func GetAuthorizedUserName(c *gin.Context) domain.UserName {
	credentials, exists := c.Get(CONTEXT_CREDENTIALS)
	lo.Assert(exists, "must use GetAuthorizedUserName only in authorized handler")

	return credentials.(domain.UserName)
}

func createJwt(username domain.UserName) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		CLAIMS_USERNAME: username,
		CLAIMS_EXP:      time.Now().Add(time.Hour * 72).Unix(),
	})
	tokenString, _ := token.SignedString([]byte(SECRET))

	return tokenString
}

func extractClaims(tokenString string) (domain.UserName, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(SECRET), nil
	})
	if err != nil {
		return "", err
	}

	if !token.Valid {
		return "", fmt.Errorf("authorization token invalid")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", fmt.Errorf("can't parse token claims")
	}

	username, ok := claims[CLAIMS_USERNAME]
	if !ok {
		return "", fmt.Errorf("claims must contain")
	}

	return domain.UserName(username.(string)), nil
}
