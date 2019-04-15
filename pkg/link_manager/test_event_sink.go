package link_manager

import (
	om "github.com/the-gigi/delinkcious/pkg/object_model"
)

type linkManagerEventsSink struct {
	addLinkEvents     map[string][]*om.Link
	updateLinkEvents  map[string][]*om.Link
	deletedLinkEvents map[string][]string
}

func (s *linkManagerEventsSink) OnLinkAdded(username string, link *om.Link) {
	if s.addLinkEvents[username] == nil {
		s.addLinkEvents[username] = []*om.Link{}
	}
	s.addLinkEvents[username] = append(s.addLinkEvents[username], link)
}

func (s *linkManagerEventsSink) OnLinkUpdated(username string, link *om.Link) {
	if s.updateLinkEvents[username] == nil {
		s.updateLinkEvents[username] = []*om.Link{}
	}
	s.updateLinkEvents[username] = append(s.updateLinkEvents[username], link)
}

func (s *linkManagerEventsSink) OnLinkDeleted(username string, url string) {
	if s.deletedLinkEvents[username] == nil {
		s.deletedLinkEvents[username] = []string{}
	}
	s.deletedLinkEvents[username] = append(s.deletedLinkEvents[username], url)
}

func newLinkManagerEventsSink() *linkManagerEventsSink {
	return &linkManagerEventsSink{
		map[string][]*om.Link{},
		map[string][]*om.Link{},
		map[string][]string{},
	}
}
