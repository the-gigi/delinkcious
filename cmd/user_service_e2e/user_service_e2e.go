package main

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
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
	mask := "host=%s port=%d user=%s password=%s dbname=user_manager sslmode=disable"
	dcn := fmt.Sprintf(mask, "localhost", 5432, "postgres", "postgres")
	db, err := sql.Open("postgres", dcn)
	if err != nil {
		return
	}

	_, err = db.Exec("DELETE from sessions")
	check(err)
	_, err = db.Exec("DELETE from users")
	check(err)
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
	runDB()

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
