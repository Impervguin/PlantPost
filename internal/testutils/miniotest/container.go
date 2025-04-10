package miniotest

import (
	"context"
	"fmt"
	"time"

	"github.com/docker/go-connections/nat"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type TestLogConsumer struct {
	Msgs []string // store the logs as a slice of strings
}

func (g *TestLogConsumer) Accept(l testcontainers.Log) {
	g.Msgs = append(g.Msgs, string(l.Content))
}

func NewTestMinio(ctx context.Context) (testcontainers.Container, *MinioCredentials, error) {
	conf := GetConfig()
	strPort := fmt.Sprintf("%d/tcp", conf.Port)
	g := TestLogConsumer{
		Msgs: make([]string, 0),
	}
	req := testcontainers.ContainerRequest{
		Image:        conf.Image,
		ExposedPorts: []string{strPort},
		Env: map[string]string{
			"API_PORT":            fmt.Sprintf("%d", conf.Port),
			"MINIO_ROOT_USER":     conf.User,
			"MINIO_ROOT_PASSWORD": conf.Password,
		},
		Cmd: []string{"server", "/data", "--address", fmt.Sprintf(":%d", conf.Port)},
		WaitingFor: wait.ForAll(
			wait.ForLog("API:"),
			wait.ForListeningPort(nat.Port(strPort)),
		),
		LogConsumerCfg: &testcontainers.LogConsumerConfig{
			Opts:      []testcontainers.LogProductionOption{testcontainers.WithLogProductionTimeout(10 * time.Second)},
			Consumers: []testcontainers.LogConsumer{&g},
		},
	}
	cnt, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})

	logMsg := ""
	for _, msg := range g.Msgs {
		logMsg += msg
	}

	if err != nil {
		return nil, nil, fmt.Errorf("error creating minio container: %w %s", err, logMsg)
	}
	host, err := cnt.Host(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("error creating minio container: %w %s", err, logMsg)
	}
	port, err := cnt.MappedPort(ctx, nat.Port(strPort))
	if err != nil {
		return nil, nil, fmt.Errorf("error creating minio container: %w %s", err, logMsg)
	}
	time.Sleep(10 * time.Second)

	return cnt, NewMinioCredentials(conf.User, conf.Password, conf.Bucket, host, uint16(port.Int())), nil
}
