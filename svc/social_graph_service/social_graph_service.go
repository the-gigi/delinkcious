package main

import (
	"log"
	"net/http"

	httptransport "github.com/go-kit/kit/transport/http"
	sgm "github.com/the-gigi/delinkcious/pkg/social_graph_manager"
)

func main() {
	store, err := sgm.NewDbSocialGraphStore("localhost", 5432, "posgres", "postgres")
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

	http.Handle("/follow", followHandler)
	http.Handle("/unfollow", unfollowHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
