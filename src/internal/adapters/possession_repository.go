package adapters

import (
	"avito-test/internal/domain"
	"context"
	"errors"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/samber/lo"

	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv4/v2"
)

type possessionRepository struct {
	db     *pgxpool.Pool
	getter *trmpgx.CtxGetter
}

func NewPosessionRepository(db *pgxpool.Pool, c *trmpgx.CtxGetter) domain.PossessionRepository {
	return &possessionRepository{db, c}
}

func (r *possessionRepository) GetPossessionByUsernameAndItemName(ctx context.Context, username domain.UserName, itemName domain.ItemName) (*domain.Possession, error) {
	query := `SELECT amount FROM user_items WHERE username=$1 AND item_name=$2`

	row := r.getter.DefaultTrOrDB(ctx, r.db).QueryRow(ctx, query, username, itemName)

	var amount uint

	err := row.Scan(&amount)
	if err != nil {
		mapped_error := mapPgxError(err)

		if !errors.Is(mapped_error, domain.ErrEntityDoesNotExist) {
			return nil, mapped_error
		}

		amount = 0
	}

	return lo.ToPtr(domain.Possession{Username: username, Item: itemName, Amount: amount}), nil
}

func (r *possessionRepository) GetUserInventory(ctx context.Context, username domain.UserName) (domain.UserInventory, error) {
	query := `SELECT item_name, amount FROM user_items WHERE username=$1`

	rows, err := r.getter.DefaultTrOrDB(ctx, r.db).Query(ctx, query, username)
	if err != nil {
		return nil, mapPgxError(err)
	}
	defer rows.Close()

	var possessions []*domain.Possession

	for rows.Next() {
		// TODO: this is a good place for separating domain from infrastructure
		var (
			item   domain.ItemName
			amount uint
		)

		err := rows.Scan(
			&item,
			&amount)
		if err != nil {
			return nil, mapPgxError(err)
		}

		possessions = append(possessions, lo.ToPtr(domain.Possession{Username: username, Item: item, Amount: amount}))
	}

	if err := rows.Err(); err != nil {
		return nil, mapPgxError(err)
	}

	return *domain.NewUserInventory(possessions), nil
}

func (r *possessionRepository) Save(ctx context.Context, possession *domain.Possession) error {
	if possession.Amount == 0 {
		query := `
        DELETE FROM user_items WHERE username = $1;
    `

		if _, err := r.getter.DefaultTrOrDB(ctx, r.db).Exec(ctx, query, possession.Username); err != nil {
			mapped_error := mapPgxError(err)

			if errors.Is(mapped_error, domain.ErrEntityDoesNotExist) {
				return nil
			}

			return mapped_error
		}

		return nil
	}

	query := `
        INSERT INTO user_items (username, item_name, amount) 
        VALUES ($1, $2, $3)
        ON CONFLICT (username, item_name) 
        DO UPDATE SET
				amount = EXCLUDED.amount
    `

	if _, err := r.getter.DefaultTrOrDB(ctx, r.db).Exec(ctx, query, possession.Username, possession.Item, possession.Amount); err != nil {
		return mapPgxError(err)
	}

	return nil
}
