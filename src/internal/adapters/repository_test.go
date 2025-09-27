package adapters

import (
	"avito-test/internal/domain"
	"avito-test/pkg/utils"
	"context"
	"fmt"
	"testing"

	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv4/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/jackc/pgx/v4/pgxpool"
)

func Test(t *testing.T) {
	ctx := context.Background()

	uri := fmt.Sprintf("postgres://%s:%s@%s:%d/%s",
		"postgres", "password", "localhost", 5432, "shop",
	)
	pool, err := pgxpool.Connect(ctx, uri)
	utils.CheckErr(err)
	defer pool.Close()

	item_repo := NewItemRepository(pool, trmpgx.DefaultCtxGetter)

	_, err = item_repo.GetByName(ctx, "test")
	require.Equal(t, err, domain.ErrEntityDoesNotExist)

	item, err := item_repo.GetByName(ctx, "pink-hoody")
	require.NoError(t, err)
	require.Equal(t, item.Price, 500)

	passHash := []byte("passhash")
	var expected_user = domain.User{
		Username:     domain.UserName(uuid.NewString()),
		PasswordHash: passHash,
		Balance:      114,
	}

	user_repo := NewUserRepository(pool, trmpgx.DefaultCtxGetter)

	err = user_repo.Save(ctx, &expected_user)
	require.NoError(t, err)

	user, err := user_repo.GetByUsername(ctx, expected_user.Username)
	require.NoError(t, err)
	require.Equal(t, expected_user, *user)

	transaction_repo := NewTransactionRepository(pool, trmpgx.DefaultCtxGetter)

	expected_transaction := domain.Transaction{
		From:   expected_user.Username,
		To:     expected_user.Username,
		Amount: 500,
	}

	err = transaction_repo.Save(ctx, &expected_transaction)
	require.NoError(t, err)

	transaction, err := transaction_repo.GetById(ctx, expected_transaction.Id)
	require.NoError(t, err)
	require.Equal(t, expected_transaction, *transaction)

	history, err := transaction_repo.GetAllByUsername(ctx, user.Username)
	require.NoError(t, err)
	require.True(t, len(history) == 1)
	require.Equal(t, expected_transaction, *(history[0]))

	expected_second_transaction := domain.Transaction{
		From:   expected_user.Username,
		To:     expected_user.Username,
		Amount: 600,
	}

	err = transaction_repo.Save(ctx, &expected_second_transaction)
	require.NoError(t, err)

	history, err = transaction_repo.GetAllByUsername(ctx, user.Username)
	require.NoError(t, err)
	require.True(t, len(history) == 2)
	require.Equal(t, expected_transaction, *(history[0]))
	require.Equal(t, expected_second_transaction, *(history[1]))

	possession_repo := NewPosessionRepository(pool, trmpgx.DefaultCtxGetter)

	inventory, err := possession_repo.GetUserInventory(ctx, expected_user.Username)
	require.NoError(t, err)
	require.True(t, len(inventory) == 0)

	first_possession := domain.Possession{
		Username: expected_user.Username,
		Item:     "pen",
		Amount:   10,
	}

	err = possession_repo.Save(ctx, &first_possession)
	require.NoError(t, err)

	inventory, err = possession_repo.GetUserInventory(ctx, expected_user.Username)
	require.NoError(t, err)
	require.True(t, len(inventory) == 1)
	require.Equal(t, first_possession.Amount, inventory[first_possession.Item])

	possession, err := possession_repo.GetPossessionByUsernameAndItemName(ctx, expected_user.Username, first_possession.Item)
	require.NoError(t, err)
	require.Equal(t, first_possession, *possession)
}
