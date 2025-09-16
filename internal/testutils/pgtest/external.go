package pgtest

import (
	"context"
	"fmt"
)

func NewTestExternalPostgres(ctx context.Context, config *PostgresConfig) (*PostgresCredentials, error) {
	if !config.External {
		return nil, fmt.Errorf("postgres is not configured")
	}
	creds := NewPostgresCredentials(config.User, config.Password, config.Database, config.Host, config.Port)
	return &creds, nil
}
