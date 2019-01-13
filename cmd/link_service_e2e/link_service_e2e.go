package main

import (
	"context"
	_ "github.com/lib/pq"
	"github.com/the-gigi/delinkcious/pkg/db_util"
	"github.com/the-gigi/delinkcious/pkg/link_manager_client"
	om "github.com/the-gigi/delinkcious/pkg/object_model"
	"log"
	"os"
	"os/exec"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func initDB() {
	db, err := db_util.RunLocalDB("link_manager")
	check(err)

	tables := []string{"tags", "links"}
	for _, table := range tables {
		err = db_util.DeleteFromTableIfExist(db, table)
		check(err)
	}
}

// Build and run a service in a target directory
func runService(ctx context.Context, targetDir string, service string) {
	// Save and restore lsater current working dir
	wd, err := os.Getwd()
	check(err)
	defer os.Chdir(wd)

	// Build the server if needed
	_, err = os.Stat("./" + service)
	if os.IsNotExist(err) {
		out, err := exec.Command("go", "build", ".").CombinedOutput()
		log.Println(out)
		check(err)
	}

	cmd := exec.CommandContext(ctx, "./"+service)
	err = cmd.Start()
	check(err)
}

func runLinkService(ctx context.Context) {
	runService(ctx, ".", "link_service")
}

func runSocialGraphService(ctx context.Context) {
	runService(ctx, "../social_graph_service", "link_service")
}

func killServer(ctx context.Context) {
	ctx.Done()
}

func main() {
	initDB()

	ctx := context.Background()
	defer killServer(ctx)
	runSocialGraphService(ctx)
	runLinkService(ctx)

	// Run some tests with the client
	cli, err := link_manager_client.NewClient("localhost:8080")
	check(err)

	links, err := cli.GetLinks(om.GetLinksRequest{Username: "gigi"})
	check(err)
	log.Print("gigi's links:", links)

	err = cli.AddLink(om.AddLinkRequest{Username: "gigi",
		Url:   "https://github.com/the-gigi",
		Title: "Gigi on Github",
		Tags:  map[string]bool{"programming": true}})
	check(err)
	links, err = cli.GetLinks(om.GetLinksRequest{Username: "gigi"})
	check(err)
	log.Print("gigi's links:", links)

	err = cli.UpdateLink(om.UpdateLinkRequest{Username: "gigi",
		Url:         "https://github.com/the-gigi",
		Description: "Most of my open source code is here"},
	)

	check(err)
	links, err = cli.GetLinks(om.GetLinksRequest{Username: "gigi"})
	check(err)
	log.Print("gigi's links:", links)
}
