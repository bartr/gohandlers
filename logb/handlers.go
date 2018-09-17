package logb

import (
	"net/http"
	"time"
)

//Handler - http handler that writes to log file(s)
func Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		wr := &ResponseLogger{
			ResponseWriter: w,
			status:         0,
			start:          time.Now().UTC()}

		if next != nil {
			next.ServeHTTP(wr, r)
		}

		reqLog.Println(wr.status, time.Now().UTC().Sub(wr.start).Nanoseconds()/100000, r.Method, r.URL.Path, r.URL.RawQuery)
	})
}
