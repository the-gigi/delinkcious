package main

import (
	"context"
	_ "github.com/lib/pq"
	"github.com/the-gigi/delinkcious/pkg/db_util"
	om "github.com/the-gigi/delinkcious/pkg/object_model"
	"github.com/the-gigi/delinkcious/pkg/user_client"
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
	db, err := db_util.RunLocalDB("user_manager")
	if err != nil {
		return
	}

	tables := []string{"sessions", "users"}
	for _, table := range tables {
		err = db_util.DeleteFromTableIfExist(db, table)
		check(err)
	}
}

func runServer(ctx context.Context) {
	// Build the server if needed
	_, err := os.Stat("./user_service")
	if os.IsNotExist(err) {
		out, err := exec.Command("go", "build", ".").CombinedOutput()
		log.Println(out)
		check(err)
	}

	cmd := exec.CommandContext(ctx, "./user_service")
	err = cmd.Start()
	check(err)
}

func killServer(ctx context.Context) {
	ctx.Done()
}

func main() {
	initDB()

	ctx := context.Background()
	defer killServer(ctx)
	runServer(ctx)

	// Run some tests with the client
	cli, err := user_client.NewClient("localhost:7070")
	check(err)

	err = cli.Register(om.User{"gg@gg.com", "gigi"})
	check(err)
	log.Print("gigi has registered successfully")

	session, err := cli.Login("gigi", "secret")
	check(err)
	log.Print("gigi has logged in successfully. the session is: ", session)

	err = cli.Logout("gigi", session)
	check(err)
	log.Print("gigi has logged out successfully.")

}
