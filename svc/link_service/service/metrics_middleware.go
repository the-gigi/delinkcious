package service

import (
	"github.com/go-kit/kit/metrics"
	om "github.com/the-gigi/delinkcious/pkg/object_model"
	"time"
)

// implement function to return ServiceMiddleware
func newMetricsMiddleware(requestCounter metrics.Counter, requestLatency metrics.Histogram) linkManagerMiddleware {
	return func(next om.LinkManager) om.LinkManager {
		return metricsMiddleware{next, requestCounter, requestLatency}
	}
}

type metricsMiddleware struct {
	next           om.LinkManager
	requestCounter metrics.Counter
	requestLatency metrics.Histogram
}

func (m metricsMiddleware) recordMetrics(methodName string, begin time.Time) {
	m.requestCounter.With("method", methodName).Add(1)
	durationMilliseconds := float64(time.Since(begin).Nanoseconds() * 1000)
	m.requestLatency.With("method", methodName).Observe(durationMilliseconds)
}

func (m metricsMiddleware) GetLinks(request om.GetLinksRequest) (result om.GetLinksResult, err error) {
	defer func(begin time.Time) {
		m.recordMetrics("GetLinks", begin)
	}(time.Now())
	result, err = m.next.GetLinks(request)
	return
}

func (m metricsMiddleware) AddLink(request om.AddLinkRequest) error {
	defer func(begin time.Time) {
		m.recordMetrics("AddLink", begin)
	}(time.Now())
	return m.next.AddLink(request)
}

func (m metricsMiddleware) UpdateLink(request om.UpdateLinkRequest) error {
	defer func(begin time.Time) {
		m.recordMetrics("UpdateLink", begin)
	}(time.Now())
	return m.next.UpdateLink(request)
}

func (m metricsMiddleware) DeleteLink(username string, url string) error {
	defer func(begin time.Time) {
		m.recordMetrics("DeleteLink", begin)
	}(time.Now())
	return m.next.DeleteLink(username, url)
}
