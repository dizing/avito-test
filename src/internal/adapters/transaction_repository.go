// type TransactionRepository interface {
// 	GetById(context.Context, TransactionUUID) (*Transaction, error)
// 	GetAllByUserUUID(context.Context, TransactionUUID) (TransactionHistory, error)
// 	Save(context.Context, *Transaction) error
// }

package adapters

import (
	"avito-test/internal/domain"
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"

	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv4/v2"
)

type transactionRepository struct {
	db     *pgxpool.Pool
	getter *trmpgx.CtxGetter
}

func NewTransactionRepository(db *pgxpool.Pool, c *trmpgx.CtxGetter) domain.TransactionRepository {
	return &transactionRepository{db, c}
}

func (r *transactionRepository) GetById(ctx context.Context, id domain.TransactionUUID) (*domain.Transaction, error) {
	query := `SELECT username_from, username_to, amount FROM transactions WHERE id=$1`

	row := r.getter.DefaultTrOrDB(ctx, r.db).QueryRow(ctx, query, id)

	transaction := &domain.Transaction{Id: id}

	err := row.Scan(&transaction.From, &transaction.To, &transaction.Amount)
	if err != nil {
		return nil, mapPgxError(err)
	}

	return transaction, nil
}

func (r *transactionRepository) GetTransactionHistoryByUsername(ctx context.Context, username domain.UserName) (domain.TransactionHistory, error) {
	query := `SELECT id, username_from, username_to, amount FROM transactions WHERE username_from=$1 OR username_to=$1`

	rows, err := r.getter.DefaultTrOrDB(ctx, r.db).Query(ctx, query, username)
	if err != nil {
		return nil, mapPgxError(err)
	}
	defer rows.Close()

	var transactions []*domain.Transaction

	for rows.Next() {
		transaction := domain.Transaction{}
		err := rows.Scan(
			&transaction.Id,
			&transaction.From,
			&transaction.To,
			&transaction.Amount)
		if err != nil {
			return nil, mapPgxError(err)
		}

		transactions = append(transactions, &transaction)
	}

	if err := rows.Err(); err != nil {
		return nil, mapPgxError(err)
	}

	return transactions, nil
}

func (r *transactionRepository) Save(ctx context.Context, transaction *domain.Transaction) error {
	if transaction.Id == domain.TransactionUUID(uuid.Nil) {
		transaction.Id = domain.TransactionUUID(uuid.New())
	}

	// NOTE: Right now we assume transactions unchangeable. So we not allow update
	query := `
        INSERT INTO transactions (id, username_from, username_to, amount) 
        VALUES ($1, $2, $3, $4)
    `

	if _, err := r.getter.DefaultTrOrDB(ctx, r.db).Exec(
		ctx,
		query,
		transaction.Id,
		transaction.From,
		transaction.To,
		transaction.Amount,
	); err != nil {
		return mapPgxError(err)
	}

	return nil
}
