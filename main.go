package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	// Counter metric to track the number of HTTP requests
	httpRequestsTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "Total number of HTTP requests",
	})

	// Gauge metric for an example random value
	randomValue = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "example_random_value",
		Help: "Example gauge showing a random value between 0 and 100",
	})
)

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Print out a log when the /metrics endpoint is hit
		if r.URL.Path == "/metrics" {
			fmt.Printf("Endpoint '/metrics' hit at %v\n", time.Now())
		}
		next.ServeHTTP(w, r)
	})
}

func main() {
	// Update the random value gauge every 5 seconds
	go func() {
		for {
			randomValue.Set(float64(rand.Intn(100)))
			time.Sleep(5 * time.Second)
		}
	}()

	// Define routes
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Increment the request counter
		httpRequestsTotal.Inc()
		w.Write([]byte("Hello, World! Visit /metrics to see Prometheus metrics."))
	})

	metricsHandler := promhttp.Handler()

	// Expose metrics endpoint
	http.Handle("/metrics", loggingMiddleware(metricsHandler))

	// Start the server
	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
