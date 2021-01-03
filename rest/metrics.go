package rest

import (
	"github.com/prometheus/client_golang/prometheus"

	"github.com/bygui86/go-k8s-probes/commons"
)

const (
	namespace = "rest-api"
)

var (
	restRequests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: commons.ServiceName,
			Name:      "rest_requests_total",
			Help:      "Number of rest requests managed",
		},
		[]string{"method"},
	)

	restRequestsTiming = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: namespace,
			Subsystem: commons.ServiceName,
			Name:      "rest_requests_execution_time_milliseconds",
			Help:      "Execution time of rest requests in milliseconds",
			Buckets:   []float64{1e-10, 1e-8, 1e-6, 1e-4, 1e-2, 0.025, 0.05, 0.075, 0.1, 0.125, 0.25, 0.5, 1, 1.5, 2, 2.5, 5, 7.5, 10, 25, 50, 100, 250, 500, 750, 1000, 2500, 5000, 10000},
		},
		[]string{"method"},
	)
)

func RegisterCustomMetrics() {
	prometheus.MustRegister(
		restRequests,
		restRequestsTiming,
	)
}

func IncreaseRestRequests(method string) {
	go restRequests.WithLabelValues(method).Inc()
}

func ObserveRestRequestsTime(method string, timing float64) {
	restRequestsTiming.WithLabelValues(method).Observe(timing)
}
