package social_graph_manager

import (
	"errors"
	om "github.com/the-gigi/delinkcious/pkg/object_model"
)

type Followers map[string]bool
type Following map[string]bool

type SocialUser struct {
	Username  string
	Followers Followers
	Following Following
}

func NewSocialUser(username string) (user *SocialUser, err error) {
	if username == "" {
		err = errors.New("user name can't be empty")
		return
	}

	user = &SocialUser{Username: username, Followers: Followers{}, Following: Following{}}
	return
}

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

func (m *InMemorySocialGraphStore) GetFollowing(username string) (map[string]bool, error) {
	user := m.socialGraph[username]
	if user == nil {
		return map[string]bool{}, nil
	}

	return user.Following, nil
}

func (m *InMemorySocialGraphStore) GetFollowers(username string) (map[string]bool, error) {
	user := m.socialGraph[username]
	if user == nil {
		return map[string]bool{}, nil
	}

	return user.Followers, nil
}
