package adapters

import (
	"avito-test/internal/domain"
	"context"

	"github.com/jackc/pgx/v4/pgxpool"

	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv4/v2"
)

type userRepository struct {
	db     *pgxpool.Pool
	getter *trmpgx.CtxGetter
}

func NewUserRepository(db *pgxpool.Pool, c *trmpgx.CtxGetter) domain.UserRepository {
	return &userRepository{db, c}
}

func (r *userRepository) GetByUsername(ctx context.Context, username domain.UserName) (*domain.User, error) {
	query := `SELECT password, balance FROM Users WHERE username=$1`

	row := r.getter.DefaultTrOrDB(ctx, r.db).QueryRow(ctx, query, username)

	user := &domain.User{Username: username}

	err := row.Scan(&user.PasswordHash, &user.Balance)
	if err != nil {
		return nil, mapPgxError(err)
	}

	return user, nil
}

func (r *userRepository) Save(ctx context.Context, user *domain.User) error {
	query := `
        INSERT INTO Users (username, password, balance) 
        VALUES ($1, $2, $3)
        ON CONFLICT (username) 
        DO UPDATE SET
				password = EXCLUDED.password,
				balance = EXCLUDED.balance
    `

	if _, err := r.getter.DefaultTrOrDB(ctx, r.db).Exec(ctx, query, user.Username, user.PasswordHash, user.Balance); err != nil {
		return mapPgxError(err)
	}

	return nil
}
