package database

import (
	"database/sql"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

//TODO: Add a way to handle errors in the database connection (e.g. if the database is down, or if the credentials are wrong, etc.
//TODO: Add rows.Close() to each query to prevent memory leaks
//TODO: clean code. Add comments to each function and struct

type Database struct {
	Pool *sql.DB
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

	instance := &Database{
		Pool: db,
	}

	return instance, nil
}

// create a function to easily query the database
func (db *Database) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return db.Pool.Query(query, args...)
}

func (db *Database) Execute(query string, args ...interface{}) (sql.Result, error) {
	return db.Pool.Exec(query, args)
}

func (db *Database) CloseDb() error {
	return db.Pool.Close()
}
