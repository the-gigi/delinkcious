package news_manager

import (
	"errors"
	om "github.com/the-gigi/delinkcious/pkg/object_model"
)

const maxPageSize = 10

// User events are a map of username:userEvents
type userEvents map[string][]*om.Event

// InMemoryNewsStore manages a UserEvents data structure
type InMemoryNewsStore struct {
	userEvents userEvents
}

func (m *InMemoryNewsStore) GetNews(username string, startIndex int) (events []*om.Event, nextIndex int, err error) {
	userEvents := m.userEvents[username]
	if startIndex > len(userEvents) {
		err = errors.New("Index out of bounds")
		return
	}

	pageSize := len(userEvents) - startIndex
	if pageSize > maxPageSize {
		pageSize = maxPageSize
		nextIndex = startIndex + maxPageSize
	} else {
		nextIndex = -1
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
	return &InMemoryNewsStore{userEvents{}}
}
