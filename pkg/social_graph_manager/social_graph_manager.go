package social_graph_manager

import (
	"errors"
	om "github.com/the-gigi/delinkcious/pkg/object_model"
)

type SocialGraphManager struct {
	store om.SocialGraphManager
}

func NewSocialGraphManager(store om.SocialGraphManager) (om.SocialGraphManager, error) {
	if store == nil {
		return nil, errors.New("store can't be nil")
	}
	return &SocialGraphManager{store: store}, nil
}

func (m *SocialGraphManager) Follow(followed string, follower string) (err error) {
	if followed == "" || follower == "" {
		err = errors.New("followed and follower can't be empty")
		return
	}

	return m.store.Follow(followed, follower)
}

func (m *SocialGraphManager) Unfollow(followed string, follower string) (err error) {
	if followed == "" || follower == "" {
		err = errors.New("followed and follower can't be empty")
		return
	}

	return m.store.Unfollow(followed, follower)
}

func (m *SocialGraphManager) AcceptFollowRequest(followed string, follower string) (err error) {
	if followed == "" || follower == "" {
		err = errors.New("followed and follower can't be empty")
		return
	}

	// All request are accepted automatically
	return nil
}

func (m *SocialGraphManager) RejectFollowRequest(followed string, follower string) (err error) {
	if followed == "" || follower == "" {
		err = errors.New("followed and follower can't be empty")
		return
	}

	// All request are accepted automatically
	return nil
}

func (m *SocialGraphManager) KickFollower(followed string, follower string) (err error) {
	if followed == "" || follower == "" {
		err = errors.New("followed and follower can't be empty")
		return
	}
	// No kicking implemented
	return nil
}

func (m *SocialGraphManager) GetFollowing(username string) (map[string]bool, error) {
	return m.store.GetFollowing(username)
}

func (m *SocialGraphManager) GetFollowers(username string) (map[string]bool, error) {
	return m.store.GetFollowers(username)
}
