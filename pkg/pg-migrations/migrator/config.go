package migrator

import (
	"fmt"
)

type Config struct {
	MigrationDir string
	Database     string
	User         string
	Password     string
	Host         string
	Port         int
}

func (c Config) Validate() error {
	if c.MigrationDir == "" {
		return fmt.Errorf("MigrationDir is required")
	}
	if c.Database == "" {
		return fmt.Errorf("Database is required")
	}
	if c.User == "" {
		return fmt.Errorf("User is required")
	}
	if c.Password == "" {
		return fmt.Errorf("Password is required")
	}
	if c.Host == "" {
		return fmt.Errorf("Host is required")
	}
	if c.Port <= 0 {
		return fmt.Errorf("Port is required")
	}
	return nil
}
