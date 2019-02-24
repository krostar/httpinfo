package httpinfo

import "net/http"

// nolint: gocyclo
func (r *responseRecorder) wrapped() http.ResponseWriter {
	var (
		hj, implHijacker = r.writer.(http.Hijacker)
		pu, implPusher   = r.writer.(http.Pusher)
		fl, implFlusher  = r.writer.(http.Flusher)
	)

	if implFlusher {
		fl = &wrapFlusher{flusher: fl, recorder: r}
	}

	switch {
	case !implHijacker && !implPusher && !implFlusher:
		return struct {
			http.ResponseWriter
		}{r}
	case !implHijacker && !implPusher && implFlusher:
		return struct {
			http.ResponseWriter
			http.Flusher
		}{r, fl}
	case !implHijacker && implPusher && !implFlusher:
		return struct {
			http.ResponseWriter
			http.Pusher
		}{r, pu}
	case implHijacker && !implPusher && !implFlusher:
		return struct {
			http.ResponseWriter
			http.Hijacker
		}{r, hj}
	case !implHijacker && implPusher && implFlusher:
		return struct {
			http.ResponseWriter
			http.Pusher
			http.Flusher
		}{r, pu, fl}
	case implHijacker && !implPusher && implFlusher:
		return struct {
			http.ResponseWriter
			http.Hijacker
			http.Flusher
		}{r, hj, fl}
	case implHijacker && implPusher && !implFlusher:
		return struct {
			http.ResponseWriter
			http.Hijacker
			http.Pusher
		}{r, hj, pu}
	case implHijacker && implPusher && implFlusher:
		return struct {
			http.ResponseWriter
			http.Hijacker
			http.Pusher
			http.Flusher
		}{r, hj, pu, fl}
	}

	return nil
}

type wrapFlusher struct {
	flusher  http.Flusher
	recorder *responseRecorder
}

// Flush implements http.Flusher.
func (cf *wrapFlusher) Flush() {
	cf.recorder.WriteHeaderNow()
	cf.flusher.Flush()
}
