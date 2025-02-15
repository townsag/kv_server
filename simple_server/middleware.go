package main

import (
	"net/http"
	"github.com/google/uuid"
	"log/slog"
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


func LoggingMiddleware(logger *slog.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.Info("received request", "method", r.Method, "path", r.URL.Path, requestIdHeader, r.Header.Get(requestIdHeader))
		next.ServeHTTP(w, r)
	})
}

// wrap the logging middleware function with a function that takes the logger 
// as an argument and passes that logger via a closure to the anonymous function
// which conforms to the function signature of taking a handler as an in put and
// returning a handler
func NewLoggingMiddleware(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return LoggingMiddleware(logger, h)
	}
}


// func 