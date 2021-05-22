/*
 * Copyright 2021 Eric Evans <eevans@wikimedia.org> and Wikimedia Foundation
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package middleware

import (
	"bufio"
	"io"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/stretchr/testify/require"
)

func TestPrometheusInstrumentationMiddleware(t *testing.T) {
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

	prometheus.MustRegister(reqCounter, durationHisto)

	// Hello, world
	hello := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "<html><body>Hello World!</body></html>")
	})

	// Wrap our handler function in middleware
	handler := PrometheusInstrumentationMiddleware(reqCounter, durationHisto)(hello)

	server := httptest.NewServer(handler)
	defer server.Close()

	// Make a request to generate some metrics
	_, err := http.Get(server.URL)
	require.Nil(t, err)

	metrics := httptest.NewServer(promhttp.Handler())
	defer metrics.Close()

	// Make a request against the prometheus handler, to read the metrics
	res, err := http.Get(metrics.URL)
	require.Nil(t, err)
	defer res.Body.Close()

	// The output we're looking for should be something like:
	//
	//   ...
	//   http_requests_total{code="200",method="GET"} 1
	//   http_request_duration_seconds_count{code="200",method="GET"} 1
	//   ...
	//
	var i int = 0
	var scanner *bufio.Scanner
	var statusOk, methodGet *regexp.Regexp

	statusOk = regexp.MustCompile(`code="200"`)
	methodGet = regexp.MustCompile(`method="GET"`)

	scanner = bufio.NewScanner(res.Body)

	// Process the output line-by-line
	for scanner.Scan() {
		line := scanner.Text()
		// Match where code=200 & method=GET
		if statusOk.MatchString(line) && methodGet.MatchString(line) {
			// Match the counter and histogram metrics
			if strings.HasPrefix(line, "http_requests_total") || strings.HasPrefix(line, "http_request_duration_seconds_count") {
				// The total for each count should be exactly: 1
				if strings.HasSuffix(line, "1") {
					i++
				}
			}
		}

		// Two matches is win.
		if i > 1 {
			return
		}
	}

	// If we get here, something has gone wrong.S
	require.Nil(t, scanner.Err())
	t.FailNow()
}
