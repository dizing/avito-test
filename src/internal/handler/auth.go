package handler

import (
	"avito-test/internal/domain"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type authHandler struct {
	userRepo domain.UserRepository
}

func NewAuthHandler(userRepo domain.UserRepository) *authHandler {
	return &authHandler{
		userRepo: userRepo,
	}
}

func (h *authHandler) Register(r gin.IRoutes) *authHandler {
	r.POST("/api/auth", h.login)

	return h
}

func (h *authHandler) login(c *gin.Context) {
	var request AuthRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		SetInvalidRequestError(c, err)
		return
	}

	user, err := h.userRepo.GetByUsername(c, domain.UserName(request.Username))
	if err != nil {
		if errors.Is(err, domain.ErrEntityDoesNotExist) {
			passwordHash, _ := bcrypt.GenerateFromPassword([]byte(request.Password), 12)
			user = domain.NewUser(domain.UserName(request.Username), passwordHash)

			h.userRepo.Save(c, user)
		} else {
			SetInternalError(c, err)
			return
		}
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(request.Password)); err != nil {
		SetInvalidRequestError(c, fmt.Errorf("password or username is invalid"))
		return
	}

	c.JSON(http.StatusOK, AuthResponse{
		Token: createJwt(user),
	})
}
