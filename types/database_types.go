package types

import "database/sql"

type Database struct {
	Pool *sql.DB
}

type DatabaseOptions struct {
	User     string
	Password string
	Host     string
	Port     string
	Name     string
}
