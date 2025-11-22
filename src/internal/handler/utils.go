package handler

import (
	"avito-test/internal/domain"

	"github.com/gin-gonic/gin"
	"github.com/samber/lo"
)

func GetAuthorizedUserName(c *gin.Context) domain.UserName {
	credentials, exists := c.Get(domain.CONTEXT_CREDENTIALS)
	lo.Assert(exists, "must use GetAuthorizedUserName only in authorized handler")

	return credentials.(domain.UserName)
}
