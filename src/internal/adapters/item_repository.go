package adapters

import (
	"avito-test/internal/domain"
	"context"

	"github.com/jackc/pgx/v4/pgxpool"

	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv4/v2"
)

type ItemRepository struct {
	db     *pgxpool.Pool
	getter *trmpgx.CtxGetter
}

func NewItemRepository(db *pgxpool.Pool, c *trmpgx.CtxGetter) domain.ItemRepository {
	return &ItemRepository{db, c}
}

// NOTE: items is constant right now. So item is never expire.
var item_cache = map[domain.ItemName]*domain.Item{}

func (r *ItemRepository) GetByName(ctx context.Context, name domain.ItemName) (*domain.Item, error) {
	if i, ok := item_cache[name]; ok {
		return i, nil
	}

	query := `SELECT * FROM Items WHERE name=$1`

	row := r.getter.DefaultTrOrDB(ctx, r.db).QueryRow(ctx, query, name)

	item := &domain.Item{}

	err := row.Scan(&item.Name, &item.Price)
	if err != nil {
		return nil, mapPgxError(err)
	}

	item_cache[name] = item

	return item, nil
}
