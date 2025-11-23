package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type healthHandler struct {
}

func NewHealthHandler() *healthHandler {
	return &healthHandler{}
}

func (h *healthHandler) Register(r gin.IRoutes) *healthHandler {
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "UP",
		})
	})

	return h
}
