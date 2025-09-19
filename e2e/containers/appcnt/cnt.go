package appcnt

import (
	"context"
	"fmt"

	"github.com/docker/go-connections/nat"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

const (
	AppPort = "23450/tcp"
)

func NewTestApp(config *AppConfig, network string) (testcontainers.Container, *AppConfig, error) {
	var ParsedPort uint16
	fmt.Println("Parsing port")
	fmt.Sscanf(AppPort, "%d", &ParsedPort)
	var exposedPorts []string
	exposedPorts = []string{AppPort}
	req := testcontainers.ContainerRequest{
		FromDockerfile: testcontainers.FromDockerfile{
			Context:    config.BuildContext,
			Dockerfile: config.DockerFilePath,
		},
		ExposedPorts: exposedPorts,
		Files: []testcontainers.ContainerFile{
			{
				ContainerFilePath: config.AppConfigPath,
				FileMode:          0644,
				Reader:            config.GetApiYamlReader(ParsedPort),
			},
		},
		WaitingFor: wait.ForAll(
			wait.ForListeningPort(nat.Port(AppPort)),
		),
		Networks: []string{network},
		Hostname: config.Host,
	}
	cnt, err := testcontainers.GenericContainer(context.Background(), testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, nil, err
	}
	host, err := cnt.Host(context.Background())
	if err != nil {
		return nil, nil, err
	}
	port, err := cnt.MappedPort(context.Background(), AppPort)
	if err != nil {
		return nil, nil, err
	}
	config.OuterHost = &host
	p := uint16(port.Int())
	config.OuterPort = &p

	return cnt, config, nil
}
