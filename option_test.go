package httpinfo

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWithRouteGetterFunc(t *testing.T) {
	var rr responseRecorder

	WithRouteGetterFunc(func(r *http.Request) string { return "" })(&rr)
	assert.NotNil(t, rr.routeGetter)
}
