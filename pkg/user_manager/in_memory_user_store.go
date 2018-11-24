package user_manager

import (
	"errors"
	om "github.com/the-gigi/delinkcious/pkg/object_model"
	"math/rand"
	"strconv"
)

type InMemoryUserStore struct {
	registeredUsers map[string]bool
	loggedInUsers   map[string]bool
	sessions        map[string]string
}

func NewInMemoryUserStore() om.UserManager {
	return &InMemoryUserStore{
		registeredUsers: map[string]bool{},
		loggedInUsers:   map[string]bool{},
		sessions:        map[string]string{},
	}
}

func (m *InMemoryUserStore) Register(user om.User) (err error) {
	if m.registeredUsers[user.Name] {
		return errors.New("user already registered")
	}

	m.registeredUsers[user.Name] = true
	return
}

func (m *InMemoryUserStore) Login(username string, authToken string) (session string, err error) {
	if !m.registeredUsers[username] {
		err = errors.New("user not registered")
		return
	}

	if m.loggedInUsers[username] {
		err = errors.New("user already logged in")
		return
	}

	m.loggedInUsers[username] = true
	session = strconv.Itoa(rand.Int())
	m.sessions[session] = username
	return
}

func (m *InMemoryUserStore) Logout(username string, session string) (err error) {
	if !m.loggedInUsers[username] {
		err = errors.New("User is not logged in")
		return
	}

	if m.sessions[session] != username {
		err = errors.New("Invalid session")
		return
	}

	delete(m.sessions, session)
	delete(m.loggedInUsers, username)

	return
}
