package main

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/the-gigi/delinkcious/pkg/social_graph_client"
	"log"
	"os"
	"os/exec"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func runDB() {
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
}

func runServer(ctx context.Context) {
	// Build the server if needed
	_, err := os.Stat("./social_graph_service")
	if os.IsNotExist(err) {
		out, err := exec.Command("go", "build", ".").CombinedOutput()
		log.Println(out)
		check(err)
	}

	cmd := exec.CommandContext(ctx, "./social_graph_service")
	err = cmd.Start()
	check(err)
}

func killServer(ctx context.Context) {
	ctx.Done()
}

func main() {
	runDB()

	ctx := context.Background()
	defer killServer(ctx)
	runServer(ctx)

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
