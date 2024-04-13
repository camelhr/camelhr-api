package middleware

import (
	"net/http"
	"time"

	"github.com/camelhr/log"
	"github.com/go-chi/chi/v5/middleware"
)

// ChiRequestLoggerMiddleware returns a middleware that logs the start and end of each request,
// along with some useful request and response information.
func ChiRequestLoggerMiddleware() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			defer func() {
				status := ww.Status()
				l := log.With(
					"status", status,
					"bytes_written", ww.BytesWritten(),
					"proto", r.Proto,
					"method", r.Method,
					"path", r.URL.Path,
					"query", r.URL.RawQuery,
					"remote_addr", r.RemoteAddr,
					"user_agent", r.UserAgent(),
					"latency", time.Since(start),
				)

				switch {
				case status >= http.StatusInternalServerError:
					l.Error("internal_server_error")
				case status >= http.StatusBadRequest && status < http.StatusInternalServerError:
					l.Warn("invalid_request")
				default:
					l.Info("request_completed")
				}
			}()

			next.ServeHTTP(ww, r)
		}

		return http.HandlerFunc(fn)
	}
}
