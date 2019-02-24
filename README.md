# httpinfo

[![godoc](https://img.shields.io/badge/godoc-reference-blue.svg?style=for-the-badge)](https://godoc.org/github.com/krostar/httpinfo)
[![Licence](https://img.shields.io/github/license/krostar/httpinfo.svg?style=for-the-badge)](https://tldrlegal.com/license/mit-license)
![Latest version](https://img.shields.io/github/tag/krostar/httpinfo.svg?style=for-the-badge)

[![Build Status](https://img.shields.io/travis/krostar/httpinfo/master.svg?style=for-the-badge)](https://travis-ci.org/krostar/httpinfo)
[![Code quality](https://img.shields.io/codacy/grade/6b68f760c0dc4ba6abf34078f30c5a87/master.svg?style=for-the-badge)](https://app.codacy.com/project/krostar/httpinfo/dashboard)
[![Code coverage](https://img.shields.io/codacy/coverage/6b68f760c0dc4ba6abf34078f30c5a87.svg?style=for-the-badge)](https://app.codacy.com/project/krostar/httpinfo/dashboard)

You don't need to write `http.ResponseWriter` wrappers anymore.

## Motivation

It has become a common thing to write a wrapper of `http.ResponseWriter` because at some point it was a need to get a request response information like the response status. In a complete request flow, a lot of middlewares require something (some requires the status, the number of bytes wrote, the route pattern used, ...). Moreover, some middleware are also interacting with the response (like a panic or a timeout handler that sets the response status) causing unwanted behaviour (like a net/http log saying the response status should only we wrote one time). The naive approach of wrapping the `http.ResponseWriter` introduce some flaws and/or does not help to fix the already existing ones.

For example here:

```go
type responseWriterWrapper struct{
    writer http.ResponseWriter
    status int
}

func (rww *responseWriterWrapper) WriteHeader(status int) {
    rww.status = status
    writer.WriteHeader(status)
}

// ...
```

If the original `http.ResponseWriter` was implementing `http.Flusher`, it is not the case anymore. We can add the `Flush` method:

```go
func (rww *responseWriterWrapper) Flush() {
    if f, ok := (rww.writer).(http.Flusher); ok {
        f.Flush()
    }
}
```

but now our wrapper **always** implements the `http.Flusher` interface which can cause unwanted behaviour.

For all these reasons I decided to write my last wrapper of `http.ResponseWriter` and it has the following features:

-   records the http **response status**, the **number of bytes wrote**, the **execution time** of the next handler, and helps to retrieve the **route matching pattern**.
-   writes the **response status** at the last possible time, to prevent the multiple status wrote error
-   keeps the **same net/http interfaces implementation of the wrapped `http.ResponseWriter`**
-   **heavily tested**
-   **multi-thread safe**
-   **super simple to use**

## Usage / examples

```go
// during the router setup...
router.Use(
    httpinfo.Record(),
    // other middlewares goes after, even the panic recover one
    myMiddleware,
)

func myMiddleware (rw http.ResponseWriter, r *http.Request ) {
    // call the next handler
    next.ServeHTTP(w, r)

    // within any request handler, you're now able to get response info
    var (
        status        = httpinfo.Status(r)
        route         = httpinfo.Route(r)
        contentLength = httpinfo.ContentLength(r)
        latency       = httpinfo.ExecutionTime(r)
    )
    // ...
}
```

More doc and examples in the httpinfo's [godoc](https://godoc.org/github.com/krostar/httpinfo)

## License

This project is under the MIT licence, please see the LICENCE file.
