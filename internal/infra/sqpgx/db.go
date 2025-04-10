package sqpgx

import (
	"PlantSite/internal/infra/sqdb"
	"context"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SquirrelPgx struct {
	db *pgxpool.Pool
}

type squirrelPgxRow struct {
	pgx.Row
}

func (r *squirrelPgxRow) Scan(dest ...interface{}) error {
	err := r.Row.Scan(dest...)
	if err != nil {
		return MapError(err)
	}
	return nil
}

var _ sqdb.SquirrelDatabase = &SquirrelPgx{}

func NewSquirrelPgx(ctx context.Context, conf *SqpgxConfig) (*SquirrelPgx, error) {
	pgxConf, err := pgxpool.ParseConfig(conf.GetConnectionString())
	if err != nil {
		return nil, fmt.Errorf("parse config failed failed %w", err)
	}
	db, err := pgxpool.NewWithConfig(ctx, pgxConf)
	if err != nil {
		return nil, fmt.Errorf("connect to db failed %w", err)
	}
	err = db.Ping(ctx)
	if err != nil {
		return nil, fmt.Errorf("ping db failed %w", err)
	}
	return &SquirrelPgx{db: db}, nil
}

func (p *SquirrelPgx) QueryRow(ctx context.Context, sqb squirrel.SelectBuilder) (sqdb.Row, error) {
	sql, args, err := sqb.PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		return nil, fmt.Errorf("sql build failed %w", err)
	}
	row := p.db.QueryRow(ctx, sql, args...)
	return &squirrelPgxRow{
		Row: row,
	}, nil
}

func (p *SquirrelPgx) Query(ctx context.Context, sqb squirrel.SelectBuilder) (sqdb.Rows, error) {
	sql, args, err := sqb.PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		return nil, fmt.Errorf("sql build failed %w", err)
	}
	rows, err := p.db.Query(ctx, sql, args...)
	if err != nil {
		return nil, MapError(err)
	}
	return rows, nil
}

func (p *SquirrelPgx) Update(ctx context.Context, sqb squirrel.UpdateBuilder) (sqdb.CommandTag, error) {
	sql, args, err := sqb.PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		return pgconn.NewCommandTag(""), fmt.Errorf("sql build failed %w", err)
	}
	tag, err := p.db.Exec(ctx, sql, args...)
	if err != nil {
		return pgconn.NewCommandTag(""), MapError(err)
	}
	if tag.RowsAffected() == 0 {
		return tag, sqdb.ErrNoRows
	}
	return tag, nil
}

func (p *SquirrelPgx) Insert(ctx context.Context, sqb squirrel.InsertBuilder) (sqdb.CommandTag, error) {
	sql, args, err := sqb.PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		return pgconn.NewCommandTag(""), fmt.Errorf("sql build failed %w", err)
	}
	tag, err := p.db.Exec(ctx, sql, args...)
	if err != nil {
		return pgconn.NewCommandTag(""), MapError(err)
	}
	if tag.RowsAffected() == 0 {
		return tag, sqdb.ErrNoRows
	}
	return tag, nil
}

func (p *SquirrelPgx) Delete(ctx context.Context, sqb squirrel.DeleteBuilder) (sqdb.CommandTag, error) {
	sql, args, err := sqb.PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		return pgconn.NewCommandTag(""), fmt.Errorf("sql build failed %w", err)
	}
	tag, err := p.db.Exec(ctx, sql, args...)
	if err != nil {
		return pgconn.NewCommandTag(""), MapError(err)
	}
	if tag.RowsAffected() == 0 {
		return tag, sqdb.ErrNoRows
	}
	return tag, nil
}

func (p *SquirrelPgx) Transaction(ctx context.Context, tFunc func(sqdb.SquirrelQuirier) error) error {
	tx, err := p.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return fmt.Errorf("BeginTx failed %w", err)
	}
	defer tx.Rollback(ctx)
	sqtx := NewSquirrelTx(tx)
	err = tFunc(sqtx)

	if err != nil {
		return fmt.Errorf("transaction failed %w", err)
	}
	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("commit failed %w", err)
	}
	return nil
}
