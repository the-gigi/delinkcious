package user_client

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/go-kit/kit/endpoint"
	om "github.com/the-gigi/delinkcious/pkg/object_model"
	"net/http"
)

type SimpleResponse struct {
	Err string
}

func decodeSimpleResponse(_ context.Context, r *http.Response) (interface{}, error) {
	if r.StatusCode != http.StatusOK {
		return nil, errors.New(r.Status)
	}
	var resp SimpleResponse
	err := json.NewDecoder(r.Body).Decode(&resp)
	return resp, err
}

type EndpointSet struct {
	RegisterEndpoint endpoint.Endpoint
	LoginEndpoint    endpoint.Endpoint
	LogoutEndpoint   endpoint.Endpoint
}

type registerRequest struct {
	Email    string
	Username string
}

func (s EndpointSet) Register(user om.User) (err error) {
	resp, err := s.RegisterEndpoint(context.Background(),
		registerRequest{Email: user.Email, Username: user.Name})
	if err != nil {
		return err
	}
	response := resp.(SimpleResponse)

	if response.Err != "" {
		err = errors.New(response.Err)
	}
	return
}

type loginRequest struct {
	Username  string
	AuthToken string
}

type loginResponse struct {
	Session string
	Err     string
}

func decodeLoginResponse(_ context.Context, r *http.Response) (interface{}, error) {
	if r.StatusCode != http.StatusOK {
		return nil, errors.New(r.Status)
	}
	var resp loginResponse
	err := json.NewDecoder(r.Body).Decode(&resp)
	return resp, err
}

func (s EndpointSet) Login(username string, authToken string) (session string, err error) {
	resp, err := s.LoginEndpoint(context.Background(),
		loginRequest{Username: username, AuthToken: authToken})
	if err != nil {
		return
	}
	response := resp.(loginResponse)
	if response.Err != "" {
		err = errors.New(response.Err)
		return
	}
	session = response.Session
	return
}

type logoutRequest struct {
	Username string
	Session  string
}

func (s EndpointSet) Logout(username string, session string) (err error) {
	resp, err := s.LogoutEndpoint(context.Background(),
		logoutRequest{Username: username, Session: session})
	if err != nil {
		return
	}

	response := resp.(SimpleResponse)
	if response.Err != "" {
		err = errors.New(response.Err)
	}
	return
}
