package service

import (
	"github.com/go-kit/kit/log"
	om "github.com/the-gigi/delinkcious/pkg/object_model"
	"time"
)

// implement function to return ServiceMiddleware
func newLoggingMiddleware(logger log.Logger) linkManagerMiddleware {
	return func(next om.LinkManager) om.LinkManager {
		return loggingMiddleware{next, logger}
	}
}

// Make a new type and wrap into Service interface
// Add logger property to this type
type loggingMiddleware struct {
	linkManager om.LinkManager
	logger      log.Logger
}

func (m loggingMiddleware) GetLinks(request om.GetLinksRequest) (result om.GetLinksResult, err error) {
	defer func(begin time.Time) {
		m.logger.Log(
			"method", "GetLinks",
			"request", request,
			"result", result,
			"duration", time.Since(begin),
		)
	}(time.Now())
	result, err = m.linkManager.GetLinks(request)
	return
}

func (m loggingMiddleware) AddLink(request om.AddLinkRequest) error {
	return m.linkManager.AddLink(request)
}

func (m loggingMiddleware) UpdateLink(request om.UpdateLinkRequest) error {
	return m.linkManager.UpdateLink(request)
}

func (m loggingMiddleware) DeleteLink(username string, url string) error {
	return m.linkManager.DeleteLink(username, url)
}
