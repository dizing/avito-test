package handler

import (
	"avito-test/internal/domain"
	"fmt"
	"net/http"

	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"
	"github.com/gin-gonic/gin"
)

type infoHandler struct {
	trManager    *manager.Manager
	userInfoRepo domain.UserInfoRepository
}

func NewInfoHandler(trManager *manager.Manager, userInfoRepo domain.UserInfoRepository) *infoHandler {
	return &infoHandler{trManager, userInfoRepo}
}

func (h *infoHandler) Register(r gin.IRoutes) *infoHandler {
	r.GET("/api/info", h.GetInfo)

	return h
}

func (h *infoHandler) GetInfo(c *gin.Context) {
	username := GetAuthorizedUserName(c)

	userInfo, err := h.userInfoRepo.GetByUsername(c, username)
	if err != nil {
		SetInternalError(c, fmt.Errorf("can't find user from valid jwt token: %w", err))
		return
	}

	c.JSON(http.StatusOK, InfoResponse{
		Coins:       userInfo.User.Balance,
		Inventory:   mapInventory(userInfo.Inventory),
		CoinHistory: mapCoinHistory(userInfo.Transactions, userInfo.User.Username),
	})
}

func mapInventory(model domain.UserInventory) []Inventory {
	dto := make([]Inventory, 0, len(model))

	for itemName, count := range model {
		dto = append(dto, Inventory{
			Type:     string(itemName),
			Quantity: count,
		})
	}

	return dto
}

func mapCoinHistory(model domain.TransactionHistory, currentUserName domain.UserName) CoinHistory {
	var dto CoinHistory

	for _, transaction := range model {
		if transaction.To == currentUserName {
			dto.Received = append(dto.Received, ReceivedTransaction{
				FromUser: string(transaction.From),
				Amount:   transaction.Amount,
			})
		}

		if transaction.From == currentUserName {
			dto.Sent = append(dto.Sent, SentTransaction{
				ToUser: string(transaction.To),
				Amount: transaction.Amount,
			})
		}
	}

	return dto
}
