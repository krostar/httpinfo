package httpinfo

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestResponseRecorder_WriteHeader(t *testing.T) {
	var (
		rw     responseWriterMock
		rr     = responseRecorder{writer: &rw}
		status = http.StatusCreated
	)

	rr.WriteHeader(status)

	assert.True(t, rr.statusSet)
	assert.False(t, rr.statusWrote)
	assert.Equal(t, status, rr.status)

	rw.AssertNotCalled(t, "WriteHeader", status)
	rw.AssertExpectations(t)
}

func TestResponseRecorder_Write(t *testing.T) {
	var (
		rw responseWriterMock
		rr = responseRecorder{writer: &rw}

		status = http.StatusCreated
		b      = []byte("hello world")
	)

	rr.WriteHeader(status)

	rw.On("WriteHeader", status).Once()
	rw.On("Write", b).Return(11, nil).Once()

	len, err := rr.Write(b)

	require.NoError(t, err)
	assert.True(t, rr.statusSet)
	assert.True(t, rr.statusWrote)
	assert.Equal(t, status, rr.status)
	assert.Equal(t, 11, len)
	assert.Equal(t, 11, rr.length)

	rw.AssertExpectations(t)
}

func TestResponseRecorder_Header(t *testing.T) {
	var (
		rw responseWriterMock
		rr = responseRecorder{writer: &rw}

		headers = http.Header{"hello": []string{"wo", "r", "ld"}}
	)

	rw.On("Header").Return(headers).Once()
	h := rr.Header()

	assert.Equal(t, headers, h)

	rw.AssertExpectations(t)
}

func TestResponseRecorder_WriteHeaderNow(t *testing.T) {
	t.Run("no status set", func(t *testing.T) {
		t.Parallel()

		var (
			rw responseWriterMock
			rr = responseRecorder{writer: &rw}
		)

		rw.On("WriteHeader", http.StatusOK).Once()
		rr.WriteHeaderNow()

		assert.True(t, rr.statusSet)
		assert.True(t, rr.statusWrote)

		assert.Equal(t, http.StatusOK, rr.status)
		rw.AssertExpectations(t)
	})

	t.Run("status set not wrote", func(t *testing.T) {
		t.Parallel()

		var (
			rw     responseWriterMock
			rr     = responseRecorder{writer: &rw}
			status = http.StatusOK
		)

		assert.False(t, rr.statusSet)
		assert.False(t, rr.statusWrote)

		rr.WriteHeader(status)
		rw.On("WriteHeader", status).Once()
		rr.WriteHeaderNow()

		assert.True(t, rr.statusSet)
		assert.True(t, rr.statusWrote)

		assert.Equal(t, http.StatusOK, rr.status)
		rw.AssertExpectations(t)
	})

	t.Run("status set already wrote", func(t *testing.T) {
		t.Parallel()

		var (
			rw responseWriterMock
			rr = responseRecorder{
				writer:      &rw,
				statusSet:   true,
				statusWrote: true,
				status:      http.StatusOK,
			}
		)

		rr.WriteHeaderNow()
		rw.AssertNotCalled(t, "WriteHeader", http.StatusOK)
		rw.AssertExpectations(t)
	})
}
