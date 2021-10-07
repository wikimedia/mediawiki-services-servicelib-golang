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
	"net/http"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

// statusObserver wraps a ResponseWriter in order to track the status code for later use.
type statusObserver struct {
	http.ResponseWriter
	status int
}

// WriteHeader writes an HTTP response status code to the ResponseWriter and status observer.
func (r *statusObserver) WriteHeader(code int) {
	r.status = code
	r.ResponseWriter.WriteHeader(code)
}

// Returns a new statusObserver with a default status
func newStatusObserver(w http.ResponseWriter) *statusObserver {
	return &statusObserver{w, 200}
}

// PrometheusInstrumentationMiddleware is an HTTP middleware that wraps the provided http.Handler
// to count and observe the request and its duration with the provided CounterVec and HistogramVec.
func PrometheusInstrumentationMiddleware(reqCounter *prometheus.CounterVec, latencyHist *prometheus.HistogramVec) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var start time.Time = time.Now()
			var observer *statusObserver = newStatusObserver(w)

			next.ServeHTTP(observer, r)

			latencyHist.WithLabelValues(strconv.Itoa(observer.status), r.Method).Observe(time.Since(start).Seconds())
			reqCounter.WithLabelValues(strconv.Itoa(observer.status), r.Method).Inc()
		})
	}
}
