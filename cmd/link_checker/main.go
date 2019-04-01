package main

import (
	"errors"
	"fmt"
	"github.com/the-gigi/delinkcious/pkg/link_checker"
	"github.com/the-gigi/delinkcious/pkg/link_checker_events"
	om "github.com/the-gigi/delinkcious/pkg/object_model"
	"log"
	"os"
)

func main() {
	natsUrl := ""
	natsHostname := os.Getenv("NATS_CLUSTER_SERVICE_HOST")
	natsPort := os.Getenv("NATS_CLUSTER_SERVICE_PORT")
	if natsHostname != "" {
		natsUrl = fmt.Sprintf("%s:%s", natsHostname, natsPort)
	}

	username := os.Getenv("USERNAME")
	url := os.Getenv("URL")

	if username == "" || url == "" {
		log.Fatal(errors.New("missing environment variable USERNAME and/or URL"))
	}

	status := om.LinkStatusValid
	err := link_checker.CheckLink(url)
	if err != nil {
		status = om.LinkStatusInvalid
	}

	if natsUrl == "" {
		return
	}

	sender, err := link_checker_events.NewEventSender(natsUrl)
	if err != nil {
		log.Fatal(err)
	}

	sender.OnLinkChecked(username, url, status)
}
