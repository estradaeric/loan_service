package middleware

import (
	"net/http"
	"time"

	"loan-service/logger"
)

// LoggingMiddleware logs each incoming HTTP request and its duration using logrus
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		rec := &statusRecorder{
			ResponseWriter: w,
			status:         http.StatusOK, // default status
		}

		next.ServeHTTP(rec, r)

		duration := time.Since(start)

		logger.Log.WithFields(map[string]interface{}{
			"method":       r.Method,
			"path":         r.URL.Path,
			"status":       rec.status,
			"elapsed_ms":   duration.Milliseconds(),
			"bytes_written": rec.bytesWritten,
		}).Info("HTTP Request")
	})
}

// RecoveryMiddleware recovers from any panic and logs the error with logrus
func RecoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				logger.Log.WithField("error", err).Error("Panic recovered")
				http.Error(w, `{"error": "internal server error"}`, http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

// statusRecorder helps capture response status code and bytes written
type statusRecorder struct {
	http.ResponseWriter
	status       int
	bytesWritten int
}

func (rec *statusRecorder) WriteHeader(code int) {
	rec.status = code
	rec.ResponseWriter.WriteHeader(code)
}

func (rec *statusRecorder) Write(b []byte) (int, error) {
	size, err := rec.ResponseWriter.Write(b)
	rec.bytesWritten += size
	return size, err
}