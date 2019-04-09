package news_manager_client

import (
	"context"
	"errors"
	"github.com/go-kit/kit/endpoint"
	"github.com/the-gigi/delinkcious/pb/news_service/pb"
	om "github.com/the-gigi/delinkcious/pkg/object_model"
	"time"
)

type EndpointSet struct {
	GetNewsEndpoint endpoint.Endpoint
}

func encodeGetNewsRequest(_ context.Context, r interface{}) (interface{}, error) {
	req := r.(om.GetNewsRequest)
	return &pb.GetNewsRequest{
		Username:   req.Username,
		StartToken: req.StartToken,
	}, nil
}

func newEvent(e *pb.Event) (event *om.LinkManagerEvent) {
	return &om.LinkManagerEvent{
		EventType: (om.LinkManagerEventTypeEnum)(e.EventType),
		Username:  e.Username,
		Url:       e.Url,
		Timestamp: time.Unix(e.Timestamp.GetSeconds(), (int64)(e.Timestamp.GetNanos())),
	}
}

func decodeGetNewsResponse(_ context.Context, r interface{}) (interface{}, error) {
	gnr := r.(*pb.GetNewsResponse)
	if gnr.Err != "" {
		return nil, errors.New(gnr.Err)
	}

	res := &om.GetNewsResult{
		NextToken: gnr.NextToken,
	}

	for _, e := range gnr.Events {
		res.Events = append(res.Events, newEvent(e))
	}
	return res, nil
}

func (s EndpointSet) GetNews(req om.GetNewsRequest) (result om.GetNewsResult, err error) {
	resp, err := s.GetNewsEndpoint(context.Background(), req)
	if err != nil {
		return
	}
	result = *(resp.(*om.GetNewsResult))
	return
}
