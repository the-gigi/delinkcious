package service

import (
	"context"
	"encoding/json"
	"github.com/go-kit/kit/endpoint"
	om "github.com/the-gigi/delinkcious/pkg/object_model"
	"net/http"
)

type registerRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
}

type registerResponse struct {
	Err string `json:"err"`
}

type loginRequest struct {
	Username  string `json:"username"`
	AuthToken string `json:"authToken"`
}

type loginResponse struct {
	Session string `json:"session"`
	Err     string `json:"err"`
}

type logoutRequest struct {
	Username string `json:"username"`
	Session  string `json:"session"`
}

type logoutResponse struct {
	Err string `json:"err"`
}

func decodeRegisterRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request registerRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		return nil, err
	}
	return request, nil
}

func decodeLoginRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request loginRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		return nil, err
	}
	return request, nil
}

func decodeLogoutRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request logoutRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		return nil, err
	}
	return request, nil
}

func encodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}

func makeRegisterEndpoint(svc om.UserManager) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(registerRequest)
		err := svc.Register(om.User{req.Email, req.Username})
		res := registerResponse{}
		if err != nil {
			res.Err = err.Error()
		}
		return res, nil
	}
}

func makeLoginEndpoint(svc om.UserManager) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(loginRequest)
		session, err := svc.Login(req.Username, req.AuthToken)
		res := loginResponse{Session: session}
		if err != nil {
			res.Err = err.Error()
		}
		return res, nil
	}
}

func makeLogoutEndpoint(svc om.UserManager) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(logoutRequest)
		err := svc.Logout(req.Username, req.Session)
		res := logoutResponse{}
		if err != nil {
			res.Err = err.Error()
		}
		return res, nil
	}
}
