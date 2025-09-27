package handler

import (
	"avito-test/internal/domain"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type buyHandler struct {
	itemRepo      domain.ItemRepository
	userRepo      domain.UserRepository
	posessionRepo domain.PossessionRepository
}

func NewBuyHandler(
	itemRepo domain.ItemRepository,
	userRepo domain.UserRepository,
	posessionRepo domain.PossessionRepository) *buyHandler {
	return &buyHandler{
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

	credentials := GetAuthorizedUser(c)

	user, err := h.userRepo.GetByUsername(c, credentials.Username)
	if err != nil {
		SetInternalError(c, fmt.Errorf("can't find user from valid jwt token: %w", err))
		return
	}

	if user.Balance < item.Price {
		SetInvalidRequestError(c, fmt.Errorf("not enough money"))
		return
	}

	possession, err := h.posessionRepo.GetPossessionByUsernameAndItemName(c, user.Username, item.Name)
	if err != nil {
		if errors.Is(err, domain.ErrEntityDoesNotExist) {
			possession = domain.NewEmptyPossession(user.Username, item.Name)
		} else {
			SetInternalError(c, err)
			return
		}
	}

	user.Balance -= item.Price
	possession.Amount += 1

	err = h.userRepo.Save(c, user)
	if err != nil {
		SetInternalError(c, err)
	}

	err = h.posessionRepo.Save(c, possession)
	if err != nil {
		SetInternalError(c, err)
	}

	c.Status(http.StatusOK)
}
