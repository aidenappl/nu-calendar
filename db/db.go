package db

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/aidenappl/nu-calendar/env"
	"github.com/go-sql-driver/mysql"
)

const (
	DefaultListLimit = 50
	MaximumListLimit = 100

	ErNoReferencedRow     = 1215
	ErDupEntry            = 1062
	ErDupEntryWithKeyName = 1586
)

func PingDB(db *sql.DB) error {
	fmt.Print("Connecting CoreClient to Core DB...")
	err := db.Ping()
	fmt.Println(" ✅ Done")
	return err
}

var DB = func() *sql.DB {
	db, err := sql.Open("mysql", env.CoreDBDSN)
	if err != nil {
		panic(fmt.Errorf("error opening database: %w", err))
	}

	return db
}()

type Queryable interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	Prepare(query string) (*sql.Stmt, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
}

func ExtractDBErrorCode(err error) uint16 {
	var sqlErr *mysql.MySQLError
	if errors.As(err, &sqlErr) {
		return sqlErr.Number
	} else {
		return 0
	}
}
