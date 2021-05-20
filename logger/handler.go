package logger

import (
	"context"
	"net"
	"net/http"
)

type contextKey int

const ScopedLogger = iota

// LoggerInjectingMiddleware injects a RequestScopedLogger into a Handler, and pre-populates
// common values from the request (trace ID, client address and port, and the network
// forwarded IP if applicable).
func LoggerInjectingMiddleware(log *Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var ctx context.Context
			var err error
			var forward string
			var id string
			var address string
			var port string
			var reqLog = log.Request()

			if id = r.Header.Get("X-Request-ID"); id != "" {
				reqLog.Trace(id)
			}

			if address, port, err = net.SplitHostPort(r.RemoteAddr); err == nil {
				reqLog.ClientIP(address).ClientPort(port)
			} else {
				log.Error("Unable to parse %q as IP:port", r.RemoteAddr)
			}

			if forward = r.Header.Get("X-Forwarded-For"); forward != "" {
				reqLog.NetworkForwardedIP(forward)
			}

			ctx = context.WithValue(r.Context(), ScopedLogger, reqLog)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
