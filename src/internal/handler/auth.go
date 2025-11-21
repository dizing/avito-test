package handler

import (
	"avito-test/internal/domain"
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type authHandler struct {
	trManager *manager.Manager
	userRepo  domain.UserRepository
}

func NewAuthHandler(trManager *manager.Manager, userRepo domain.UserRepository) *authHandler {
	return &authHandler{
		trManager: trManager,
		userRepo:  userRepo,
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

	var username domain.UserName

	if h.trManager.Do(c, func(ctx context.Context) error {
		user, err := h.userRepo.GetByUsername(c, domain.UserName(request.Username))
		if err != nil {
			if errors.Is(err, domain.ErrEntityDoesNotExist) {
				return h.trManager.Do(c, func(ctx context.Context) error {
					passwordHash, _ := bcrypt.GenerateFromPassword([]byte(request.Password), 12)
					user = domain.NewUser(domain.UserName(request.Username), passwordHash)

					h.userRepo.Save(c, user)
					username = user.Username
					return nil
				})
			}

			SetInternalError(c, err)
			return err
		}

		if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(request.Password)); err != nil {
			SetInvalidRequestError(c, fmt.Errorf("password or username is invalid"))
			return err
		}

		username = user.Username
		return nil
	}) != nil {
		return
	}

	c.JSON(http.StatusOK, AuthResponse{
		Token: createJwt(username),
	})
}
