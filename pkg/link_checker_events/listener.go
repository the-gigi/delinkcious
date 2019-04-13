package link_checker_events

import (
	om "github.com/the-gigi/delinkcious/pkg/object_model"
)

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
