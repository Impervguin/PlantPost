package miniocnt

import (
	"context"
	"fmt"

	"github.com/docker/go-connections/nat"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

const (
	MinioPort = "9000/tcp"
)

func NewTestMinio(ctx context.Context, configPath string, network string) (testcontainers.Container, *MinioConfig, error) {
	config, err := GetConfig(configPath)
	if err != nil {
		return nil, nil, err
	}
	// var ParsedPort uint16
	// fmt.Sscanf(MinioPort, "%d", &ParsedPort)
	req := testcontainers.ContainerRequest{
		Image:        config.Image,
		ExposedPorts: []string{MinioPort},
		Cmd:          []string{"server", "/data", "--address", fmt.Sprintf(":%d", config.Port)},
		Env: map[string]string{
			"API_PORT":            fmt.Sprintf("%d", config.Port),
			"MINIO_ROOT_USER":     config.User,
			"MINIO_ROOT_PASSWORD": config.Password,
		},
		WaitingFor: wait.ForAll(
			wait.ForLog("API:"),
			wait.ForListeningPort(nat.Port(MinioPort)),
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
	port, err := cnt.MappedPort(ctx, MinioPort)
	if err != nil {
		return nil, nil, err
	}
	config.OuterHost = &host
	p := uint16(port.Int())
	config.OuterPort = &p

	return cnt, config, nil
}
