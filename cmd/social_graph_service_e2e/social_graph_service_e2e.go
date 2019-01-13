package main

import (
	"context"
	_ "github.com/lib/pq"
	"github.com/the-gigi/delinkcious/pkg/db_util"
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

func initDB() {
	db, err := db_util.RunLocalDB("social_graph_manager")
	check(err)

	// Ignore if table doesn't exist (will be created by service)
	err = db_util.DeleteFromTableIfExist(db, "social_graph")
	check(err)
}

func main() {
	initDB()

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
