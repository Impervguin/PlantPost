package filters

import (
	"github.com/Masterminds/squirrel"
)

type PostgresPlantFilter squirrel.Sqlizer

type PostgresPostFilter squirrel.Sqlizer

type PostgresPlantSearch squirrel.And

type PostgresPostSearch squirrel.And

func NewPostgresPlantSearch() (*PostgresPlantSearch, error) {
	return &PostgresPlantSearch{}, nil
}

func (ps *PostgresPlantSearch) AddFilter(filter PostgresPlantFilter) error {
	*ps = append(*ps, filter)
	return nil
}

func (ps *PostgresPlantSearch) ToSql() (string, []interface{}, error) {
	and := squirrel.And{}
	for _, filter := range *ps {
		and = append(and, filter)
	}
	query, args, err := and.ToSql()
	return query, args, err
}

func NewPostgresPostSearch() (*PostgresPostSearch, error) {
	return &PostgresPostSearch{}, nil
}

func (ps *PostgresPostSearch) AddFilter(filter PostgresPostFilter) error {
	*ps = append(*ps, filter)
	return nil
}

func (ps *PostgresPostSearch) ToSql() (string, []interface{}, error) {
	and := squirrel.And{}
	for _, filter := range *ps {
		and = append(and, filter)
	}
	query, args, err := and.ToSql()
	return query, args, err
}
