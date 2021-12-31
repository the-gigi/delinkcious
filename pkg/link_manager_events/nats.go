package link_manager_events

import (
	"github.com/nats-io/nats.go"
)

const (
	subject = "link-events"
	queue   = "the-queue"
)

func connect(url string) (encodedConn *nats.EncodedConn, err error) {
	conn, err := nats.Connect(url)
	if err != nil {
		return
	}

	encodedConn, err = nats.NewEncodedConn(conn, nats.JSON_ENCODER)
	return
}
