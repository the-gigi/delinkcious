package service

import (
	"github.com/opentracing/opentracing-go"
	om "github.com/the-gigi/delinkcious/pkg/object_model"
)

// implement function to return ServiceMiddleware
func newTracingMiddleware(tracer opentracing.Tracer) linkManagerMiddleware {
	return func(next om.LinkManager) om.LinkManager {
		return tracingMiddleware{next, tracer}
	}
}

type tracingMiddleware struct {
	next   om.LinkManager
	tracer opentracing.Tracer
}

func (m tracingMiddleware) GetLinks(request om.GetLinksRequest) (result om.GetLinksResult, err error) {
	defer func(span opentracing.Span) {
		span.Finish()
	}(m.tracer.StartSpan("GetLinks"))
	result, err = m.next.GetLinks(request)
	return
}

func (m tracingMiddleware) AddLink(request om.AddLinkRequest) error {
	return m.next.AddLink(request)
}

func (m tracingMiddleware) UpdateLink(request om.UpdateLinkRequest) error {
	return m.next.UpdateLink(request)
}

func (m tracingMiddleware) DeleteLink(username string, url string) error {
	return m.next.DeleteLink(username, url)
}
