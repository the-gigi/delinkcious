package db_util

import (
	"database/sql"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	_ "github.com/lib/pq"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type dbParams struct {
	Host     string
	Port     int
	User     string
	Password string
	DbName   string
}

func defaultDbParams() dbParams {
	return dbParams{
		Host:     "localhost",
		Port:     5432,
		User:     "postgres",
		Password: "postgres",
	}
}

func RunLocalDB(dbName string) (db *sql.DB, err error) {
	// Launch the DB if not running
	out, err := exec.Command("docker", "ps", "-f", "name=postgres", "--format", "{{.Names}}").CombinedOutput()
	if err != nil {
		return
	}

	s := string(out)
	if s == "" {
		out, err = exec.Command("docker", "restart", "postgres").CombinedOutput()
		if err != nil {
			log.Print(string(out))
			_, err = exec.Command("docker", "run", "-d", "--name", "postgres",
				"-p", "5432:5432",
				"-e", "POSTGRES_PASSWORD=postgres",
				"postgres:alpine").CombinedOutput()

		}
		if err != nil {
			return
		}
	}

	p := defaultDbParams()
	db, err = EnsureDB(p.Host, p.Port, p.User, p.Password, dbName)
	return
}

func connectToDB(host string, port int, username string, password string, dbName string) (db *sql.DB, err error) {
	mask := "host=%s port=%d user=%s password=%s dbname=%s sslmode=disable"
	dcn := fmt.Sprintf(mask, host, port, username, password, dbName)
	db, err = sql.Open("postgres", dcn)
	return
}

// Make sure the database exists (creates it if it doesn't)
func EnsureDB(host string, port int, username string, password string, dbName string) (db *sql.DB, err error) {
	// Connect to the postgres DB
	postgresDb, err := connectToDB(host, port, username, password, "postgres")
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

	db, err = connectToDB(host, port, username, password, dbName)
	return
}

func DeleteFromTableIfExist(db *sql.DB, table string) (err error) {
	_, err = db.Exec("DELETE from " + table)
	if err != nil {
		if err.Error() != fmt.Sprint("pq: relation \"%s\" does not exist", table) {
			return
		}
	}
	return
}

func GetDbEndpoint(dbName string) (host string, port int, err error) {
	hostEnvVar := strings.ToUpper(dbName) + "_DB_SERVICE_HOST"
	host = os.Getenv(hostEnvVar)
	if host == "" {
		host = "localhost"
	}

	portEnvVar := strings.ToUpper(dbName) + "_DB_SERVICE_PORT"
	dbPort := os.Getenv(portEnvVar)
	if dbPort == "" {
		dbPort = "5432"
	}

	port, err = strconv.Atoi(dbPort)
	return
}
