package sqpgx

import (
	"PlantSite/internal/infra/sqdb"
	"errors"

	"github.com/jackc/pgx/v5"
)

func MapError(err error) error {
	if errors.Is(err, pgx.ErrNoRows) {
		return sqdb.ErrNoRows
	}
	return err
}
