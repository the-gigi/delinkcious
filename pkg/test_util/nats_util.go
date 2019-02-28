package test_util

import (
	"log"
	"os/exec"
)

func RunLocalNatsServer() (err error) {
	// Launch the DB if not running
	out, err := exec.Command("docker", "ps", "-f", "name=gnatsd", "--format", "{{.Names}}").CombinedOutput()
	if err != nil {
		return
	}

	s := string(out)
	if s == "" {
		out, err = exec.Command("docker", "restart", "postgres").CombinedOutput()
		if err != nil {
			log.Print(string(out))
			_, err = exec.Command("docker", "run", "-d", "--name", "gnatsd",
				"-p", "4222:4222",
				"-p", "6222:6222",
				"-p", "8222:8222",
				"nats:latest").CombinedOutput()

		}
		if err != nil {
			return
		}
	}

	return
}
