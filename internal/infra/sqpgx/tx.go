package sqpgx

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type SquirrelTx struct {
	tx pgx.Tx
}

func (tx *SquirrelTx) QueryRow(ctx context.Context, sqb squirrel.SelectBuilder) pgx.Row {
	sql, args, _ := sqb.PlaceholderFormat(squirrel.Dollar).ToSql()
	return tx.tx.QueryRow(ctx, sql, args...)
}

func (tx *SquirrelTx) Query(ctx context.Context, sqb squirrel.SelectBuilder) (pgx.Rows, error) {
	sql, args, err := sqb.PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		return nil, err
	}
	return tx.tx.Query(ctx, sql, args...)
}

func (tx *SquirrelTx) Update(ctx context.Context, sqb squirrel.UpdateBuilder) (pgconn.CommandTag, error) {
	sql, args, err := sqb.PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		return pgconn.NewCommandTag(""), err
	}
	return tx.tx.Exec(ctx, sql, args...)
}

func (tx *SquirrelTx) Insert(ctx context.Context, sqb squirrel.InsertBuilder) (pgconn.CommandTag, error) {
	sql, args, err := sqb.PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		return pgconn.NewCommandTag(""), err
	}
	return tx.tx.Exec(ctx, sql, args...)
}

func (tx *SquirrelTx) Delete(ctx context.Context, sqb squirrel.DeleteBuilder) (pgconn.CommandTag, error) {
	sql, args, err := sqb.PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		return pgconn.NewCommandTag(""), err
	}
	return tx.tx.Exec(ctx, sql, args...)
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
