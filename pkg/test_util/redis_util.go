package test_util

import (
	"log"
	"os/exec"
)

func RunLocalRedisServer() (err error) {
	// Launch redis if not running
	out, err := exec.Command("docker", "ps", "-f", "name=redis", "--format", "{{.Names}}").CombinedOutput()
	if err != nil {
		return
	}

	s := string(out)
	if s == "" {
		out, err = exec.Command("docker", "restart", "redis").CombinedOutput()
		if err != nil {
			log.Print(string(out))
			_, err = exec.Command("docker", "run", "-d", "--name", "redis",
				"-p", "6379:6379",
				"redis:latest").CombinedOutput()

		}
		if err != nil {
			return
		}
	}

	return
}
