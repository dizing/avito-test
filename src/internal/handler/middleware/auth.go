package middleware

import (
	"avito-test/internal/domain"
	"avito-test/internal/handler"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
)

func NewAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			handler.SetUnauthorizedError(c, fmt.Errorf("'Authorization' header required"))
			return
		}

		const BearerPrefix = "Bearer "

		for strings.HasPrefix(tokenString, BearerPrefix) {
			tokenString = strings.TrimPrefix(tokenString, BearerPrefix)
		}

		username, err := domain.ExtractClaimsFromJwt(tokenString)
		if err != nil {
			handler.SetUnauthorizedError(c, err)
			return
		}

		c.Set(domain.CONTEXT_CREDENTIALS, username)

		c.Next()
	}
}
