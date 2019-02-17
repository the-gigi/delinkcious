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
	// Save and restore later current working dir
	wd, err := os.Getwd()
	check(err)
	defer os.Chdir(wd)

	// Build the server if needed
	os.Chdir(targetDir)
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
	// Set environment
	err := os.Setenv("PORT", "8080")
	check(err)

	err = os.Setenv("MAX_LINKS_PER_USER", "10")
	check(err)

	runService(ctx, ".", "link_service")
}

func runSocialGraphService(ctx context.Context) {
	err := os.Setenv("PORT", "9090")
	check(err)

	runService(ctx, "../social_graph_service", "social_graph_service")
}

func killServer(ctx context.Context) {
	ctx.Done()
}

func main() {
	// Turn on authentication
	err := os.Setenv("DELINKCIOUS_MUTUAL_AUTH", "true")
	check(err)

	initDB()

	ctx := context.Background()
	defer killServer(ctx)

	if os.Getenv("RUN_SOCIAL_GRAPH_SERVICE") == "true" {
		runSocialGraphService(ctx)
	}

	if os.Getenv("RUN_LINK_SERVICE") == "true" {
		runLinkService(ctx)
	}

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
