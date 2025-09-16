package pgcnt

import (
	"context"

	"github.com/docker/go-connections/nat"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

const (
	PostgresPort                = "5432/tcp"
	PostgresPortInternal uint16 = 5432
)

func NewTestPostgres(ctx context.Context, configPath string, network string) (testcontainers.Container, *PostgresConfig, error) {
	config, err := GetConfig(configPath)
	if err != nil {
		return nil, nil, err
	}
	req := testcontainers.ContainerRequest{
		FromDockerfile: testcontainers.FromDockerfile{},
		Image:          config.Image,
		ExposedPorts:   []string{PostgresPort},
		Env: map[string]string{
			"POSTGRES_USER":     config.User,
			"POSTGRES_PASSWORD": config.Password,
			"POSTGRES_DB":       config.Database,
		},
		WaitingFor: wait.ForAll(
			wait.ForLog("database system is ready to accept connections"),
			wait.ForListeningPort(nat.Port(PostgresPort)),
		),
		Networks: []string{network},
		Hostname: config.Host,
	}
	cnt, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, nil, err
	}
	host, err := cnt.Host(ctx)
	if err != nil {
		return nil, nil, err
	}
	port, err := cnt.MappedPort(ctx, PostgresPort)
	if err != nil {
		return nil, nil, err
	}
	config.OuterHost = &host
	p := uint16(port.Int())
	config.OuterPort = &p
	config.Port = PostgresPortInternal

	return cnt, config, nil
}
