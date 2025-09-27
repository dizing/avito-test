package handler

import (
	"avito-test/internal/domain"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type sendHandler struct {
	userRepo         domain.UserRepository
	transactionsRepo domain.TransactionRepository
}

func NewSendHandler(userRepo domain.UserRepository, transactionsRepo domain.TransactionRepository) *sendHandler {
	return &sendHandler{
		userRepo:         userRepo,
		transactionsRepo: transactionsRepo,
	}
}

func (h *sendHandler) Register(r gin.IRoutes) *sendHandler {
	r.POST("/api/sendCoin", h.SendCoins)

	return h
}

func (h *sendHandler) SendCoins(c *gin.Context) {
	var request SendCoinRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		SetInvalidRequestError(c, err)
		return
	}

	credentials := GetAuthorizedUser(c)

	from_user, err := h.userRepo.GetByUsername(c, credentials.Username)
	if err != nil {
		SetInternalError(c, fmt.Errorf("can't find user from valid jwt token: %w", err))
		return
	}

	if from_user.Balance < request.Amount {
		SetInvalidRequestError(c, fmt.Errorf("not enough money to send"))
		return
	}

	to_user, err := h.userRepo.GetByUsername(c, domain.UserName(request.ToUser))
	if err != nil {
		SetInvalidRequestError(c, fmt.Errorf("invalid destination user: %w", err))
		return
	}

	transaction := domain.NewTransaction(from_user.Username, to_user.Username, request.Amount)
	from_user.Balance -= request.Amount
	to_user.Balance += request.Amount

	err = h.transactionsRepo.Save(c, transaction)
	if err != nil {
		SetInternalError(c, err)
	}
	err = h.userRepo.Save(c, from_user)
	if err != nil {
		SetInternalError(c, err)
	}
	err = h.userRepo.Save(c, to_user)
	if err != nil {
		SetInternalError(c, err)
	}

	c.Status(http.StatusOK)
}
