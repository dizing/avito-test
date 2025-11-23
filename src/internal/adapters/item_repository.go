package adapters

import (
	"avito-test/internal/domain"
	"context"
	"sync"

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
var itemCache = map[domain.ItemName]*domain.Item{}
var itemCacheMtx = sync.RWMutex{}

func getItemFromCache(name domain.ItemName) *domain.Item {
	itemCacheMtx.RLock()
	defer itemCacheMtx.RUnlock()

	if i, ok := itemCache[name]; ok {
		item := *i
		return &item
	}

	return nil
}

func cacheItem(name domain.ItemName, item domain.Item) {
	itemCacheMtx.Lock()
	defer itemCacheMtx.Unlock()

	itemCache[name] = &item
}

func (r *ItemRepository) GetByName(ctx context.Context, name domain.ItemName) (*domain.Item, error) {
	if item := getItemFromCache(name); item != nil {
		return item, nil
	}

	query := `SELECT * FROM items WHERE name=$1`

	row := r.getter.DefaultTrOrDB(ctx, r.db).QueryRow(ctx, query, name)

	var item domain.Item

	err := row.Scan(&item.Name, &item.Price)
	if err != nil {
		return nil, mapPgxError(err)
	}

	cacheItem(name, item)

	return &item, nil
}
