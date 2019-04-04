package main

import (
	"errors"
	"fmt"
	"github.com/nuclio/nuclio-sdk-go"
	"github.com/the-gigi/delinkcious/pkg/link_checker"
	"github.com/the-gigi/delinkcious/pkg/link_checker_events"
	om "github.com/the-gigi/delinkcious/pkg/object_model"
	"os"
)

//func Handler(context *nuclio.Context, event nuclio.Event) (interface{}, error) {
//	context.Logger.Info("This is an unstrucured %s", "log")
//
//	return nuclio.Response{
//		StatusCode:  200,
//		ContentType: "application/text",
//		Body:        []byte("Hello, from nuclio :]"),
//	}, nil
//}

func Handler(context *nuclio.Context, event nuclio.Event) (interface{}, error) {
	natsUrl := ""
	natsHostname := os.Getenv("NATS_CLUSTER_SERVICE_HOST")
	natsPort := os.Getenv("NATS_CLUSTER_SERVICE_PORT")
	if natsHostname != "" {
		natsUrl = fmt.Sprintf("%s:%s", natsHostname, natsPort)
	}

	r := nuclio.Response{
		StatusCode:  200,
		ContentType: "application/text",
	}

	if natsUrl == "" {
		r.Body = []byte("No NATS url. Skipping link check")
		return r, nil
	}

	username := os.Getenv("USERNAME")
	url := os.Getenv("URL")

	if username == "" || url == "" {
		msg := fmt.Sprintf("missing USERNAME ('%s') and/or URL ('%s')", username, url)
		context.Logger.Error(msg)

		r.StatusCode = 500
		r.Body = []byte(fmt.Sprintf(msg, username, url))
		return r, errors.New(msg)
	}

	status := om.LinkStatusValid
	err := link_checker.CheckLink(url)
	if err != nil {
		status = om.LinkStatusInvalid
	}

	sender, err := link_checker_events.NewEventSender(natsUrl)
	if err != nil {
		context.Logger.Error(err.Error())

		r.StatusCode = 500
		r.Body = []byte(err.Error())
		return r, err
	}

	sender.OnLinkChecked(username, url, status)
	return r, nil
}

func main() {

}
