package main

import (
	"log"
	"net/http"

	httptransport "github.com/go-kit/kit/transport/http"
	sgm "github.com/the-gigi/delinkcious/pkg/social_graph_manager"
)

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

	http.Handle("/follow", followHandler)
	http.Handle("/unfollow", unfollowHandler)
	http.Handle("/following", getFollowingHandler)
	http.Handle("/followers", getFollowersHandler)

	log.Fatal(http.ListenAndServe(":9090", nil))
}
