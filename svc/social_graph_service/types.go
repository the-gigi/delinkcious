package main

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

type GetFollowingResponse struct {
	Err string `json:"err"`
}
