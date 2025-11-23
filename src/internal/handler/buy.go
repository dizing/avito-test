package handler

import (
	"avito-test/internal/domain"
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"
	"github.com/gin-gonic/gin"
)

type buyHandler struct {
	trManager     *manager.Manager
	itemRepo      domain.ItemRepository
	userRepo      domain.UserRepository
	posessionRepo domain.PossessionRepository
}

func NewBuyHandler(
	trManager *manager.Manager,
	itemRepo domain.ItemRepository,
	userRepo domain.UserRepository,
	posessionRepo domain.PossessionRepository) *buyHandler {
	return &buyHandler{
		trManager:     trManager,
		itemRepo:      itemRepo,
		userRepo:      userRepo,
		posessionRepo: posessionRepo,
	}
}

func (h *buyHandler) Register(r gin.IRoutes) *buyHandler {
	r.GET("/api/buy/:item", h.BuyItem)

	return h
}

func (h *buyHandler) BuyItem(c *gin.Context) {
	item_name := domain.ItemName(c.Param("item"))

	item, err := h.itemRepo.GetByName(c, item_name)
	if err != nil {
		if errors.Is(err, domain.ErrEntityDoesNotExist) {
			SetInvalidRequestError(c, err)
		} else {
			SetInternalError(c, err)
		}
		return
	}

	username := GetAuthorizedUserName(c)

	if h.trManager.Do(c, func(ctx context.Context) error {
		user, err := h.userRepo.GetByUsername(c, username)
		if err != nil {
			SetInternalError(c, fmt.Errorf("can't find user from valid jwt token: %w", err))
			return err
		}

		if user.Balance < int(item.Price) {
			SetInvalidRequestError(c, fmt.Errorf("not enough money"))
			return err
		}

		possession, err := h.posessionRepo.GetPossessionByUsernameAndItemName(c, user.Username, item.Name)
		if err != nil {
			SetInternalError(c, err)
			return err
		}

		user.Balance -= int(item.Price)
		possession.Amount += 1

		err = h.userRepo.Save(c, user)
		if err != nil {
			SetInternalError(c, err)
		}

		err = h.posessionRepo.Save(c, possession)
		if err != nil {
			SetInternalError(c, err)
		}

		return nil
	}) != nil {
		return
	}

	c.Status(http.StatusOK)
}
