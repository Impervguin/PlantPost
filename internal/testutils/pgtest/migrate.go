package pgtest

import (
	"context"
	"fmt"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func Migrate(ctx context.Context, db *PostgresCredentials, database string) error {
	config := GetConfig()
	sourceUrl := fmt.Sprintf("file://%s", config.MigrationDir)
	dbUrl := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", db.User, db.Password, db.Host, db.Port, database)
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

func MigrateDown(ctx context.Context, db *PostgresCredentials, database string) error {
	config := GetConfig()
	sourceUrl := fmt.Sprintf("file://%s", config.MigrationDir)
	dbUrl := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", db.User, db.Password, db.Host, db.Port, database)
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

func getVersionFromFilename(filename string) (uint, error) {
	var version uint
	_, err := fmt.Sscanf(filename, "%d", &version)
	if err != nil {
		return 0, err
	}
	return version, nil
}

func getMaxSourceVersion(ctx context.Context, db *PostgresCredentials) (uint, error) {
	config := GetConfig()
	entries, err := os.ReadDir(config.MigrationDir)
	if err != nil {
		return 0, err
	}
	var maxVersion uint = 0
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		version, err := getVersionFromFilename(entry.Name())
		if err != nil {
			continue
		}
		if version > maxVersion {
			maxVersion = version
		}
	}
	return maxVersion, nil
}

func CheckMigrationVersion(ctx context.Context, db *PostgresCredentials, database string) error {
	config := GetConfig()
	sourceUrl := fmt.Sprintf("file://%s", config.MigrationDir)
	dbUrl := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", db.User, db.Password, db.Host, db.Port, database)
	m, err := migrate.New(sourceUrl, dbUrl)
	if err != nil {
		return err
	}
	defer m.Close()
	dbVersion, _, err := m.Version()
	if err != nil {
		return err
	}

	sourceVersion, err := getMaxSourceVersion(ctx, db)
	if err != nil {
		return err
	}
	if sourceVersion != dbVersion {
		return fmt.Errorf("migration version mismatch: source version %d, db version %d", sourceVersion, dbVersion)
	}
	return nil
}
