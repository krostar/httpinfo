package httpinfo

import (
	"net/http"

	"github.com/stretchr/testify/mock"
)

type responseWriterMock struct{ mock.Mock }

func (m *responseWriterMock) Header() http.Header {
	args := m.Called()
	return args.Get(0).(http.Header)
}

func (m *responseWriterMock) Write(b []byte) (int, error) {
	args := m.Called(b)
	return args.Int(0), args.Error(1)
}

func (m *responseWriterMock) WriteHeader(statusCode int) {
	m.Called(statusCode)
}

type flusherMock struct{ mock.Mock }

func (m *flusherMock) Flush() {
	m.Called()
}
