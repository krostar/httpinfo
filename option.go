package httpinfo

import "net/http"

// Option defines the way to configure the response recorder.
type Option func(rr *responseRecorder)

// WithRouteGetterFunc overrides the default route getter function.
func WithRouteGetterFunc(routeGetter func(r *http.Request) string) Option {
	return func(rr *responseRecorder) {
		rr.routeGetter = routeGetter
	}
}
