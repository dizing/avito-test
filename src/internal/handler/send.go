package handler

import (
	"avito-test/internal/domain"
	"context"
	"fmt"
	"net/http"

	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"
	"github.com/gin-gonic/gin"
	"github.com/samber/lo"
)

type sendHandler struct {
	trManager        *manager.Manager
	userRepo         domain.UserRepository
	transactionsRepo domain.TransactionRepository
}

func NewSendHandler(trManager *manager.Manager, userRepo domain.UserRepository, transactionsRepo domain.TransactionRepository) *sendHandler {
	return &sendHandler{
		trManager:        trManager,
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

	username := GetAuthorizedUserName(c)

	if username == domain.UserName(request.ToUser) {
		SetInvalidRequestError(c, fmt.Errorf("can't send to itself"))
		return
	}

	if h.trManager.Do(c, func(ctx context.Context) error {
		from_user, err := h.userRepo.GetByUsername(c, username)
		if err != nil {
			SetInternalError(c, fmt.Errorf("can't find user from valid jwt token: %w", err))
			return err
		}

		if from_user.Balance < int(request.Amount) {
			SetInvalidRequestError(c, fmt.Errorf("not enough money to send"))
			return err
		}

		to_user, err := h.userRepo.GetByUsername(c, domain.UserName(request.ToUser))
		if err != nil {
			SetInvalidRequestError(c, fmt.Errorf("invalid destination user: %w", err))
			return err
		}

		from_user.Balance -= int(request.Amount)
		to_user.Balance += int(request.Amount)

		err = h.transactionsRepo.Save(c, lo.ToPtr(domain.Transaction{
			Id:     domain.TransactionUUID{},
			From:   from_user.Username,
			To:     to_user.Username,
			Amount: request.Amount}))
		if err != nil {
			SetInternalError(c, err)
			return err
		}
		err = h.userRepo.Save(c, from_user)
		if err != nil {
			SetInternalError(c, err)
			return err
		}
		err = h.userRepo.Save(c, to_user)
		if err != nil {
			SetInternalError(c, err)
			return err
		}

		return nil
	}) != nil {
		return
	}

	c.Status(http.StatusOK)
}
