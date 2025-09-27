package handler

import (
	"avito-test/internal/domain"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type infoHandler struct {
	userInfoRepo domain.UserInfoRepository
}

func NewInfoHandler(userInfoRepo domain.UserInfoRepository) *infoHandler {
	return &infoHandler{userInfoRepo}
}

func (h *infoHandler) Register(r gin.IRoutes) *infoHandler {
	r.GET("/api/info", h.GetInfo)

	return h
}

func (h *infoHandler) GetInfo(c *gin.Context) {
	credentials := GetAuthorizedUser(c)

	userInfo, err := h.userInfoRepo.GetByUsername(c, credentials.Username)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		log.Println("incorrect request" + err.Error())
		return
		// TODO user doesn't exists
	}

	c.JSON(http.StatusOK, InfoResponse{
		Coins:       userInfo.User.Balance,
		Inventory:   mapInventory(userInfo.Inventory),
		CoinHistory: mapCoinHistory(userInfo.Transactions, userInfo.User.Username),
	})
}

func mapInventory(model domain.UserInventory) []Inventory {
	var dto []Inventory

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
