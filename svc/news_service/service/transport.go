package service

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	"github.com/the-gigi/delinkcious/pb/news_service/pb"
	om "github.com/the-gigi/delinkcious/pkg/object_model"
)

func newEvent(e *om.Event) (event *pb.Event) {
	event = &pb.Event{
		EventType: (pb.EventType)(e.EventType),
		Username:  e.Username,
		Url:       e.Url,
	}

	event.Timestamp.Seconds = e.Timestamp.Unix()
	event.Timestamp.Nanos = (int32(e.Timestamp.UnixNano() - 1e9*event.Timestamp.Seconds))

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
	resp := r.(om.GetNewsResponse)
	result := &pb.GetNewsResponse{
		NextToken: resp.NextToken,
	}
	for _, e := range resp.Events {
		event := newEvent(e)
		result.Events = append(result.Events, event)
	}
	return result, nil
}

func makeGetNewsEndpoint(svc om.NewsManager) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(om.GetNewsRequest)
		return svc.GetNews(req)
	}
}

type handler struct {
	getNews grpctransport.Handler
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

func (s *handler) GetNews(ctx context.Context, r *pb.GetNewsRequest) (response *pb.GetNewsResponse, err error) {
	_, resp, err := s.getNews.ServeGRPC(ctx, r)
	if err != nil {
		return nil, err
	}
	return resp.(*pb.GetNewsResponse), nil
}
