package news_manager

import (
	om "github.com/the-gigi/delinkcious/pkg/object_model"
)

type Store interface {
	GetNews(username string, startIndex int) (events []*om.Event, nextIndex int, err error)
	AddEvent(username string, event *om.Event) (err error)
}
