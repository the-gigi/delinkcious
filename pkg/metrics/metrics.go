package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

func NewCounter(service string, name string, help string) prometheus.Counter {
	opts := prometheus.CounterOpts{
		Namespace: "",
		Subsystem: service,
		Name:      name,
		Help:      help,
	}
	counter := promauto.NewCounter(opts)
	return counter
}

func NewSummary(service string, name string, help string) prometheus.Summary {
	opts := prometheus.SummaryOpts{
		Namespace: "",
		Subsystem: service,
		Name:      name,
		Help:      help,
	}

	summary := promauto.NewSummary(opts)
	return summary
}
