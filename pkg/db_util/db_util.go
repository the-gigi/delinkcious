package db_util

import (
	"database/sql"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	_ "github.com/lib/pq"
)

func connectToDb(host string, port int, username string, password string, dbName string) (db *sql.DB, err error) {
	mask := "host=%s port=%d user=%s password=%s dbname=%s sslmode=disable"
	dcn := fmt.Sprintf(mask, host, port, username, password, dbName)
	db, err = sql.Open("postgres", dcn)
	return
}

// Make sure the database exists (creates it if it doesn't)
func EnsureDB(host string, port int, username string, password string, dbName string) (db *sql.DB, err error) {
	// Connect to the postgres DB
	postgresDb, err := connectToDb(host, port, username, password, "postgres")
	if err != nil {
		return
	}

	// Check if the DB exists in the list of databases
	var count int
	sb := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	q := sb.Select("count(*)").From("pg_database").Where(sq.Eq{"datname": dbName})
	err = q.RunWith(postgresDb).QueryRow().Scan(&count)
	if err != nil {
		return
	}

	// If it doesn't exist create it
	if count == 0 {
		_, err = postgresDb.Exec("CREATE database " + dbName)
		if err != nil {
			return
		}
	}

	db, err = connectToDb(host, port, username, password, dbName)
	return
}
