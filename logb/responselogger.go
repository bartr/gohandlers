package logb

import (
	"net/http"
	"time"
)

// ResponseLogger - wrap http.ResponseWriter to include status and size
type ResponseLogger struct {
	http.ResponseWriter
	status int
	size   int
	start  time.Time
}

// WriteHeader - wraps http.WriteHeader
func (r *ResponseLogger) WriteHeader(status int) {
	// store status for logging
	r.status = status

	r.ResponseWriter.WriteHeader(status)
}

// Write - wraps http.Write
func (r *ResponseLogger) Write(buf []byte) (int, error) {
	n, err := r.ResponseWriter.Write(buf)

	// store bytes written for logging
	if err == nil {
		r.size += n
	}

	return n, err
}
