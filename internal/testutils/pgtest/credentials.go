package pgtest

import "fmt"

type PostgresCredentials struct {
	User             string
	Password         string
	Database         string
	TemplateDatabase string
	Host             string
	Port             uint16
}

func getTemplateDatabase(database string) string {
	return fmt.Sprintf("template1_%s", database)
}

func NewPostgresCredentials(user, password, database, host string, port uint16) PostgresCredentials {
	return PostgresCredentials{
		User:             user,
		Password:         password,
		Database:         database,
		TemplateDatabase: getTemplateDatabase(database),
		Host:             host,
		Port:             port,
	}
}
