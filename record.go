package httpinfo

import (
	"context"
	"net/http"
	"time"
)

type ctxKey string

const ctxKeyRR = ctxKey("rr")

// Record records the http response information and helps to reach
// them from any other middleware. See examples on how to use it.
func Record(opts ...Option) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			var (
				ctx = r.Context()
				rr  = &responseRecorder{
					writer: rw,
					routeGetter: func(r *http.Request) string {
						return r.Method + " " + r.URL.Path
					},
					start: time.Now(),
				}
			)
			defer rr.WriteHeaderNow()

			for _, opt := range opts {
				opt(rr)
			}

			ctx = context.WithValue(ctx, ctxKeyRR, rr)
			rw = rr.wrapped()
			r = r.WithContext(ctx)

			next.ServeHTTP(rw, r)
		})
	}
}
