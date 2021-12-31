package link_checker_events

import "github.com/nats-io/nats.go"

const (
	subject = "link-check-events"
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
