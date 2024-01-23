package database

import (
	"database/sql"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

type Database struct {
	pool *sql.DB
}

type databaseOptions struct {
	User     string
	Password string
	Host     string
	Port     string
	Name     string
}

type ControllerInterface interface {
}

func Connect() (*Database, error) {

	databaseOptions := databaseOptions{
		User:     os.Getenv("DATABASE_USER"),
		Password: os.Getenv("DATABASE_PASSWORD"),
		Host:     os.Getenv("DATABASE_HOST"),
		Port:     os.Getenv("DATABASE_PORT"),
		Name:     os.Getenv("DATABASE"),
	}

	db, err := sql.Open("mysql", databaseOptions.User+":"+databaseOptions.Password+"@tcp("+databaseOptions.Host+":"+databaseOptions.Port+")/"+databaseOptions.Name)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return &Database{
		pool: db,
	}, nil
}

// create a function to easily query the database
func (db *Database) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return db.pool.Query(query, args...)
}
