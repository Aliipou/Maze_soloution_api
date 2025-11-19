package middleware

import (
	"net/http"
	"strings"
)

type Middleware func(http.Handler) http.Handler

func ChainMiddleware(h http.Handler, middlewares ...Middleware) http.Handler {
	for _, mw := range middlewares {
		h = mw(h)
	}
	return h
}

func CommonMiddleware(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// * If type is Option return, no need for authentication or content type validation *
		if r.Method == http.MethodOptions {
			next.ServeHTTP(w, r)
			return
		}

		// * The request body should be JSON, and the Content-Type header must start with: application/json *
		if !strings.HasPrefix(r.Header.Get("Content-Type"), "application/json") {
			w.WriteHeader(http.StatusUnsupportedMediaType)
			w.Write([]byte(`{"error": "Content-Type header should be set to: application/json."}`))
			return
		}

		// * Set the Content-Type header of the response to application/json for all responses
		// * On http.Error("..."), the Content-Type header will be set to text/plain; charset=utf-8
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		next.ServeHTTP(w, r)
	})
}
