package link_manager_client

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

import (
	"context"
	httptransport "github.com/go-kit/kit/transport/http"
	om "github.com/the-gigi/delinkcious/pkg/object_model"
)

func NewClient(baseURL string) (om.LinkManager, error) {
	// Quickly sanitize the instance string.
	if !strings.HasPrefix(baseURL, "http") {
		baseURL = "http://" + baseURL
	}
	u, err := url.Parse(baseURL)
	if err != nil {
		return nil, err
	}

	getLinksEndpoint := httptransport.NewClient(
		"GET",
		copyURL(u, "/links"),
		encodeGetLinksRequest,
		decodeGetLinksResponse).Endpoint()

	addLinkEndpoint := httptransport.NewClient(
		"POST",
		copyURL(u, "/links"),
		encodeHTTPGenericRequest,
		decodeSimpleResponse).Endpoint()

	updateLinkEndpoint := httptransport.NewClient(
		"PUT",
		copyURL(u, "/links"),
		encodeHTTPGenericRequest,
		decodeSimpleResponse).Endpoint()

	deleteLinkEndpoint := httptransport.NewClient(
		"DELETE",
		copyURL(u, "/links"),
		encodeHTTPGenericRequest,
		decodeSimpleResponse).Endpoint()

	// Returning the EndpointSet as an interface relies on the
	// EndpointSet implementing the Service methods. That's just a simple bit
	// of glue code.
	return EndpointSet{
		GetLinksEndpoint:   getLinksEndpoint,
		AddLinkEndpoint:    addLinkEndpoint,
		UpdateLinkEndpoint: updateLinkEndpoint,
		DeleteLinkEndpoint: deleteLinkEndpoint,
	}, nil
}

func copyURL(base *url.URL, path string) *url.URL {
	next := *base
	next.Path = path
	return &next
}

// encodeHTTPGenericRequest is a transport/http.EncodeRequestFunc that
// JSON-encodes any request to the request body. Primarily useful in a client.
func encodeHTTPGenericRequest(_ context.Context, r *http.Request, request interface{}) error {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(request); err != nil {
		return err
	}
	r.Body = ioutil.NopCloser(&buf)
	return nil
}

// Extract the request details from the incoming request and add them as query arguments
func encodeGetLinksRequest(ctx context.Context, req *http.Request, request interface{}) error {
	r := request.(om.GetLinksRequest)
	urlRegex := url.QueryEscape(r.UrlRegex)
	titleRegex := url.QueryEscape(r.TitleRegex)
	descriptionRegex := url.QueryEscape(r.DescriptionRegex)
	username := url.QueryEscape(r.Username)
	tag := url.QueryEscape(r.Tag)
	startToken := url.QueryEscape(r.StartToken)

	q := req.URL.Query()
	q.Add("url", urlRegex)
	q.Add("title", titleRegex)
	q.Add("description", descriptionRegex)
	q.Add("username", username)
	q.Add("tag", tag)
	q.Add("start", startToken)
	req.URL.RawQuery = q.Encode()
	return encodeHTTPGenericRequest(ctx, req, request)
}

type deleteRequest struct {
	Username string
	Url      string
}

// Extract the request details from the incoming request and add them as query arguments
func encodeDeleteRequest(ctx context.Context, req *http.Request, request interface{}) error {
	r := request.(deleteRequest)

	username := url.PathEscape(r.Username)
	toDeleteUrl := url.PathEscape(r.Url)

	req.URL.Path += "/" + username + "/" + toDeleteUrl
	return encodeHTTPGenericRequest(ctx, req, request)
}
