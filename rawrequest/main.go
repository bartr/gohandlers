package rawrequest

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"strings"
)

// DisplayRawRequests - display the raw requests
func DisplayRawRequests(w http.ResponseWriter, r *http.Request) {
	s := fmt.Sprintln(requests)

	// trim the []
	s = s[1 : len(s)-2]

	w.WriteHeader(200)
	w.Header().Add("Content-Type", "text/plain")
	// for debug, turn caching off
	w.Header().Add("Cache-Control", "no-cache")
	w.Write([]byte(s))
}

var requests []string
var maxLength = 10

// Handler - http handler that saves the raw request
// this handler can be chained with other handlers
func Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		s := strings.ToLower(r.URL.Path)

		if !strings.HasSuffix(s, "favicon.ico") {

			// get the full request
			b, err := httputil.DumpRequest(r, true)

			if err != nil {
				log.Printf("Error: %s\n", b)
				return
			}

			s := string(b) + "====================\n"

			// prepend to array
			arr := make([]string, len(requests)+1)
			arr[0] = s
			copy(arr[1:], requests)
			requests = arr

			// slice array to maxLength
			if len(requests) > maxLength {
				requests = requests[:maxLength]
			}
		}

		// call the next handler
		if next != nil {
			next.ServeHTTP(w, r)
		}
	})
}
