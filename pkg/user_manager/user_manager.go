package user_manager

import (
	"errors"
	om "github.com/the-gigi/delinkcious/pkg/object_model"
)

type UserManager struct {
	userStore om.UserManager
}

func NewUserManager(userStore om.UserManager) (userManager om.UserManager, err error) {
	if userStore == nil {
		return nil, errors.New("user store can't be nil")
	}
	return &UserManager{userStore: userStore}, nil
}

func (m *UserManager) Register(user om.User) error {
	if user.Name == "" {
		return errors.New("invalid user name")
	}

	return m.userStore.Register(user)
}

func (m *UserManager) Login(username string, authToken string) (session string, err error) {
	if username == "" {
		return "", errors.New("username can't be empty")
	}

	return m.userStore.Login(username, authToken)
}

func (m *UserManager) Logout(username string, session string) error {
	return m.userStore.Logout(username, session)
}
