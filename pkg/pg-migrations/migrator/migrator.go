package migrator

import (
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type Migrator struct {
	config Config
}

func NewMigrator(config *Config) (*Migrator, error) {
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}
	return &Migrator{
		config: *config,
	}, nil
}

func (m *Migrator) Up() error {
	sourceURL := fmt.Sprintf("file://%s", m.config.MigrationDir)
	dbURL := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=disable",
		m.config.User,
		m.config.Password,
		m.config.Host,
		m.config.Port,
		m.config.Database,
	)
	mig, err := migrate.New(sourceURL, dbURL)
	if err != nil {
		return err
	}
	defer mig.Close()

	log := NewLogger()

	mig.Log = log
	err = mig.Up()
	if err != nil {
		return err
	}
	return nil
}

func (m *Migrator) Down() error {
	sourceURL := fmt.Sprintf("file://%s", m.config.MigrationDir)
	dbURL := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=disable",
		m.config.User,
		m.config.Password,
		m.config.Host,
		m.config.Port,
		m.config.Database,
	)
	mig, err := migrate.New(sourceURL, dbURL)
	if err != nil {
		return err
	}
	defer mig.Close()

	log := NewLogger()

	mig.Log = log
	err = mig.Down()
	if err != nil {
		return err
	}
	return nil
}
