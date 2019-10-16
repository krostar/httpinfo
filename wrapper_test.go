package httpinfo

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

// nolint: funlen
func TestResponseRecorder_wrapped(t *testing.T) {
	t.Run("w", func(t *testing.T) {
		t.Parallel()
		type w interface {
			http.ResponseWriter
		}

		var (
			i   w
			orw struct{ w } // nolint: unused
			rr  = responseRecorder{writer: &orw}
		)
		require.Implements(t, &i, rr.wrapped())
	})
	t.Run("wf", func(t *testing.T) {
		t.Parallel()
		type wf interface {
			http.ResponseWriter
			http.Flusher
		}

		var (
			i   wf
			orw struct{ wf } // nolint: unused
			rr  = responseRecorder{writer: &orw}
		)
		require.Implements(t, &i, rr.wrapped())
	})
	t.Run("wp", func(t *testing.T) {
		t.Parallel()
		type wp interface {
			http.ResponseWriter
			http.Pusher
		}

		var (
			i   wp
			orw struct{ wp } // nolint: unused
			rr  = responseRecorder{writer: &orw}
		)
		require.Implements(t, &i, rr.wrapped())
	})
	t.Run("wh", func(t *testing.T) {
		t.Parallel()
		type wh interface {
			http.ResponseWriter
			http.Hijacker
		}

		var (
			i   wh
			orw struct{ wh } // nolint: unused
			rr  = responseRecorder{writer: &orw}
		)
		require.Implements(t, &i, rr.wrapped())
	})
	t.Run("wpf", func(t *testing.T) {
		t.Parallel()
		type wpf interface {
			http.ResponseWriter
			http.Pusher
			http.Flusher
		}

		var (
			i   wpf
			orw struct{ wpf } // nolint: unused
			rr  = responseRecorder{writer: &orw}
		)
		require.Implements(t, &i, rr.wrapped())
	})
	t.Run("whf", func(t *testing.T) {
		t.Parallel()
		type whf interface {
			http.ResponseWriter
			http.Hijacker
			http.Flusher
		}

		var (
			i   whf
			orw struct{ whf } // nolint: unused
			rr  = responseRecorder{writer: &orw}
		)
		require.Implements(t, &i, rr.wrapped())
	})
	t.Run("whp", func(t *testing.T) {
		t.Parallel()
		type whp interface {
			http.ResponseWriter
			http.Hijacker
			http.Pusher
		}

		var (
			i   whp
			orw struct{ whp } // nolint: unused
			rr  = responseRecorder{writer: &orw}
		)
		require.Implements(t, &i, rr.wrapped())
	})
	t.Run("whpf", func(t *testing.T) {
		t.Parallel()
		type whpf interface {
			http.ResponseWriter
			http.Hijacker
			http.Pusher
			http.Flusher
		}

		var (
			i   whpf
			orw struct{ whpf } // nolint: unused
			rr  = responseRecorder{writer: &orw}
		)
		require.Implements(t, &i, rr.wrapped())
	})
}

func TestWrapFlusher_Flush(t *testing.T) {
	var (
		rw responseWriterMock
		f  flusherMock

		rr = responseRecorder{
			writer: &rw,
		}
		wf = wrapFlusher{
			flusher:  &f,
			recorder: &rr,
		}
	)

	rw.On("WriteHeader", http.StatusOK).Once()
	f.On("Flush").Once()

	wf.Flush()

	rw.AssertExpectations(t)
	f.AssertExpectations(t)
}
