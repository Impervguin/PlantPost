package sqpgx

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SquirrelPgx struct {
	db *pgxpool.Pool
}

func NewSquirrelPgx(db *pgxpool.Pool) *SquirrelPgx {
	return &SquirrelPgx{db: db}
}

func (p *SquirrelPgx) QueryRow(ctx context.Context, sqb squirrel.SelectBuilder) pgx.Row {
	sql, args, _ := sqb.PlaceholderFormat(squirrel.Dollar).ToSql()
	return p.db.QueryRow(ctx, sql, args...)
}

func (p *SquirrelPgx) Query(ctx context.Context, sqb squirrel.SelectBuilder) (pgx.Rows, error) {
	sql, args, err := sqb.PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		return nil, err
	}
	return p.db.Query(ctx, sql, args...)
}

func (p *SquirrelPgx) Update(ctx context.Context, sqb squirrel.UpdateBuilder) (pgconn.CommandTag, error) {
	sql, args, err := sqb.PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		return pgconn.NewCommandTag(""), err
	}
	return p.db.Exec(ctx, sql, args...)
}

func (p *SquirrelPgx) Insert(ctx context.Context, sqb squirrel.InsertBuilder) (pgconn.CommandTag, error) {
	sql, args, err := sqb.PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		return pgconn.NewCommandTag(""), err
	}
	return p.db.Exec(ctx, sql, args...)
}

func (p *SquirrelPgx) Delete(ctx context.Context, sqb squirrel.DeleteBuilder) (pgconn.CommandTag, error) {
	sql, args, err := sqb.PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		return pgconn.NewCommandTag(""), err
	}
	return p.db.Exec(ctx, sql, args...)
}

func (p *SquirrelPgx) BeginTx(ctx context.Context, opts pgx.TxOptions) (*SquirrelTx, error) {
	tx, err := p.db.BeginTx(ctx, opts)
	if err != nil {
		return nil, err
	}
	return NewSquirrelTx(tx), nil
}
