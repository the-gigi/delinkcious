package service

import (
	"github.com/gorilla/mux"
	"github.com/the-gigi/delinkcious/pkg/db_util"
	"log"
	"net/http"

	httptransport "github.com/go-kit/kit/transport/http"
	lm "github.com/the-gigi/delinkcious/pkg/link_manager"
	sgm "github.com/the-gigi/delinkcious/pkg/social_graph_client"
)

func Run() {
	dbHost, dbPort, err := db_util.GetDbEndpoint("social_graph")
	if err != nil {
		log.Fatal(err)
	}

	store, err := lm.NewDbLinkStore(dbHost, dbPort, "postgres", "postgres")
	if err != nil {
		log.Fatal(err)
	}

	socialGraphClient, err := sgm.NewClient("localhost:9090")
	if err != nil {
		log.Fatal(err)
	}

	svc, err := lm.NewLinkManager(store, socialGraphClient, nil)
	if err != nil {
		log.Fatal(err)
	}

	getLinksHandler := httptransport.NewServer(
		makeGetLinksEndpoint(svc),
		decodeGetLinksRequest,
		encodeResponse,
	)

	addLinkHandler := httptransport.NewServer(
		makeAddLinkEndpoint(svc),
		decodeAddLinkRequest,
		encodeResponse,
	)

	updateLinkHandler := httptransport.NewServer(
		makeUpdateLinkEndpoint(svc),
		decodeUpdateLinkRequest,
		encodeResponse,
	)

	deleteLinkHandler := httptransport.NewServer(
		makeDeleteLinkEndpoint(svc),
		decodeDeleteLinkRequest,
		encodeResponse,
	)

	r := mux.NewRouter()
	r.Methods("GET").Path("/links").Handler(getLinksHandler)
	r.Methods("POST").Path("/links").Handler(addLinkHandler)
	r.Methods("PUT").Path("/links").Handler(updateLinkHandler)
	r.Methods("DELETE").Path("/links").Handler(deleteLinkHandler)

	log.Println("Listening on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", r))
}
