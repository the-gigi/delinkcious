package main

import (
	"context"
	"fmt"
	"github.com/go-kit/kit/endpoint"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/the-gigi/delinkcious/pb/news_service/pb"
	nm "github.com/the-gigi/delinkcious/pkg/news_manager"
	om "github.com/the-gigi/delinkcious/pkg/object_model"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
)

func main() {
	fmt.Println("Started")
	Run()
}

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

	fmt.Println("33333333")

	var store nm.Store
	if redisHostname == "" {
		fmt.Println("44444444 in-memory")
		store = nm.NewInMemoryNewsStore()
	} else {
		fmt.Println("5555555 redis")
		address := fmt.Sprintf("%s:%s", redisHostname, redisPort)
		store, err = nm.NewRedisNewsStore(address)
		if err != nil {
			log.Fatal(err)
		}
	}

	natsHostname := os.Getenv("NATS_CLUSTER_SERVICE_HOST")
	natsPort := os.Getenv("NATS_CLUSTER_SERVICE_PORT")

	fmt.Println("66666 nats:", natsHostname, ":", natsPort)

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

func newEvent(e *om.Event) (event *pb.Event) {
	event = &pb.Event{
		EventType: (pb.EventType)(e.EventType),
		Username:  e.Username,
		Url:       e.Url,
	}

	seconds := e.Timestamp.Unix()
	nanos := (int32(e.Timestamp.UnixNano() - 1e9*seconds))
	event.Timestamp = &timestamp.Timestamp{Seconds: seconds, Nanos: nanos}
	return
}

func decodeGetNewsRequest(_ context.Context, r interface{}) (interface{}, error) {
	request := r.(*pb.GetNewsRequest)
	return om.GetNewsRequest{
		Username:   request.Username,
		StartToken: request.StartToken,
	}, nil
}

func encodeGetNewsResponse(_ context.Context, r interface{}) (interface{}, error) {
	return r, nil
}

func makeGetNewsEndpoint(svc om.NewsManager) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(om.GetNewsRequest)
		r, err := svc.GetNews(req)
		res := &pb.GetNewsResponse{
			Events:    []*pb.Event{},
			NextToken: r.NextToken,
		}
		if err != nil {
			res.Err = err.Error()
		}
		for _, e := range r.Events {
			event := newEvent(e)
			res.Events = append(res.Events, event)
		}
		return res, nil
	}
}

type handler struct {
	getNews grpctransport.Handler
}

func (s *handler) GetNews(ctx context.Context, r *pb.GetNewsRequest) (*pb.GetNewsResponse, error) {
	_, resp, err := s.getNews.ServeGRPC(ctx, r)
	if err != nil {
		return nil, err
	}

	return resp.(*pb.GetNewsResponse), nil
}

func newNewsServer(svc om.NewsManager) pb.NewsServer {
	return &handler{
		getNews: grpctransport.NewServer(
			makeGetNewsEndpoint(svc),
			decodeGetNewsRequest,
			encodeGetNewsResponse,
		),
	}
}
