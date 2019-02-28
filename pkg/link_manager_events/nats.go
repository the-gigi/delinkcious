package link_manager_events

import "github.com/nats-io/go-nats"

func connect(url string) (encoddedConn *nats.EncodedConn, err error) {
	conn, err := nats.Connect(url)
	if err != nil {
		return
	}

	encoddedConn, err = nats.NewEncodedConn(conn, nats.JSON_ENCODER)
	return
}
