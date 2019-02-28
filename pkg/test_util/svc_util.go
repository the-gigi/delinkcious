package test_util

import (
	"context"
	"os"
	"os/exec"
)

func Check(err error) {
	if err != nil {
		panic(err)
	}
}

// Build and run a service in a target directory
func RunService(ctx context.Context, targetDir string, service string) {
	// Save and restore later current working dir
	wd, err := os.Getwd()
	Check(err)
	defer os.Chdir(wd)

	// Build the server if needed
	os.Chdir(targetDir)
	_, err = os.Stat("./" + service)
	if os.IsNotExist(err) {
		_, err := exec.Command("go", "build", ".").CombinedOutput()
		Check(err)
	}

	cmd := exec.CommandContext(ctx, "./"+service)
	err = cmd.Start()
	Check(err)
}

func KillServer(ctx context.Context) {
	ctx.Done()
}
