package main

import (
	"avito-test/internal/adapters"
	"avito-test/internal/handler"
	"avito-test/pkg/utils"
	"context"
	"fmt"
	"net/http"

	"github.com/jackc/pgx/v4/pgxpool"

	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv4/v2"
	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"

	"github.com/gin-gonic/gin"
)

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func main() {
	ctx := context.Background()
	r := gin.Default()

	uri := fmt.Sprintf("postgres://%s:%s@%s:%d/%s",
		"postgres", "password", "postgres", 5432, "shop",
	)
	pool, err := pgxpool.Connect(ctx, uri)
	utils.CheckErr(err)
	defer pool.Close()

	ctxGetter := trmpgx.DefaultCtxGetter

	user_repo := adapters.NewUserRepository(pool, ctxGetter)
	item_repo := adapters.NewItemRepository(pool, ctxGetter)
	possession_repo := adapters.NewPosessionRepository(pool, ctxGetter)
	transactions_repo := adapters.NewTransactionRepository(pool, ctxGetter)
	user_info_repo := adapters.NewInfoRepository(user_repo, possession_repo, transactions_repo)

	r.Use(CORSMiddleware())

	trManager := manager.Must(trmpgx.NewDefaultFactory(pool))

	handler.NewAuthHandler(trManager, user_repo).Register(r)

	authorize_group := r.Group("/")
	authorize_group.Use(handler.NewAuthMiddleware())

	handler.NewBuyHandler(trManager, item_repo, user_repo, possession_repo).Register(authorize_group)
	handler.NewInfoHandler(trManager, user_info_repo).Register(authorize_group)
	handler.NewSendHandler(trManager, user_repo, transactions_repo).Register(authorize_group)

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "UP",
		})
	})

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
