package adapters

import (
	"avito-test/internal/domain"

	"github.com/jackc/pgx/v4"
)

func mapPgxError(pgxErr error) error {
	if pgxErr == pgx.ErrNoRows {
		return domain.ErrEntityDoesNotExist
	}

	// TODO default data sanitizing
	return pgxErr
}
