package pgtest

type PostgresCredentials struct {
	User     string
	Password string
	Database string
	Host     string
	Port     uint16
}

func NewPostgresCredentials(user, password, database, host string, port uint16) PostgresCredentials {
	return PostgresCredentials{
		User:     user,
		Password: password,
		Database: database,
		Host:     host,
		Port:     port,
	}
}
