package minioclient

import "errors"

type MinioConfig struct {
	Endpoint string
	User     string
	Password string
	Bucket   string
}

var (
	ErrIncorrectEndpoint = errors.New("incorrect endpoint")
	ErrIncorrectUser     = errors.New("incorrect user")
	ErrIncorrectPassword = errors.New("incorrect password")
	ErrIncorrectBucket   = errors.New("incorrect bucket")
)

func NewMinioConfig(endpoint, user, password, bucket string) (*MinioConfig, error) {
	if endpoint == "" {
		return nil, ErrIncorrectEndpoint
	}
	if user == "" {
		return nil, ErrIncorrectUser
	}
	if password == "" {
		return nil, ErrIncorrectPassword
	}
	if bucket == "" {
		return nil, ErrIncorrectBucket
	}
	return &MinioConfig{
		Endpoint: endpoint,
		User:     user,
		Password: password,
		Bucket:   bucket,
	}, nil
}
