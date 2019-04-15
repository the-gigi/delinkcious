package link_manager

type mockSocialGraphManager struct {
	followers map[string]bool
}

func (m *mockSocialGraphManager) Follow(followed string, follower string) error {
	return nil
}

func (m *mockSocialGraphManager) Unfollow(followed string, follower string) error {
	return nil
}

func (m *mockSocialGraphManager) GetFollowing(username string) (map[string]bool, error) {
	return nil, nil
}

func (m *mockSocialGraphManager) GetFollowers(username string) (map[string]bool, error) {
	return m.followers, nil
}

func newMockSocialGraphManager(followers []string) *mockSocialGraphManager {
	m := &mockSocialGraphManager{
		map[string]bool{},
	}
	for _, f := range followers {
		m.followers[f] = true
	}

	return m
}
