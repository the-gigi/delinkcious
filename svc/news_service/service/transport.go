package service

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/the-gigi/delinkcious/pb/news_service/pb"
	om "github.com/the-gigi/delinkcious/pkg/object_model"
)

func newEvent(e *om.LinkManagerEvent) (event *pb.Event) {
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
