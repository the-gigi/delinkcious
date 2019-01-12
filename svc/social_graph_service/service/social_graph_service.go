package service

import (
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	sgm "github.com/the-gigi/delinkcious/pkg/social_graph_manager"
	"log"
	"net/http"
	"os"
	"strconv"
)

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func Run() {
	log.Println("Service started...")
	dbHost := os.Getenv("SOCIAL_GRAPH_DB_SERVICE_HOST")
	if dbHost == "" {
		dbHost = "localhost"
	}

	port := os.Getenv("SOCIAL_GRAPH_DB_SERVICE_PORT")
	if port == "" {
		port = "5432"
	}

	log.Println("DB host:", dbHost, "DB port:", port)

	dbPort, err := strconv.Atoi(port)
	check(err)

	store, err := sgm.NewDbSocialGraphStore(dbHost, dbPort, "postgres", "postgres")
	check(err)

	svc, err := sgm.NewSocialGraphManager(store)
	check(err)

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

	log.Println("Listening on port 9090...")
	log.Fatal(http.ListenAndServe(":9090", r))
}
