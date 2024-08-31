package metrics

import (
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
)

type Metrics struct {
	totalHTTPRequests *prometheus.CounterVec
	Registry          *prometheus.Registry
}

func Register() *Metrics {
	totalHTTPRequests := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "total_http_requests",
			Help: "Number of HTTP requests processed, labeled by status code, method, and path.",
		},
		[]string{"code", "method", "path"},
	)

	registry := prometheus.NewRegistry()
	registry.MustRegister(
		collectors.NewGoCollector(),
		collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
		totalHTTPRequests,
	)

	return &Metrics{
		totalHTTPRequests: totalHTTPRequests,
		Registry:          registry,
	}
}

func (m *Metrics) LogRequest(status int, method string, path string) {
	m.totalHTTPRequests.WithLabelValues(strconv.Itoa(status), method, path).Inc()
}
