package user_manager

import (
	om "github.com/the-gigi/delinkcious/pkg/object_model"
	"errors"
	"math/rand"
	"strconv"
)

type InMemoryUserManager struct {
	registeredUsers map[string]bool
	loggedInUsers   map[string]bool
	sessions        map[string]string
}

func NewImMemoryUserManager() om.UserManager {
	return &InMemoryUserManager{
		registeredUsers: map[string]bool{},
		loggedInUsers:   map[string]bool{},
		sessions:        map[string]string{},
	}
}

func (m *InMemoryUserManager) Register(user om.User) error {
	if user.Name == "" {
		return errors.New("invalid user name")
	}
	if m.registeredUsers[user.Name] {
		return errors.New("user already registered")
	}

	m.registeredUsers[user.Name] = true
	return nil
}

func (m *InMemoryUserManager) Login(username string, authToken string) (session string, err error) {
	if !m.registeredUsers[username] {
		return "", errors.New("user not registered")
	}

	if m.loggedInUsers[username] {
		return "", errors.New("user already logged in")
	}

	m.loggedInUsers[username] = true
	session = strconv.Itoa(rand.Int())
	m.sessions[session] = username

	return
}

func (m *InMemoryUserManager) Logout(username string, session string) error {
	if !m.loggedInUsers[username] {
		return errors.New("User is not logged in")
	}

	if m.sessions[session] != username {
		return errors.New("Invalid session")
	}

	delete(m.sessions, session)
	delete(m.loggedInUsers, username)

	return nil
}
