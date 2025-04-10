package pgtest

import (
	"context"
	"fmt"

	"github.com/docker/go-connections/nat"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func NewTestPostgres(ctx context.Context) (testcontainers.Container, PostgresCredentials, error) {
	config := GetConfig()
	strPort := fmt.Sprintf("%d/tcp", config.Port)
	req := testcontainers.ContainerRequest{
		Image:        config.Image,
		ExposedPorts: []string{strPort},
		Env: map[string]string{
			"POSTGRES_USER":     config.User,
			"POSTGRES_PASSWORD": config.Password,
			"POSTGRES_DB":       config.Database,
		},
		WaitingFor: wait.ForAll(
			wait.ForLog("database system is ready to accept connections"),
			wait.ForListeningPort(nat.Port(strPort)),
		),
	}
	cnt, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, PostgresCredentials{}, err
	}
	host, err := cnt.Host(ctx)
	if err != nil {
		return nil, PostgresCredentials{}, err
	}
	port, err := cnt.MappedPort(ctx, "5432/tcp")
	if err != nil {
		return nil, PostgresCredentials{}, err
	}
	creds := NewPostgresCredentials(config.User, config.Password, config.Database, host, uint16(port.Int()))
	return cnt, creds, nil
}
