package metrics

import "github.com/prometheus/client_golang/prometheus"

var RequestCounter = prometheus.NewCounter(
	prometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "Total number of HTTP requests",
	},
)
