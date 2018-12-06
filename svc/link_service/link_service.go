package main

import (
	"log"
	"net/http"

	httptransport "github.com/go-kit/kit/transport/http"
	lm "github.com/the-gigi/delinkcious/pkg/link_manager"
	sgm "github.com/the-gigi/delinkcious/pkg/social_graph_client"
)

func main() {
	store, err := lm.NewDbLinkStore("localhost", 5432, "postgres", "postgres")
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

	http.Handle("/links", getLinksHandler)
	http.Handle("/addLink", addLinkHandler)
	http.Handle("/updateLink", updateLinkHandler)
	http.Handle("/deleteLink", deleteLinkHandler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
