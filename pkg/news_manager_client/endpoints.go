package news_manager_client
//
//import (
//	"context"
//	"errors"
//	"github.com/go-kit/kit/endpoint"
//	"github.com/the-gigi/delinkcious/pb/news_service/pb"
//	om "github.com/the-gigi/delinkcious/pkg/object_model"
//)
//
//type EndpointSet struct {
//	GetNewsEndpoint endpoint.Endpoint
//}
//
//func encodeGetNewsRequest(_ context.Context, r interface{}) (interface{}, error) {
//	req := r.(pb.GetNewsRequest)
//	return &pb.GetNewsRequest{
//		Username:   req.Username,
//		StartToken: req.StartToken,
//	}, nil;
//}
//
//
//func newEvent(e *pb.Event) (event *om.Event) {
//	return
//}
//
//func decodeGetNewsResponse(_ context.Context, r interface{}) (interface{}, error) {
// 	gnr := r.(pb.GetNewsResponse)
// 	if gnr.Err != "" {
// 		return nil, errors.New(gnr.Err)
//	}
//
// 	res := &om.GetNewsResult {
// 		NextToken: gnr.NextToken,
//	}
//
// 	for _, e := range  gnr.Events {
// 		res.Events = append(res.Events, newEvent(e))
//	}
// 	return res, nil
//}
//
//func (s EndpointSet) GetNews(req om.GetNewsRequest) (result om.GetNewsResult, err error) {
//	res, err := s.GetNewsEndpoint(context.Background(), req)
//	if err != nil {
//		return
//	}
//	resp := res.(pb.GetNewsResponse)
//	if resp.Err != "" {
//		err = errors.New(resp.Err)
//		return
//	}
//
//	result.NextToken = resp.NextToken
//	for _, e := range resp.Events {
//		event = newEvent()
//	}
//
//
//	return
//}
