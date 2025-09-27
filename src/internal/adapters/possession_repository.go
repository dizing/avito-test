package adapters

import (
	"avito-test/internal/domain"
	"context"

	"github.com/jackc/pgx/v4/pgxpool"

	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv4/v2"
)

type possessionRepository struct {
	db     *pgxpool.Pool
	getter *trmpgx.CtxGetter
}

func NewPosessionRepository(db *pgxpool.Pool, c *trmpgx.CtxGetter) domain.PossessionRepository {
	return &possessionRepository{db, c}
}

func (r *possessionRepository) GetPossessionByUsernameAndItemName(ctx context.Context, username domain.UserName, item_name domain.ItemName) (*domain.Possession, error) {
	query := `SELECT amount FROM User_Items WHERE username=$1 AND item_name=$2`

	row := r.getter.DefaultTrOrDB(ctx, r.db).QueryRow(ctx, query, username, item_name)

	posession := domain.NewEmptyPossession(username, item_name)

	err := row.Scan(&posession.Amount)
	if err != nil {
		return nil, mapPgxError(err)
	}

	return posession, nil
}

func (r *possessionRepository) GetUserInventory(ctx context.Context, username domain.UserName) (domain.UserInventory, error) {
	query := `SELECT item_name, amount FROM User_Items WHERE user_id=$1`

	rows, err := r.getter.DefaultTrOrDB(ctx, r.db).Query(ctx, query, username)
	if err != nil {
		return nil, mapPgxError(err)
	}

	var possessions []*domain.Possession

	for rows.Next() {
		// TODO: this is a good place for separating domain from infrastructure
		possession := domain.NewEmptyPossession(username, domain.ItemName(""))
		err := rows.Scan(
			&possession.Item,
			&possession.Amount)
		if err != nil {
			return nil, mapPgxError(err)
		}

		possessions = append(possessions, possession)
	}

	if err := rows.Err(); err != nil {
		return nil, mapPgxError(err)
	}

	return *domain.NewUserInventory(possessions), nil
}

func (r *possessionRepository) Save(ctx context.Context, possession *domain.Possession) error {
	// TODO: if amount == 0 delete row

	query := `
        INSERT INTO User_Items (username, item_name, amount) 
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
