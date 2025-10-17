package pgtest

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
)

var ErrTemplateDatabaseAlreadyExists = errors.New("template database already exists")
var ErrTemplateDatabaseDoesNotExist = errors.New("template database does not exist")

func CreateTemplateDatabase(ctx context.Context, db *PostgresCredentials) error {
	conn, err := pgx.Connect(
		ctx,
		fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", db.User, db.Password, db.Host, db.Port, db.Database),
	)
	if err != nil {
		return err
	}
	defer conn.Close(ctx)
	var tmp int
	err = conn.QueryRow(ctx, fmt.Sprintf("SELECT 1 FROM pg_database WHERE datname = '%s'", db.TemplateDatabase)).
		Scan(&tmp)
	if !errors.Is(err, pgx.ErrNoRows) {
		return ErrTemplateDatabaseAlreadyExists
	}

	_, err = conn.Exec(ctx, fmt.Sprintf("CREATE DATABASE %s", db.TemplateDatabase))
	if err != nil {
		return err
	}
	return nil
}

func DropTemplateDatabase(ctx context.Context, db *PostgresCredentials) error {
	conn, err := pgx.Connect(
		ctx,
		fmt.Sprintf("postgres://%s:%s@%s:%d?sslmode=disable", db.User, db.Password, db.Host, db.Port),
	)
	if err != nil {
		return err
	}
	defer conn.Close(ctx)
	var tmp int
	err = conn.QueryRow(ctx, fmt.Sprintf("SELECT 1 FROM pg_database WHERE datname = '%s'", db.TemplateDatabase)).
		Scan(&tmp)
	if errors.Is(err, pgx.ErrNoRows) {
		return ErrTemplateDatabaseDoesNotExist
	}
	_, err = conn.Exec(ctx, fmt.Sprintf("DROP DATABASE %s", db.TemplateDatabase))
	if err != nil {
		return err
	}
	return nil
}

func ApplyTemplate(ctx context.Context, db *PostgresCredentials) error {
	conn, err := pgx.Connect(
		ctx,
		fmt.Sprintf(
			"postgres://%s:%s@%s:%d/%s?sslmode=disable",
			db.User,
			db.Password,
			db.Host,
			db.Port,
			db.TemplateDatabase,
		),
	)
	if err != nil {
		return err
	}
	defer conn.Close(ctx)
	var tmp int
	err = conn.QueryRow(ctx, fmt.Sprintf("SELECT 1 FROM pg_database WHERE datname = '%s'", db.TemplateDatabase)).
		Scan(&tmp)
	if errors.Is(err, pgx.ErrNoRows) {
		return ErrTemplateDatabaseDoesNotExist
	}
	_, err = conn.Exec(ctx, fmt.Sprintf("DROP DATABASE %s with (FORCE)", db.Database))
	if err != nil {
		return err
	}
	_, err = conn.Exec(ctx, fmt.Sprintf("CREATE DATABASE %s WITH TEMPLATE %s", db.Database, db.TemplateDatabase))
	if err != nil {
		return err
	}
	return nil
}
