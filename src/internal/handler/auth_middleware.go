package handler

import (
	"avito-test/internal/domain"
	"fmt"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func NewAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			SetUnauthorizedError(c, fmt.Errorf("'Authorization' header required"))
			return
		}
		tokenString = strings.TrimPrefix(tokenString, "Bearer ")

		user, err := extractClaims(tokenString)
		if err != nil {
			SetUnauthorizedError(c, err)
			return
		}
		c.Set("credentials", user)

		c.Next()
	}
}

func GetAuthorizedUser(c *gin.Context) *domain.User {
	credentials, exists := c.Get("credentials")
	if !exists {
		return nil
	}

	return credentials.(*domain.User)
}

const (
	CLAIMS_USERNAME = "Username"
	CLAIMS_BALANCE  = "Balance"
	CLAIMS_EXP      = "exp"
)

func createJwt(user *domain.User) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		CLAIMS_USERNAME: user.Username,
		CLAIMS_BALANCE:  user.Balance,
		CLAIMS_EXP:      time.Now().Add(time.Hour * 72).Unix(),
	})
	tokenString, _ := token.SignedString([]byte("secret"))

	return tokenString
}

func extractClaims(tokenString string) (*domain.User, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte("secret"), nil
	})
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("authorization token invalid")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("can't parse token claims")
	}

	return &domain.User{
		Username: domain.UserName(claims[CLAIMS_USERNAME].(string)),
		Balance:  int(claims[CLAIMS_BALANCE].(float64)),
	}, nil
}
