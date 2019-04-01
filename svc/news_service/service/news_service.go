package service

import (
	"fmt"
	"github.com/the-gigi/delinkcious/pb/news_service/pb"
	nm "github.com/the-gigi/delinkcious/pkg/news_manager"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
)

func Run() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "6060"
	}

	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatal(err)
	}

	redisHostname := os.Getenv("NEWS_MANAGER_REDIS_SERVICE_HOST")
	redisPort := os.Getenv("NEWS_MANAGER_REDIS_SERVICE_PORT")

	var store nm.Store
	if redisHostname == "" {
		store = nm.NewInMemoryNewsStore()
	} else {
		address := fmt.Sprintf("%s:%s", redisHostname, redisPort)
		store, err = nm.NewRedisNewsStore(address)
		if err != nil {
			log.Fatal(err)
		}
	}

	natsHostname := os.Getenv("NATS_CLUSTER_SERVICE_HOST")
	natsPort := os.Getenv("NATS_CLUSTER_SERVICE_PORT")

	svc, err := nm.NewNewsManager(store, natsHostname, natsPort)
	if err != nil {
		log.Fatal(err)
	}

	gRPCServer := grpc.NewServer()
	pb.RegisterNewsServer(gRPCServer, newNewsServer(svc))

	fmt.Printf("News service is listening on port %s...\n", port)
	err = gRPCServer.Serve(listener)
	fmt.Println("Serve() failed", err)
}
