package sqpgx

import (
	"errors"
	"fmt"
	"time"
)

type SqpgxConfig struct {
	User     string
	Password string
	DbName   string
	Port     uint16
	Host     string

	MaxConnections         int
	MaxConnectionsLifetime time.Duration
}

const (
	DefaultMaxConnectionsLifetime = time.Minute
	DefaultMaxConnections         = 10
)

var (
	ErrIncorrectUser     = errors.New("incorrect user")
	ErrIncorrectPassword = errors.New("incorrect password")
	ErrIncorrectDbName   = errors.New("incorrect dbname")
	ErrIncorrectHost     = errors.New("incorrect host")
)

func NewSqpgxConfig(user, password, dbname, host string, port int) (*SqpgxConfig, error) {
	if user == "" {
		return nil, ErrIncorrectUser
	}
	if password == "" {
		return nil, ErrIncorrectPassword
	}
	if dbname == "" {
		return nil, ErrIncorrectDbName
	}
	if host == "" {
		return nil, ErrIncorrectHost
	}

	config := &SqpgxConfig{
		User:                   user,
		Password:               password,
		DbName:                 dbname,
		Host:                   host,
		Port:                   uint16(port),
		MaxConnections:         DefaultMaxConnections,
		MaxConnectionsLifetime: DefaultMaxConnectionsLifetime,
	}
	return config, nil
}

func (c *SqpgxConfig) GetConnectionString() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s pool_max_conns=%d pool_max_conn_lifetime=%s", c.Host, c.Port, c.User, c.Password, c.DbName, c.MaxConnections, c.MaxConnectionsLifetime.String())
}
