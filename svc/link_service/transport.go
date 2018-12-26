package main

import (
	"context"
	"encoding/json"
	"github.com/go-kit/kit/endpoint"
	om "github.com/the-gigi/delinkcious/pkg/object_model"
	"net/http"
	"time"
)

type link struct {
	Url         string
	Title       string
	Description string
	Tags        map[string]bool
	CreatedAt   string
	UpdatedAt   string
}

func newLink(source om.Link) link {
	return link{
		Url:         source.Url,
		Title:       source.Title,
		Description: source.Description,
		Tags:        source.Tags,
		CreatedAt:   source.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   source.UpdatedAt.Format(time.RFC3339),
	}
}

type getLinksResponse struct {
	Links []link `json:"links"`
	Err   string `json:"err"`
}

type deleteLinkRequest struct {
	Username string
	Url      string
}

type SimpleResponse struct {
	Err string `json:"err"`
}

func decodeGetLinksRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request om.GetLinksRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		return nil, err
	}
	return request, nil
}

func decodeAddLinkRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request om.AddLinkRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		return nil, err
	}
	return request, nil
}

func decodeUpdateLinkRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request om.UpdateLinkRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		return nil, err
	}
	return request, nil
}

func decodeDeleteLinkRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request deleteLinkRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		return nil, err
	}
	return request, nil
}

func encodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}

func makeGetLinksEndpoint(svc om.LinkManager) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(om.GetLinksRequest)
		result, err := svc.GetLinks(req)
		res := getLinksResponse{}
		for _, link := range result.Links {
			res.Links = append(res.Links, newLink(link))
		}
		if err != nil {
			res.Err = err.Error()
		}
		return res, nil
	}
}

func makeAddLinkEndpoint(svc om.LinkManager) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(om.AddLinkRequest)
		err := svc.AddLink(req)
		res := SimpleResponse{}
		if err != nil {
			res.Err = err.Error()
		}
		return res, nil
	}
}

func makeUpdateLinkEndpoint(svc om.LinkManager) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(om.UpdateLinkRequest)
		err := svc.UpdateLink(req)
		res := SimpleResponse{}
		if err != nil {
			res.Err = err.Error()
		}
		return res, nil
	}
}

func makeDeleteLinkEndpoint(svc om.LinkManager) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(deleteLinkRequest)
		err := svc.DeleteLink(req.Username, req.Url)
		res := SimpleResponse{}
		if err != nil {
			res.Err = err.Error()
		}
		return res, nil
	}
}
