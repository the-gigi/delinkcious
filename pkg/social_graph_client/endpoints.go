package social_graph_client

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/go-kit/kit/endpoint"
	"net/http"
)

type SimpleResponse struct {
	Err error
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
	FollowEndpoint    endpoint.Endpoint
	UnfollowEndpoint endpoint.Endpoint
	GetFollowingEndpoint endpoint.Endpoint
	GetFollowersEndpoint endpoint.Endpoint
}

type FollowRequest struct {
	Followed string
	Follower string
}


func (s EndpointSet) Follow(followed string, follower string) error {
	resp, err := s.FollowEndpoint(context.Background(), FollowRequest{Followed: followed, Follower: follower})
	if err != nil {
		return err
	}
	response := resp.(SimpleResponse)
	return response.Err
}

type UnfollowRequest struct {
	Followed string
	Follower string
}

func (s EndpointSet) Unfollow(followed string, follower string) error {
	resp, err := s.UnfollowEndpoint(context.Background(), UnfollowRequest{Followed: followed, Follower: follower})
	if err != nil {
		return err
	}
	response := resp.(SimpleResponse)
	return response.Err
}

type GetFollowingRequest struct {
	Username string
}

type GetFollowingResponse struct {
	Following map[string]bool
	Err error
}

func decodeGetFollowingResponse(_ context.Context, r *http.Response) (interface{}, error) {
	if r.StatusCode != http.StatusOK {
		return nil, errors.New(r.Status)
	}
	var resp GetFollowingResponse
	err := json.NewDecoder(r.Body).Decode(&resp)
	return resp, err
}

func (s EndpointSet) GetFollowing(username string) (following map[string]bool, err error) {
	resp, err := s.GetFollowingEndpoint(context.Background(), GetFollowingRequest{Username: username})
	if err != nil {
		return
	}

	response := resp.(GetFollowingResponse)
	return response.Following, response.Err
}

type GetFollowersRequest struct {
	Username string
}

type GetFollowersResponse struct {
	Following map[string]bool
	Err error
}

func decodeGetFollowersResponse(_ context.Context, r *http.Response) (interface{}, error) {
	if r.StatusCode != http.StatusOK {
		return nil, errors.New(r.Status)
	}
	var resp GetFollowingResponse
	err := json.NewDecoder(r.Body).Decode(&resp)
	return resp, err
}


func (s EndpointSet) GetFollowers(username string) (following map[string]bool, err error) {
	resp, err := s.GetFollowersEndpoint(context.Background(), GetFollowersRequest{Username: username})
	if err != nil {
		return
	}

	response := resp.(GetFollowingResponse)
	return response.Following, response.Err
}

