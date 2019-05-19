package service

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/the-gigi/delinkcious/pkg/metrics"
	om "github.com/the-gigi/delinkcious/pkg/object_model"
	"strings"
	"time"
)


// implement function to return ServiceMiddleware
func newMetricsMiddleware() linkManagerMiddleware {
	return func(next om.LinkManager) om.LinkManager {
		m := metricsMiddleware{next,
			map[string]prometheus.Counter{},
			map[string]prometheus.Summary{}}
		methodNames := []string{"GetLinks", "AddLink", "UpdateLink", "DeleteLink"}
		for _, name := range methodNames {
			m.requestCounter[name] = metrics.NewCounter("link", strings.ToLower(name)+"_count", "count # of requests")
			m.requestLatency[name] = metrics.NewSummary("link", strings.ToLower(name)+"_summary", "request summary in milliseconds")

		}
		return m
	}
}

type metricsMiddleware struct {
	next           om.LinkManager
	requestCounter map[string]prometheus.Counter
	requestLatency map[string]prometheus.Summary
}

func (m metricsMiddleware) recordMetrics(name string, begin time.Time) {
	m.requestCounter[name].Inc()
	durationMilliseconds := float64(time.Since(begin).Nanoseconds() * 1000000)
	m.requestLatency[name].Observe(durationMilliseconds)
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
