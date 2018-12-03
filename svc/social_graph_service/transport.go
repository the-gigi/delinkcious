package main

import (
	"context"
	"encoding/json"
	"github.com/go-kit/kit/endpoint"
	om "github.com/the-gigi/delinkcious/pkg/object_model"
	"net/http"
)

type followRequest struct {
	Followed string `json:"followed"`
	Follower string `json:"follower"`
}

type followResponse struct {
	Err string `json:"err"`
}

type unfollowRequest struct {
	Followed string `json:"followed"`
	Follower string `json:"follower"`
}

type unfollowResponse struct {
	Err string `json:"err"`
}

type getFollowersRequest struct {
	Username string `json:"followed"`
}

type getFollowersResponse struct {
	Followers []string `json:"followers"`
	Err       string   `json:"err"`
}

type getFollowingRequest struct {
	Username string `json:"followed"`
}

type getFollowingResponse struct {
	Following []string `json:"following"`
	Err       string   `json:"err"`

}

func decodeFollowRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request followRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		return nil, err
	}
	return request, nil
}

func decodeUnfollowRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request unfollowRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		return nil, err
	}
	return request, nil
}

func decodeGetFollowingRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request getFollowersRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		return nil, err
	}
	return request, nil
}

func decodeGetFollowersRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request getFollowersRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		return nil, err
	}
	return request, nil
}


func encodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}

func makeFollowEndpoint(svc om.SocialGraphManager) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(followRequest)
		err := svc.Follow(req.Followed, req.Follower)
		res := followResponse{}
		if err != nil {
			res.Err = err.Error()
		}
		return res, nil
	}
}

func makeUnfollowEndpoint(svc om.SocialGraphManager) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(unfollowRequest)
		err := svc.Unfollow(req.Followed, req.Follower)
		res := unfollowResponse{}
		if err != nil {
			res.Err = err.Error()
		}
		return res, nil
	}
}

func makeGetFollowingEndpoint(svc om.SocialGraphManager) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(getFollowingRequest)
		followingMap, err := svc.GetFollowing(req.Username)

		following := []string{}
		for f, _ := range followingMap {
			following = append(following, f)
		}

		res := getFollowingResponse{Following: following}
		if err != nil {
			res.Err = err.Error()
		}
		return res, nil
	}
}

func makeGetFollowersEndpoint(svc om.SocialGraphManager) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(getFollowersRequest)
		followersMap, err := svc.GetFollowers(req.Username)

		followers := []string{}
		for f, _ := range followersMap {
			followers = append(followers, f)
		}

		res := getFollowersResponse{Followers: followers}
		if err != nil {
			res.Err = err.Error()
		}
		return res, nil
	}
}
