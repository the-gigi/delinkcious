package user_client

import (
	"bytes"
	"context"
	"encoding/json"
	httptransport "github.com/go-kit/kit/transport/http"
	om "github.com/the-gigi/delinkcious/pkg/object_model"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

func NewClient(baseURL string) (om.UserManager, error) {
	// Quickly sanitize the instance string.
	if !strings.HasPrefix(baseURL, "http") {
		baseURL = "http://" + baseURL
	}
	u, err := url.Parse(baseURL)
	if err != nil {
		return nil, err
	}

	registerEndpoint := httptransport.NewClient(
		"POST",
		copyURL(u, "/register"),
		encodeHTTPGenericRequest,
		decodeSimpleResponse).Endpoint()

	loginEndpoint := httptransport.NewClient(
		"POST",
		copyURL(u, "/login"),
		encodeHTTPGenericRequest,
		decodeLoginResponse).Endpoint()

	logoutEndpoint := httptransport.NewClient(
		"POST",
		copyURL(u, "/logout"),
		encodeHTTPGenericRequest,
		decodeSimpleResponse).Endpoint()

	// Returning the EndpointSet as an interface relies on the
	// EndpointSet implementing the Service methods. That's just a simple bit
	// of glue code.
	return EndpointSet{
		RegisterEndpoint: registerEndpoint,
		LoginEndpoint:    loginEndpoint,
		LogoutEndpoint:   logoutEndpoint,
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
