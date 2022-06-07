package driver

import (
	"database/sql"
	"time"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

//DB holds the database connection pool
type DB struct {
	SQL *sql.DB
}

var dbConn = &DB{}

const maxOpenDbConn = 10
const maxIdleDbConn = 5
const maxDbLifetime = 5 * time.Minute

func ConnectSQL(dsn string) (*DB, error) {

	dbConnection, err := NewDatabase(dsn)
	if err != nil {
		panic(err)
	}

	dbConnection.SetMaxOpenConns(maxOpenDbConn)
	dbConnection.SetConnMaxIdleTime(maxIdleDbConn)
	dbConnection.SetConnMaxLifetime(maxDbLifetime)

	dbConn.SQL = dbConnection

	err = testDB(dbConnection)
	if err != nil {
		return nil, err
	}

	return dbConn, nil

}

// testDB tries to ping the database
func testDB(d *sql.DB) error {
	err := d.Ping()
	if err != nil {
		return err
	}

	return nil

}

// NewDatabase creates a new database for the application
func NewDatabase(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)

	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
