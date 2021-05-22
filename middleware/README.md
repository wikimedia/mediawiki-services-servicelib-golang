# middleware

## PrometheusInstrumentationMiddleware

```golang
package main

import (
    "net/http"

    "github.com/eevans/servicelib-golang/middleware"
    "github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
    reqCounter = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "http_requests_total",
            Help: "Count of HTTP requests processed, partitioned by status code and HTTP method.",
        },
        []string{"code", "method"},
    )

    durationHisto = prometheus.NewHistogramVec(
        prometheus.HistogramOpts{
            Name:    "http_request_duration_seconds",
            Help:    "A histogram of latencies for requests, partitioned by status code and HTTP method.",
            Buckets: []float64{.001, .0025, .0050, .01, .025, .050, .10, .25, .50, 1},
        },
        []string{"code", "method"},
    )
)

func init() {
    prometheus.MustRegister(reqCounter, durationHisto)
}

func main() {
    // Hello, world
    hello := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        io.WriteString(w, "<html><body>Hello World!</body></html>")
    })

    // Wrap our handler function in middleware
    handler := PrometheusInstrumentationMiddleware(reqCounter, durationHisto)(hello)

    http.Handle("/hw", handler)
    http.Handle("/metrics", promhttp.Handler())

    http.ListenAndServe(":8090", nil)
}
```
