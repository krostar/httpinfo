package httpinfo_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/krostar/httpinfo"
)

func myMiddleware(next http.Handler) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(rw, r)

		if httpinfo.IsUsed(r) {
			fmt.Printf("status      = %d\n", httpinfo.Status(r))
			fmt.Printf("bytes wrote = %d\n", httpinfo.ContentLength(r))
		}
	}
}

func myHandler(rw http.ResponseWriter, _ *http.Request) {
	rw.Write([]byte("Hello world")) // nolint: errcheck, gosec
	rw.WriteHeader(http.StatusAlreadyReported)
}

func Example() {
	var srv = httptest.NewServer(
		httpinfo.Record()(
			myMiddleware(http.HandlerFunc(myHandler)),
		),
	)
	defer srv.Close()

	_, err := http.DefaultClient.Get(srv.URL)
	if err != nil {
		panic(err)
	}

	// Output:
	// status      = 208
	// bytes wrote = 11
}
