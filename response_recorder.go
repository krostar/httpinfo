package httpinfo

import (
	"net/http"
	"sync"
	"time"
)

type responseRecorder struct {
	writer http.ResponseWriter
	mu     sync.Mutex

	routeGetter func(r *http.Request) string

	status      int
	statusWrote bool
	statusSet   bool

	length int
	start  time.Time
}

// WriteHeader implements http.ResponseWriter.
// This does not actually write the status but only save it.
// This does not call the underlying response writer WriteHeader method.
func (r *responseRecorder) WriteHeader(status int) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.status = status
	r.statusSet = true
}

// Write implements http.ResponseWriter.
// It writes the status if it has not already been wrote, and call
// the underlying response writer Write method.
func (r *responseRecorder) Write(b []byte) (int, error) {
	r.WriteHeaderNow()

	r.mu.Lock()
	defer r.mu.Unlock()

	n, err := r.writer.Write(b)
	r.length += n
	return n, err
}

// Header implements http.ResponseWriter.
// It simply calls the underlying response writer Header method.
func (r *responseRecorder) Header() http.Header {
	return r.writer.Header()
}

// WriteHeaderNow writes the http status code if it as not already
// been wrote, otherwise it does nothing. If no previous call to
// WriteHeader were made before calling WriteHeaderNow, the same default
// behavior of net/http is applied.
func (r *responseRecorder) WriteHeaderNow() {
	r.mu.Lock()
	defer r.mu.Unlock()

	if !r.statusWrote { // avoid multiple written status http error
		if !r.statusSet {
			r.mu.Unlock()
			r.WriteHeader(http.StatusOK) // same default behavious as net/http
			r.mu.Lock()
		}
		r.statusWrote = true
		r.writer.WriteHeader(r.status)
	}
}
