package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/nuclio/nuclio-sdk-go"
	"github.com/the-gigi/delinkcious/pkg/link_checker"
	"github.com/the-gigi/delinkcious/pkg/link_checker_events"
	om "github.com/the-gigi/delinkcious/pkg/object_model"
)

const natsUrl = "nats-cluster.default.svc.cluster.local:4222"

func Handler(context *nuclio.Context, event nuclio.Event) (interface{}, error) {
	r := nuclio.Response{
		StatusCode:  200,
		ContentType: "application/text",
	}

	body := event.GetBody()
	var e om.CheckLinkRequest
	err := json.Unmarshal(body, &e)
	if err != nil {
		msg := fmt.Sprintf("failed to unmarshal body: %v", body)
		context.Logger.Error(msg)

		r.StatusCode = 400
		r.Body = []byte(fmt.Sprintf(msg))
		return r, errors.New(msg)

	}

	username := e.Username
	url := e.Url
	if username == "" || url == "" {
		msg := fmt.Sprintf("missing USERNAME ('%s') and/or URL ('%s')", username, url)
		context.Logger.Error(msg)

		r.StatusCode = 400
		r.Body = []byte(msg)
		return r, errors.New(msg)
	}

	status := om.LinkStatusValid
	err = link_checker.CheckLink(url)
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
