package main

import (
	"avito-test/internal/adapters"
	"avito-test/internal/handler"
	"avito-test/internal/handler/middleware"
	"avito-test/pkg/utils"
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/jackc/pgx/v4/pgxpool"

	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv4/v2"
	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"

	"github.com/gin-gonic/gin"
)

func main() {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, nil)))

	ctx := context.Background()
	r := gin.Default()

	uri := fmt.Sprintf("postgres://%s:%s@%s:%d/%s",
		"postgres", "password", "postgres", 5432, "shop",
	)
	config, err := pgxpool.ParseConfig(uri)
	utils.CheckErr(err)

	config.MaxConns = 300
	pool, err := pgxpool.ConnectConfig(ctx, config)

	utils.CheckErr(err)
	defer pool.Close()

	ctxGetter := trmpgx.DefaultCtxGetter

	userRepo := adapters.NewUserRepository(pool, ctxGetter)
	itemRepo := adapters.NewItemRepository(pool, ctxGetter)
	possession_repo := adapters.NewPosessionRepository(pool, ctxGetter)
	transactions_repo := adapters.NewTransactionRepository(pool, ctxGetter)
	userInfoRepo := adapters.NewInfoRepository(userRepo, possession_repo, transactions_repo)

	r.Use(middleware.NewCORSMiddleware())

	trManager := manager.Must(trmpgx.NewDefaultFactory(pool))

	handler.NewAuthHandler(trManager, userRepo).Register(r)

	authorizeGroup := r.Group("/")
	authorizeGroup.Use(middleware.NewAuthMiddleware())

	handler.NewBuyHandler(trManager, itemRepo, userRepo, possession_repo).Register(authorizeGroup)
	handler.NewInfoHandler(trManager, userInfoRepo).Register(authorizeGroup)
	handler.NewSendHandler(trManager, userRepo, transactions_repo).Register(authorizeGroup)

	handler.NewHealthHandler().Register(r)

	r.Run("0.0.0.0:8080") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
