package adapters

import (
	"avito-test/internal/domain"
	"errors"

	"github.com/jackc/pgx/v4"
)

func mapPgxError(pgxErr error) error {
	if errors.Is(pgxErr, pgx.ErrNoRows) {
		return domain.ErrEntityDoesNotExist
	}

	return pgxErr
}
