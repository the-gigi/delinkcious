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


type dbParams struct {
	Host string
	Port int
	User string
	Password string
	DbName string
}

func defaultDbParams() dbParams {
	return dbParams{
		Host: "localhost",
		Port: 5432,
		User: "postgres",
		Password: "postgres",
	}
}

func initDb(params dbParams) {
	// Connect to target database
	mask := "host=%s port=%d user=%s password=%s dbname=%s sslmode=disable"
	dcn := fmt.Sprintf(mask, params.Host, params.Port, params.User, params.Password, params.DbName)

	db, err := sql.Open("postgres", dcn)
	if err != nil {
		return
	}

	// Ignore is table doesn't exist (will be created by service)
	_, err = db.Exec("DELETE from social_graph")
	if err != nil {
		if err.Error() != "pq: relation \"social_graph\" does not exist" {
			check(err)
		}
	}
}

func runDB() {
	// Launch the DB if not running
	out, err := exec.Command("docker", "ps", "-f", "name=postgres", "--format", "{{.Names}}").CombinedOutput()
	check(err)

	s := string(out)
	log.Print(s)

	if s == "" {
		out, err = exec.Command("docker", "restart", "postgres").CombinedOutput()
		if err != nil {
			log.Print(string(out))
			_, err = exec.Command("docker", "run", "-d", "--name", "postgres",
				                        "-p", "5432:5432",
				                        "-e", "POSTGRES_PASSWORD=postgres",
				                        "postgres:alpine").CombinedOutput()

		}
		check(err)
	}

	params := defaultDbParams()
	params.DbName = "social_graph_manager"
	initDb(params)
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
