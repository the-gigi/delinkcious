package link_manager_events

import (
	om "github.com/the-gigi/delinkcious/pkg/object_model"
)

func Listen(url string, sink om.LinkManagerEvents) (err error) {
	conn, err := connect(url)
	if err != nil {
		return
	}

	conn.QueueSubscribe(subject, queue, func(e *Event) {
		switch e.EventType {
		case om.LinkAdded:
			{
				sink.OnLinkAdded(e.Username, e.Link)
			}
		case om.LinkUpdated:
			{
				sink.OnLinkUpdated(e.Username, e.Link)
			}
		default:
			// Ignore other event types
		}
	})

	return
}
