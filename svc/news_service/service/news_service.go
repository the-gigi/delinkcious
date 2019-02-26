package service

import (
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

	svc, err := nm.NewNewsManager()
	if err != nil {
		log.Fatal(err)
	}

	gRPCServer := grpc.NewServer()
	pb.RegisterNewsServer(gRPCServer, newNewsServer(svc))
	gRPCServer.Serve(listener)
}
