package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Prometheus metrics
var (
	httpRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"status_code"},
	)
)

func init() {
	// Register metrics with Prometheus
	prometheus.MustRegister(httpRequestsTotal)
}

// middleware to track HTTP status codes
func trackStatusCode(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Create a response writer wrapper to capture status code
		wrapper := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		// Call the next handler
		next.ServeHTTP(wrapper, r)

		// Convert status code to string and increment counter
		statusCode := strconv.Itoa(wrapper.statusCode)
		httpRequestsTotal.WithLabelValues(statusCode).Inc()
	})
}

// responseWriter wrapper to capture status code
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// define helloWorld handler
func helloWorld(w http.ResponseWriter, req *http.Request) {
	hostname, err := os.Hostname()
	msg := fmt.Sprintf("Hello World from %s", hostname)
	if err != nil {
		io.WriteString(w, msg)
	}
	io.WriteString(w, msg)
}

// define request_status_code handler that responds with any requested status code
func request_status_code(w http.ResponseWriter, req *http.Request) {
	// Get status code from query parameter, default to 200
	statusCodeStr := req.URL.Query().Get("code")
	statusCode := 200 // default status code

	if statusCodeStr != "" {
		if code, err := strconv.Atoi(statusCodeStr); err == nil && code >= 100 && code <= 599 {
			statusCode = code
		}
	}

	w.WriteHeader(statusCode)
	msg := fmt.Sprintf("Response with status code: %d", statusCode)
	io.WriteString(w, msg)
}

func main() {
	// define vars
	httpPort := ":8080"

	// initialise new servemux and register http handlers
	mux := http.NewServeMux()
	mux.HandleFunc("/", trackStatusCode(helloWorld))
	mux.HandleFunc("/status", trackStatusCode(request_status_code))
	mux.Handle("/metrics", promhttp.Handler())

	// start server
	log.Printf("Starting server on %s", httpPort)
	log.Printf("Metrics available at http://localhost%s/metrics", httpPort)
	err := http.ListenAndServe(httpPort, mux)
	log.Fatal(err)
}
