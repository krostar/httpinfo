package httpinfo

import (
	"net/http"
	"time"
)

// IsUsed returns true if the Record middleware
// has been used and false otherwise.
func IsUsed(r *http.Request) bool {
	return fromRequest(r) != nil
}

// Status returns the response status. The value is likely to
// change throughout the execution of a request (each time
// ResponseWriter.WriteHeader is called). It is advised
// to only use this function after calling the next handler.
func Status(r *http.Request) int {
	if rr := fromRequest(r); rr != nil {
		defer rr.mu.Unlock()
		rr.mu.Lock()
		return rr.status
	}
	return 0
}

// ExecutionTime returns the duration since the request start. The value
// will only have sens after the next handler returned. It is advised
// to only use this function after calling the next handler.
func ExecutionTime(r *http.Request) time.Duration {
	if rr := fromRequest(r); rr != nil {
		defer rr.mu.Unlock()
		rr.mu.Lock()
		return time.Since(rr.start)
	}
	return 0
}

// BytesWrote returns the response content length. The value
// is likely to change throughout the execution of a request
// (each time ResponseWriter.Write is called). It is advised
// to only use this function after calling the next handler.
func BytesWrote(r *http.Request) int {
	if rr := fromRequest(r); rr != nil {
		defer rr.mu.Unlock()
		rr.mu.Lock()
		return rr.length
	}
	return 0
}

// RouteUsed returns the route that matches the request. The value is
// returned thanks to the WithRouteGetterFunc option. It is advised
// to only use this function after calling the next handler.
func RouteUsed(r *http.Request) string {
	if rr := fromRequest(r); rr != nil {
		defer rr.mu.Unlock()
		rr.mu.Lock()
		return rr.routeGetter(r)
	}
	return ""
}

func fromRequest(r *http.Request) *responseRecorder {
	if rr, ok := r.Context().Value(ctxKeyRR).(*responseRecorder); ok {
		return rr
	}
	return nil
}
