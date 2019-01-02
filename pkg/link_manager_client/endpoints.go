package link_manager_client

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/go-kit/kit/endpoint"
	"net/http"

	om "github.com/the-gigi/delinkcious/pkg/object_model"
)

type SimpleResponse struct {
	Err string
}

func decodeSimpleResponse(_ context.Context, r *http.Response) (interface{}, error) {
	if r.StatusCode != http.StatusOK {
		return nil, errors.New(r.Status)
	}
	var resp SimpleResponse
	err := json.NewDecoder(r.Body).Decode(&resp)
	return resp, err
}

type EndpointSet struct {
	GetLinksEndpoint   endpoint.Endpoint
	AddLinkEndpoint    endpoint.Endpoint
	UpdateLinkEndpoint endpoint.Endpoint
	DeleteLinkEndpoint endpoint.Endpoint
}

func decodeGetLinksResponse(_ context.Context, r *http.Response) (interface{}, error) {
	if r.StatusCode != http.StatusOK {
		return nil, errors.New(r.Status)
	}
	var res om.GetLinksResult
	err := json.NewDecoder(r.Body).Decode(&res)
	return res, err
}

func (s EndpointSet) GetLinks(req om.GetLinksRequest) (result om.GetLinksResult, err error) {
	res, err := s.GetLinksEndpoint(context.Background(), req)
	if err != nil {
		return
	}
	result = res.(om.GetLinksResult)

	return
}

func (s EndpointSet) AddLink(req om.AddLinkRequest) (err error) {
	resp, err := s.AddLinkEndpoint(context.Background(), req)
	if err != nil {
		return err
	}
	response := resp.(SimpleResponse)

	if response.Err != "" {
		err = errors.New(response.Err)
	}
	return
}

func (s EndpointSet) UpdateLink(req om.UpdateLinkRequest) (err error) {
	resp, err := s.UpdateLinkEndpoint(context.Background(), req)
	if err != nil {
		return err
	}
	response := resp.(SimpleResponse)

	if response.Err != "" {
		err = errors.New(response.Err)
	}
	return
}

func (s EndpointSet) DeleteLink(username string, url string) (err error) {
	resp, err := s.DeleteLinkEndpoint(context.Background(), &deleteRequest{Username: username, Url: url})
	if err != nil {
		return err
	}
	response := resp.(SimpleResponse)

	if response.Err != "" {
		err = errors.New(response.Err)
	}
	return
}
