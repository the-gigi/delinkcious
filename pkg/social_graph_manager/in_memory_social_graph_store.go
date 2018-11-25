package social_graph_manager

import (
	"errors"
	om "github.com/the-gigi/delinkcious/pkg/object_model"
)

type SocialGraph map[string]*SocialUser

type InMemorySocialGraphStore struct {
	socialGraph SocialGraph
}

func NewInMemorySocialGraphStore() om.SocialGraphManager {
	return &InMemorySocialGraphStore{
		socialGraph: SocialGraph{},
	}
}

func (m *InMemorySocialGraphStore) Follow(followed string, follower string) (err error) {
	followedUser := m.socialGraph[followed]
	if followedUser == nil {
		followedUser, _ = NewSocialUser(followed)
		m.socialGraph[followed] = followedUser
	}

	if followedUser.Followers[follower] {
		return errors.New("already following")
	}

	followedUser.Followers[follower] = true

	followerUser := m.socialGraph[follower]
	if followerUser == nil {
		followerUser, _ = NewSocialUser(follower)
		m.socialGraph[follower] = followerUser
	}

	followerUser.Following[followed] = true

	return
}

func (m *InMemorySocialGraphStore) Unfollow(followed string, follower string) (err error) {
	followedUser := m.socialGraph[followed]
	if followedUser == nil {
		err = errors.New("followed user doesn't exist")
		return
	}

	if !followedUser.Followers[follower] {
		err = errors.New("follower doesn't follow followed")
		return
	}

	followedUser.Followers[follower] = false

	followerUser := m.socialGraph[follower]
	if followerUser == nil {
		err = errors.New("follower user doesn't exist")
		return
	}

	followerUser.Following[followed] = false

	return
}

func (m *InMemorySocialGraphStore) AcceptFollowRequest(followed string, follower string) error {
	// All request are accepted automatically
	return nil
}

func (m *InMemorySocialGraphStore) RejectFollowRequest(followed string, follower string) error {
	// All request are accepted automatically
	return nil
}

func (m *InMemorySocialGraphStore) KickFollower(followed string, follower string) error {
	// No kicking allowed for in-memory social graph manager
	return nil
}

func (m *InMemorySocialGraphStore) GetFollowing(username string) map[string]bool {
	user := m.socialGraph[username]
	if user == nil {
		return map[string]bool{}
	}

	return user.Following
}

func (m *InMemorySocialGraphStore) GetFollowers(username string) map[string]bool {
	user := m.socialGraph[username]
	if user == nil {
		return map[string]bool{}
	}

	return user.Followers
}
