package logb

import (
	"log"
	"net/http"
	"os"
	"time"
)

// Logger - the Logger that is written to - default is stdout
var Logger = log.New(os.Stdout, "", log.Ldate|log.Ltime)

//Handler - http handler that writes to log file(s)
func Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		id := time.Now().UnixNano()

		wr := &ResponseLogger{
			ResponseWriter: w,
			status:         0,
			start:          time.Now().UTC(),
			duration:       0}

		Logger.Println("Request", id, r.Method, r.URL.Path, r.URL.RawQuery)
		defer Logger.Println("Response", id, wr.status, wr.duration, wr.bytes)

		if next != nil {
			next.ServeHTTP(wr, r)
		}

		wr.duration = time.Now().UTC().Sub(wr.start).Nanoseconds() / 100000
	})
}

// ResponseLogger - wrap http.ResponseWriter to include status and size
type ResponseLogger struct {
	http.ResponseWriter
	status   int
	bytes    int
	start    time.Time
	duration int64
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
		r.bytes += n
	}

	return n, err
}
