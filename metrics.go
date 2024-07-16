package check_status

import "github.com/prometheus/client_golang/prometheus"

var (
	requestCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "poller_requests_total",
			Help: "Total number of requests made to each provider",
		},
		[]string{"provider"},
	)
	requestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "poller_request_duration_seconds",
			Help:    "Duration of requests to each provider",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"provider"},
	)
	errorCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "poller_errors_total",
			Help: "Total number of errors encountered when making requests",
		},
		[]string{"provider"},
	)
	transactionsProcessed = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "poller_transactions_processed_total",
			Help: "Total number of transactions processed",
		},
		[]string{"provider"},
	)
)
