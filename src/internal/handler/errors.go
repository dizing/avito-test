package handler

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SetInvalidRequestError(c *gin.Context, err error) {
	c.JSON(http.StatusBadRequest, ErrorResponse{
		Errors: err.Error(),
	})
	c.Abort()
	log.Println("Invalid Request ", err)
}

func SetUnauthorizedError(c *gin.Context, err error) {
	c.JSON(http.StatusUnauthorized, ErrorResponse{
		Errors: err.Error(),
	})
	c.Abort()
	log.Println("Anauthorized error ", err)
}

func SetInternalError(c *gin.Context, err error) {
	c.JSON(http.StatusInternalServerError, ErrorResponse{
		Errors: err.Error(),
	})
	c.Abort()
	log.Println("Internal error ", err)
}
