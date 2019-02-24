package httpinfo

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestFromRequest(t *testing.T) {
	var request = httptest.NewRequest(http.MethodGet, "/", nil)

	assert.Nil(t, fromRequest(request))

	Record()(http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
		assert.NotNil(t, fromRequest(r))
	})).ServeHTTP(httptest.NewRecorder(), request)
}

func TestIsUsed(t *testing.T) {
	var request = httptest.NewRequest(http.MethodGet, "/", nil)

	assert.False(t, IsUsed(request))

	Record()(http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
		assert.True(t, IsUsed(r))
	})).ServeHTTP(httptest.NewRecorder(), request)
}

func TestStatus(t *testing.T) {
	Record()(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		var status = http.StatusAlreadyReported
		http.HandlerFunc(func(rw http.ResponseWriter, _ *http.Request) {
			rw.WriteHeader(status)
		}).ServeHTTP(rw, r)
		assert.Equal(t, status, Status(r))
	})).ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/", nil))
}

func TestExecutionTime(t *testing.T) {
	if testing.Short() {
		t.Skip()
		return
	}

	Record()(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {
			time.Sleep(500 * time.Millisecond)
		}).ServeHTTP(rw, r)
		assert.True(t, ExecutionTime(r) > 500*time.Millisecond)
	})).ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/", nil))
}

func TestContentLength(t *testing.T) {
	Record()(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		var hello = []byte("hello")
		http.HandlerFunc(func(rw http.ResponseWriter, _ *http.Request) {
			rw.Write(hello) // nolint: errcheck, gosec
		}).ServeHTTP(rw, r)
		assert.Equal(t, len(hello), ContentLength(r))
	})).ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/", nil))
}

func TestRoute(t *testing.T) {
	Record()(http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method+" "+r.URL.Path, Route(r))
	})).ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/", nil))
}
