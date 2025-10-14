package pgcnt

import (
	"context"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func Migrate(ctx context.Context, db *PostgresConfig) error {
	sourceUrl := fmt.Sprintf("file://%s", db.MigrationDir)
	dbUrl := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", db.User, db.Password, *db.OuterHost, *db.OuterPort, db.Database)
	m, err := migrate.New(sourceUrl, dbUrl)
	if err != nil {
		return err
	}
	defer m.Close()
	err = m.Up()
	if err != nil {
		return err
	}
	return nil
}

func MigrateDown(ctx context.Context, db *PostgresConfig) error {
	sourceUrl := fmt.Sprintf("file://%s", db.MigrationDir)
	dbUrl := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", db.User, db.Password, *db.OuterHost, *db.OuterPort, db.Database)
	m, err := migrate.New(sourceUrl, dbUrl)
	if err != nil {
		return err
	}
	defer m.Close()
	err = m.Down()
	if err != nil {
		return err
	}
	return nil
}
