package news_manager

import (
	"errors"
	om "github.com/the-gigi/delinkcious/pkg/object_model"
)

const maxPageSize = 10

// User links are a map of url:TaggedLink
type UserEvents map[string][]*om.Event

// Link store is a map of username:userEvents
type InMemoryNewsStore struct {
	userEvents UserEvents
}

func (m *InMemoryNewsStore) GetNews(username string, startIndex int) (events []*om.Event, err error) {
	userEvents := m.userEvents[username]
	if startIndex > len(userEvents) {
		err = errors.New("Index out of bounds")
		return
	}

	pageSize := len(userEvents) - startIndex
	if pageSize > maxPageSize {
		pageSize = maxPageSize
	}

	events = userEvents[startIndex : startIndex+pageSize]
	return
}

func (m *InMemoryNewsStore) AddEvent(username string, event *om.Event) (err error) {
	if username == "" {
		err = errors.New("user name can't be empty")
		return
	}

	if event == nil {
		err = errors.New("event can't be nil")
		return
	}

	if m.userEvents[username] == nil {
		m.userEvents[username] = []*om.Event{}
	}

	m.userEvents[username] = append(m.userEvents[username], event)
	return
}

func NewInMemoryNewsStore() *InMemoryNewsStore {
	return &InMemoryNewsStore{UserEvents{}}
}
