package service

import (
	"github.com/the-gigi/delinkcious/pkg/db_util"
	"log"
	"net/http"

	httptransport "github.com/go-kit/kit/transport/http"
	sgm "github.com/the-gigi/delinkcious/pkg/user_manager"
)

func Run() {
	dbHost, dbPort, err := db_util.GetDbEndpoint("user")
	if err != nil {
		log.Fatal(err)
	}

	store, err := sgm.NewDbUserStore(dbHost, dbPort, "postgres", "postgres")
	if err != nil {
		log.Fatal(err)
	}

	svc, err := sgm.NewUserManager(store)
	if err != nil {
		log.Fatal(err)
	}

	registerHandler := httptransport.NewServer(
		makeRegisterEndpoint(svc),
		decodeRegisterRequest,
		encodeResponse,
	)

	LoginHandler := httptransport.NewServer(
		makeLoginEndpoint(svc),
		decodeLoginRequest,
		encodeResponse,
	)

	LogoutHandler := httptransport.NewServer(
		makeLogoutEndpoint(svc),
		decodeLogoutRequest,
		encodeResponse,
	)

	http.Handle("/register", registerHandler)
	http.Handle("/login", LoginHandler)
	http.Handle("/logout", LogoutHandler)

	log.Println("Listening on port 7070...")
	log.Fatal(http.ListenAndServe(":7070", nil))
}
