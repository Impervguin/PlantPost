package sqpgx

import (
	"PlantSite/internal/infra/sqdb"
	"context"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type SquirrelTx struct {
	tx pgx.Tx
}

func (tx *SquirrelTx) QueryRow(ctx context.Context, sqb squirrel.SelectBuilder) (sqdb.Row, error) {
	sql, args, err := sqb.PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		return nil, fmt.Errorf("sql build failed %w", err)
	}
	row := tx.tx.QueryRow(ctx, sql, args...)

	return &squirrelPgxRow{row}, nil
}

func (tx *SquirrelTx) Query(ctx context.Context, sqb squirrel.SelectBuilder) (sqdb.Rows, error) {
	sql, args, err := sqb.PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		return nil, fmt.Errorf("sql build failed %w", err)
	}
	rows, err := tx.tx.Query(ctx, sql, args...)
	if err != nil {
		return nil, MapError(err)
	}
	return rows, nil
}

func (tx *SquirrelTx) Update(ctx context.Context, sqb squirrel.UpdateBuilder) (sqdb.CommandTag, error) {
	sql, args, err := sqb.PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		return pgconn.NewCommandTag(""), fmt.Errorf("sql build failed %w", err)
	}
	commandTag, err := tx.tx.Exec(ctx, sql, args...)
	if err != nil {
		return pgconn.NewCommandTag(""), MapError(err)
	}
	if commandTag.RowsAffected() == 0 {
		return pgconn.NewCommandTag(""), sqdb.ErrNoRows
	}
	return commandTag, nil
}

func (tx *SquirrelTx) Insert(ctx context.Context, sqb squirrel.InsertBuilder) (sqdb.CommandTag, error) {
	sql, args, err := sqb.PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		return pgconn.NewCommandTag(""), MapError(err)
	}
	commandTag, err := tx.tx.Exec(ctx, sql, args...)
	if err != nil {
		return pgconn.NewCommandTag(""), MapError(err)
	}
	if commandTag.RowsAffected() == 0 {
		return pgconn.NewCommandTag(""), sqdb.ErrNoRows
	}
	return commandTag, nil
}

func (tx *SquirrelTx) Delete(ctx context.Context, sqb squirrel.DeleteBuilder) (sqdb.CommandTag, error) {
	sql, args, err := sqb.PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		return pgconn.NewCommandTag(""), MapError(err)
	}
	commandTag, err := tx.tx.Exec(ctx, sql, args...)
	if err != nil {
		return pgconn.NewCommandTag(""), MapError(err)
	}
	if commandTag.RowsAffected() == 0 {
		return pgconn.NewCommandTag(""), sqdb.ErrNoRows
	}
	return commandTag, nil
}

func (tx *SquirrelTx) Commit(ctx context.Context) error {
	return tx.tx.Commit(ctx)
}

func (tx *SquirrelTx) Rollback(ctx context.Context) error {
	return tx.tx.Rollback(ctx)
}

func NewSquirrelTx(tx pgx.Tx) *SquirrelTx {
	return &SquirrelTx{tx: tx}
}
