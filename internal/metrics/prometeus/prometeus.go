package metrics

import (
	"context"
	"fmt"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var RequestCounter = prometheus.NewCounter(
	prometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "Total number of HTTP requests",
	},
)

func RunMetrics(host string, port int64, errCh chan error) chan struct{} {
	done := make(chan struct{})

	router := http.NewServeMux()
	router.Handle("/metrics", promhttp.Handler())

	server := http.Server{
		Addr:    fmt.Sprintf("%s:%d", host, port),
		Handler: router,
	}

	go func(d chan struct{}) {
		select {
		case <-d:
			server.Shutdown(context.Background())
		}
	}(done)

	go func(ch chan error) {
		if err := server.ListenAndServe(); err != nil {
			ch <- err
		}
	}(errCh)

	return done
}
