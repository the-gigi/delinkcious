package link_checker_events

import (
	"github.com/nats-io/go-nats"
	om "github.com/the-gigi/delinkcious/pkg/object_model"
)

const subject = "link-check-events"
const queue = "the-queue"

func connect(url string) (encodedConn *nats.EncodedConn, err error) {
	conn, err := nats.Connect(url)
	if err != nil {
		return
	}

	encodedConn, err = nats.NewEncodedConn(conn, nats.JSON_ENCODER)
	return
}

func Listen(url string, sink om.LinkCheckerEvents) (err error) {
	conn, err := connect(url)
	if err != nil {
		return
	}

	conn.QueueSubscribe(subject, queue, func(e *Event) {
		sink.OnLinkChecked(e.Username, e.Url, e.Status)
	})

	return
}
