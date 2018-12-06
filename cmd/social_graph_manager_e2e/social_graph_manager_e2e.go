package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"os/exec"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	// Launch the DB if not running
	out, err := exec.Command("docker", "ps", "-f", "name=postgres", "--format", "{{.Names}}").CombinedOutput()
	check(err)

	s := string(out)
	log.Print(s)

	if s == "" {
		_, err := exec.Command("docker", "restart", "postgres").CombinedOutput()
		check(err)
	}

	// Clear the DB
	mask := "host=%s port=%d user=%s password=%s dbname=social_graph_manager sslmode=disable"
	dcn := fmt.Sprintf(mask, "localhost", 5432, "postgres", "postgres")
	db, err := sql.Open("postgres", dcn)
	if err != nil {
		return
	}

	_, err = db.Exec("DELETE from social_graph")
	check(err)

	// Launch the server

	// Run some tests with the client
}
