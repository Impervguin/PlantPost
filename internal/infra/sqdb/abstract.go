package sqdb

import (
	"context"

	"github.com/Masterminds/squirrel"
)

type Row interface {
	Scan(dest ...interface{}) error
}

type Rows interface {
	Next() bool
	Scan(dest ...interface{}) error
	Close()
	Err() error
}

type CommandTag interface {
	RowsAffected() int64
}

type SquirrelQuirier interface {
	QueryRow(ctx context.Context, sqb squirrel.SelectBuilder) (Row, error)
	Query(ctx context.Context, sqb squirrel.SelectBuilder) (Rows, error)
	Update(ctx context.Context, sqb squirrel.UpdateBuilder) (CommandTag, error)
	Insert(ctx context.Context, sqb squirrel.InsertBuilder) (CommandTag, error)
	Delete(ctx context.Context, sqb squirrel.DeleteBuilder) (CommandTag, error)
}
type SquirrelDatabase interface {
	SquirrelQuirier
	Transaction(ctx context.Context, tx func(SquirrelQuirier) error) error
}
