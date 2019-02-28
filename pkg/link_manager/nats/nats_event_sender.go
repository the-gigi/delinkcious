package nats

import (
	"github.com/nats-io/go-nats"
	om "github.com/the-gigi/delinkcious/pkg/object_model"
	"log"
)

type event struct {
	EventType om.EventTypeEnum
	Link      om.Link
}

type eventSender struct {
	hostname string
	nats     *nats.EncodedConn
}

func connect(url string) (encoddedConn *nats.EncodedConn, err error) {
	conn, err := nats.Connect(url)
	if err != nil {
		return
	}

	encoddedConn, err = nats.NewEncodedConn(conn, nats.JSON_ENCODER)
	return
}

func (s *eventSender) OnLinkAdded(username string, link *om.Link) {
	err := s.nats.Publish("link-events", event{om.LinkAdded, *link})
	if err != nil {
		log.Fatal(err)
	}
}

func (s *eventSender) OnLinkUpdated(username string, link *om.Link) {
	err := s.nats.Publish("link-events", event{om.LinkUpdated, *link})
	if err != nil {
		log.Fatal(err)
	}
}

func (s *eventSender) OnLinkDeleted(username string, url string) {
	// Ignore link delete events
}

func NewEventSender(url string) (om.LinkManagerEvents, error) {
	ec, err := connect(url)
	if err != nil {
		return nil, err
	}
	return &eventSender{hostname: url, nats: ec}, nil
}
