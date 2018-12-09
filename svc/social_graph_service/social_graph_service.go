package main

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	om "github.com/the-gigi/delinkcious/pkg/object_model"
	sgm "github.com/the-gigi/delinkcious/pkg/social_graph_manager"
	"log"
	"net/http"
)

var (
	// return when an expected path variable is missing.
	BadRoutingError = errors.New("inconsistent mapping between route and handler")
)

type followRequest struct {
	Followed string
	Follower string
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
	Followers map[string]bool `json:"followers"`
	Err       string          `json:"err"`
}

type getFollowingRequest struct {
	Username string `json:"followed"`
}

type getFollowingResponse struct {
	Following map[string]bool `json:"following"`
	Err       string          `json:"err"`
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

// Extract the username from the request variables in the path
func getUsername(r *http.Request) (username string, err error) {
	vars := mux.Vars(r)
	username, ok := vars["username"]
	if !ok {
		err = BadRoutingError
	}
	return
}

func decodeGetFollowingRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	username, err := getUsername(r)
	return getFollowingRequest{Username: username}, err
}

func decodeGetFollowersRequest(_ context.Context, r *http.Request) (interface{}, error) {
	username, err := getUsername(r)
	return getFollowersRequest{Username: username}, err
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
		res := getFollowingResponse{Following: followingMap}
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
		res := getFollowersResponse{Followers: followersMap}
		if err != nil {
			res.Err = err.Error()
		}
		return res, nil
	}
}

func main() {
	store, err := sgm.NewDbSocialGraphStore("localhost", 5432, "postgres", "postgres")
	if err != nil {
		log.Fatal(err)
	}
	svc, err := sgm.NewSocialGraphManager(store)
	if err != nil {
		log.Fatal(err)
	}

	followHandler := httptransport.NewServer(
		makeFollowEndpoint(svc),
		decodeFollowRequest,
		encodeResponse,
	)

	unfollowHandler := httptransport.NewServer(
		makeUnfollowEndpoint(svc),
		decodeUnfollowRequest,
		encodeResponse,
	)

	getFollowingHandler := httptransport.NewServer(
		makeGetFollowingEndpoint(svc),
		decodeGetFollowingRequest,
		encodeResponse,
	)

	getFollowersHandler := httptransport.NewServer(
		makeGetFollowersEndpoint(svc),
		decodeGetFollowersRequest,
		encodeResponse,
	)

	r := mux.NewRouter()
	r.Methods("POST").Path("/follow").Handler(followHandler)
	r.Methods("POST").Path("/unfollow").Handler(unfollowHandler)
	r.Methods("GET").Path("/following/{username}").Handler(getFollowingHandler)
	r.Methods("GET").Path("/followers/{username}").Handler(getFollowersHandler)

	log.Fatal(http.ListenAndServe(":9090", r))
}
