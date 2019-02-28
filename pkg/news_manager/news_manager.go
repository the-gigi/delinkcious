package news_manager

import (
	"errors"
	"github.com/the-gigi/delinkcious/pkg/link_manager_events"
	om "github.com/the-gigi/delinkcious/pkg/object_model"
	"strconv"
	"time"
)

type NewsManager struct {
	eventStore *InMemoryNewsStore
}

func (m *NewsManager) GetNews(req om.GetNewsRequest) (resp om.GetNewsResult, err error) {
	if req.Username == "" {
		err = errors.New("user name can't be empty")
		return
	}

	startIndex := 0
	if req.StartToken != "" {
		startIndex, err := strconv.Atoi(req.StartToken)
		if err != nil || startIndex < 0 {
			err = errors.New("invalid start token: " + req.StartToken)
			return resp, err
		}
	}

	events, nextIndex, err := m.eventStore.GetNews(req.Username, startIndex)
	if err != nil {
		return
	}

	resp.Events = events
	if nextIndex != -1 {
		resp.NextToken = strconv.Itoa(nextIndex)
	}

	return
}

func (m *NewsManager) OnLinkAdded(username string, link *om.Link) {
	event := &om.Event{
		EventType: om.LinkAdded,
		Username:  username,
		Url:       link.Url,
		Timestamp: time.Now().UTC(),
	}
	m.eventStore.AddEvent(username, event)
}

func (m *NewsManager) OnLinkUpdated(username string, link *om.Link) {
	event := &om.Event{
		EventType: om.LinkUpdated,
		Username:  username,
		Url:       link.Url,
		Timestamp: time.Now().UTC(),
	}
	m.eventStore.AddEvent(username, event)
}

func (m *NewsManager) OnLinkDeleted(username string, url string) {
	event := &om.Event{
		EventType: om.LinkDeleted,
		Username:  username,
		Url:       url,
		Timestamp: time.Now().UTC(),
	}
	m.eventStore.AddEvent(username, event)
}

func NewNewsManager(natsHostname string, natsPort string) (om.NewsManager, error) {
	nm := &NewsManager{eventStore: NewInMemoryNewsStore()}
	if natsHostname != "" {
		natsUrl := natsHostname + ":" + natsPort
		err := link_manager_events.Listen(natsUrl, nm)
		if err != nil {
			return nil, err
		}
	}

	return nm, nil
}
