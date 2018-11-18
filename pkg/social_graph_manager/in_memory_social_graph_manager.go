package social_graph_manager

import (
	om "../object_model"
	"errors"
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

type InMemorySocialGraphManager struct {
	socialGraph SocialGraph
}

func NewImMemorySocilaGraphManager() om.SocialGraphManager {
	return &InMemorySocialGraphManager{
		socialGraph: SocialGraph{},
	}
}

func (m *InMemorySocialGraphManager) Follow(followed string, follower string) (err error) {
	if followed == "" || follower == "" {
		err = errors.New("followed and follower can't be empty")
		return
	}

	followedUser := m.socialGraph[followed]
	if followedUser == nil {
		followedUser, _ = NewSocialUser(followed)
	}

	followedUser.Followers[follower] = true

	followerUser := m.socialGraph[follower]
	if followerUser == nil {
		followerUser, _ = NewSocialUser(follower)
	}

	followerUser.Following[followed] = true

	return
}

func (m *InMemorySocialGraphManager) Unfollow(followed string, follower string) (err error) {
	if followed == "" || follower == "" {
		err = errors.New("followed and follower can't be empty")
		return
	}

	followedUser := m.socialGraph[followed]
	if followedUser == nil {
		err = errors.New("followed user doesn't exist")
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

func (m *InMemorySocialGraphManager) AcceptFollowRequest(followed string, follower string) error {
	// All request are accepted automatically
	return nil
}

func (m *InMemorySocialGraphManager) RejectFollowRequest(followed string, follower string) error {
	// All request are accepted automatically
	return nil
}

func (m *InMemorySocialGraphManager) KickFollower(followed string, follower string) error {
	// No kicking allowed for in-memory social graph manager
	return nil
}

func (m *InMemorySocialGraphManager) GetFollowing(username string) map[string]bool {
	user := m.socialGraph[username]
	if user == nil {
		return map[string]bool{}
	}

	return user.Following
}

func (m *InMemorySocialGraphManager) GetFollowers(username string) map[string]bool {
	user := m.socialGraph[username]
	if user == nil {
		return map[string]bool{}
	}

	return user.Followers
}
