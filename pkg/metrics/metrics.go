package metrics

import (
	kit_prom "github.com/go-kit/kit/metrics/prometheus"
	std_prom "github.com/prometheus/client_golang/prometheus"
)

func NewCounter(service string, name string, help string) *kit_prom.Counter {
	opts := std_prom.CounterOpts{
		Namespace: "",
		Subsystem: service,
		Name:      name,
		Help:      help,
	}
	counter := kit_prom.NewCounterFrom(opts, []string{"method"})
	return counter
}

func NewSummary(service string, name string, help string) *kit_prom.Summary {
	opts := std_prom.SummaryOpts{
		Namespace: "",
		Subsystem: service,
		Name:      name,
		Help:      help,
	}

	summary := kit_prom.NewSummaryFrom(opts, []string{"method"})
	return summary
}
