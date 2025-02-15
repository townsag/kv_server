package middleware

import (
	"net/http"
	"log/slog"
	"context"
	"os"
)

type contextKey string
const LoggerKey = contextKey("logger")

func LoggingMiddleware(logger *slog.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// create a request specific logger for this request
		requestLogger := logger.With(
			"method", r.Method,
			"path", r.URL.Path,
			requestIdHeader, r.Header.Get(requestIdHeader),
		)
		requestLogger.Info("received request")

		// create a new context with the request logger 
		ctx := context.WithValue(r.Context(), LoggerKey, requestLogger)
		// create a new request pointer value with the updated context
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

func GetLogger(ctx context.Context) *slog.Logger {
	if logger, ok := ctx.Value(LoggerKey).(*slog.Logger); ok {
		return logger
	}
	return slog.New(slog.NewJSONHandler(os.Stdout, nil))
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