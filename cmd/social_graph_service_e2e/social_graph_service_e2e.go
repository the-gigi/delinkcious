package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/the-gigi/delinkcious/pkg/social_graph_client"
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
	//go exec.Command("./social_graph_service").CombinedOutput()

	// Run some tests with the client
	cli, err := social_graph_client.NewClient("localhost:9090")
	check(err)

	following, err := cli.GetFollowing("gigi")
	check(err)
	log.Print("gigi is following:", following)
	followers, err := cli.GetFollowers("gigi")
	check(err)
	log.Print("gigi is followed by:", followers)

	err = cli.Follow("gigi", "liat")
	check(err)
	err = cli.Follow("gigi", "guy")
	check(err)
	err = cli.Follow("guy", "gigi")
	check(err)
	err = cli.Follow("saar", "gigi")
	check(err)
	err = cli.Follow("saar", "ophir")
	check(err)

	following, err = cli.GetFollowing("gigi")
	check(err)
	log.Print("gigi is following:", following)
	followers, err = cli.GetFollowers("gigi")
	check(err)
	log.Print("gigi is followed by:", followers)

}
