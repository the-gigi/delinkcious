package link_manager

import (
	"errors"
	om "github.com/the-gigi/delinkcious/pkg/object_model"
)

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

func NewLinkManager(linkStore LinkStore,
	socialGraphManager om.SocialGraphManager,
	eventSink om.LinkManagerEvents,
	maxLinksPerUser int64) (om.LinkManager, error) {
	if linkStore == nil {
		return nil, errors.New("link store")
	}

	if eventSink != nil && socialGraphManager == nil {
		return nil, errors.New("social graph manager can't be nil if event sink is not nil")
	}

	return &LinkManager{
		linkStore:          linkStore,
		socialGraphManager: socialGraphManager,
		eventSink:          eventSink,
		maxLinksPerUser:    maxLinksPerUser,
	}, nil
}
