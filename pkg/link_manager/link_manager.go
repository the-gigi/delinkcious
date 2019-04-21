package link_manager

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/the-gigi/delinkcious/pkg/link_checker_events"
	om "github.com/the-gigi/delinkcious/pkg/object_model"
	"net/http"
)

// Nuclio functions listen by default on port 8080 of their service IP
const link_checker_func_url = "http://link-checker.nuclio.svc.cluster.local:8080"

type LinkManager struct {
	linkStore          LinkStore
	socialGraphManager om.SocialGraphManager
	eventSink          om.LinkManagerEvents
	maxLinksPerUser    int64
}

func (m *LinkManager) GetLinks(request om.GetLinksRequest) (result om.GetLinksResult, err error) {
	if request.Username == "" {
		err = errors.New("user name can't be empty")
		return
	}

	result, err = m.linkStore.GetLinks(request)
	if result.Links == nil {
		result.Links = []om.Link{}
	}
	return
}

// Very wasteful way to count links
func (m *LinkManager) getLinkCount(username string) (linkCount int64, err error) {
	req := om.GetLinksRequest{Username: username}
	res, err := m.GetLinks(req)
	if err != nil {
		return
	}

	linkCount += int64(len(res.Links))

	for res.NextPageToken != "" {
		req = om.GetLinksRequest{Username: username, StartToken: res.NextPageToken}
		res, err = m.GetLinks(req)
		if err != nil {
			return
		}

		linkCount += int64(len(res.Links))
	}
	return
}

func triggerLinkCheck(username string, url string) {
	go func() {
		checkLinkRequest := &om.CheckLinkRequest{Username: username, Url: url}
		data, err := json.Marshal(checkLinkRequest)
		if err != nil {
			return
		}

		req, err := http.NewRequest("POST", link_checker_func_url, bytes.NewBuffer(data))
		req.Header.Set("Content-Type", "application/json")
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			return
		}
		defer resp.Body.Close()
	}()
}

func (m *LinkManager) AddLink(request om.AddLinkRequest) (err error) {
	if request.Url == "" {
		return errors.New("the URL can't be empty")
	}

	if request.Username == "" {
		return errors.New("the user name can't be empty")
	}

	linkCount, err := m.getLinkCount(request.Username)
	if err != nil {
		return
	}

	if linkCount >= m.maxLinksPerUser {
		return errors.New("the user has too many links")
	}

	link, err := m.linkStore.AddLink(request)
	if err != nil {
		return
	}

	if m.eventSink != nil {
		followers, err := m.socialGraphManager.GetFollowers(request.Username)
		if err != nil {
			return err
		}

		for follower := range followers {
			m.eventSink.OnLinkAdded(follower, link)
		}
	}

	// Trigger link check asynchronously (don't wait for result)
	triggerLinkCheck(request.Username, request.Url)
	return
}

func (m *LinkManager) UpdateLink(request om.UpdateLinkRequest) (err error) {
	if request.Url == "" {
		return errors.New("the URL can't be empty")
	}

	if request.Username == "" {
		return errors.New("the user name can't be empty")
	}

	link, err := m.linkStore.UpdateLink(request)
	if err != nil {
		return
	}

	if m.eventSink != nil {
		followers, err := m.socialGraphManager.GetFollowers(request.Username)
		if err != nil {
			return err
		}

		for follower := range followers {
			m.eventSink.OnLinkUpdated(follower, link)
		}
	}

	return
}

func (m *LinkManager) DeleteLink(username string, url string) (err error) {
	if url == "" {
		return errors.New("the URL can't be empty")
	}

	if username == "" {
		return errors.New("the user name can't be empty")
	}

	err = m.linkStore.DeleteLink(username, url)
	if err != nil {
		return
	}

	if m.eventSink != nil {
		followers, err := m.socialGraphManager.GetFollowers(username)
		if err != nil {
			return err
		}

		for follower := range followers {
			m.eventSink.OnLinkDeleted(follower, url)
		}
	}
	return
}

func (m *LinkManager) OnLinkChecked(username string, url string, status om.LinkStatus) {
	m.linkStore.SetLinkStatus(username, url, status)
}

func NewLinkManager(linkStore LinkStore,
	socialGraphManager om.SocialGraphManager,
	natsUrl string,
	eventSink om.LinkManagerEvents,
	maxLinksPerUser int64) (om.LinkManager, error) {
	if linkStore == nil {
		return nil, errors.New("link store")
	}

	if eventSink != nil && socialGraphManager == nil {
		return nil, errors.New("social graph manager can't be nil if event sink is not nil")
	}

	link_manager := &LinkManager{
		linkStore:          linkStore,
		socialGraphManager: socialGraphManager,
		eventSink:          eventSink,
		maxLinksPerUser:    maxLinksPerUser,
	}

	// Subscribe LinkManager to link checker events if nats is ocnfigured
	if natsUrl != "" {
		link_checker_events.Listen(natsUrl, link_manager)
	}

	return link_manager, nil
}
