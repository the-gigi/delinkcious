package service

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/go-kit/kit/endpoint"
	"github.com/the-gigi/delinkcious/pkg/auth_util"
	om "github.com/the-gigi/delinkcious/pkg/object_model"
	"net/http"
	"net/url"
	"os"
	"strings"
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

type getByUsernameRequest struct {
	Username string `json:"username"`
}

type getFollowersResponse struct {
	Followers map[string]bool `json:"followers"`
	Err       string          `json:"err"`
}

type getFollowingResponse struct {
	Following map[string]bool `json:"following"`
	Err       string          `json:"err"`
}

func newGetByUsernameRequest(username string) (request getByUsernameRequest, err error) {
	request.Username, err = url.PathUnescape(username)
	return
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
	parts := strings.Split(r.URL.Path, "/")
	username := parts[len(parts)-1]
	if username == "" || username == "following" {
		return nil, errors.New("user name must not be empty")
	}
	return newGetByUsernameRequest(username)
}

func decodeGetFollowersRequest(_ context.Context, r *http.Request) (interface{}, error) {
	if os.Getenv("DELINKCIOUS_MUTUAL_AUTH") != "false" {
		token := r.Header["Delinkcious-Caller-Token"]
		if len(token) == 0 || token[0] == "" {
			return nil, errors.New("missing caller token")
		}

		if !auth_util.HasCaller("link-manager", token[0]) {
			return nil, errors.New("unauthorized caller")
		}
	}
	parts := strings.Split(r.URL.Path, "/")
	username := parts[len(parts)-1]
	if username == "" || username == "followers" {
		return nil, errors.New("user name must not be empty")
	}

	return newGetByUsernameRequest(username)
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

//func wasteCPU() {
//	fmt.Println("wasteCPU() here!")
//	go func() {
//		for {
//			if rand.Int() % 8000 == 0 {
//				time.Sleep(50 * time.Microsecond)
//			}
//		}
//	}()
//}


func makeGetFollowingEndpoint(svc om.SocialGraphManager) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(getByUsernameRequest)
		followingMap, err := svc.GetFollowing(req.Username)
		res := getFollowingResponse{Following: followingMap}
		if err != nil {
			res.Err = err.Error()
		}
		return res, nil
	}
}

func makeGetFollowersEndpoint(svc om.SocialGraphManager) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		//wasteCPU()
		req := request.(getByUsernameRequest)
		followersMap, err := svc.GetFollowers(req.Username)
		res := getFollowersResponse{Followers: followersMap}
		if err != nil {
			res.Err = err.Error()
		}
		return res, nil
	}
}
