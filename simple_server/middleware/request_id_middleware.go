package middleware

import (
	"net/http"
	"github.com/google/uuid"
)

// https://www.alexedwards.net/blog/making-and-using-middleware

var requestIdHeader string = "X-Request-ID"

func RequestIdMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestId := r.Header.Get(requestIdHeader)
		if requestId == "" {
			requestId = uuid.New().String()
			r.Header.Add(requestIdHeader, requestId)
		}
		next.ServeHTTP(w, r)
	})
}





// func 