package eventgrid

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// Handler - handle the event grid message
func Handler(next func(w http.ResponseWriter, r *http.Request, env []Envelope)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var env Envelope
		var err error
		var msg []Envelope

		// validate the request
		if r.Body == nil {
			log.Println(err)
			w.WriteHeader(500)
			return
		}
		defer r.Body.Close()

		// decode the event grid message from the body
		if err = json.NewDecoder(r.Body).Decode(&msg); err != nil {
			log.Println(err)
			w.WriteHeader(500)
			return
		}

		if len(msg) > 0 {
			env = msg[0]

			// handle event grid subscription validation events
			if env.EventType == "Microsoft.EventGrid.SubscriptionValidationEvent" {
				r.URL.RawQuery = "validate"

				// handle the event grid validation event
				if err = handleValidate(w, &env); err != nil {
					log.Println(err)
					w.WriteHeader(500)
					return
				}
			}

			return
		}

		// call the next handler
		if next != nil {
			next(w, r, msg)
		}
	})
}

// handleValidate - handle the event grid webhook validation request
func handleValidate(w http.ResponseWriter, msg *Envelope) error {
	// get the validationCode from the json (that's all we care about)
	var vData struct {
		ValidationCode string `json:"validationCode"`
		ValidationURL  string `json:"validationUrl"`
	}

	// handle the json error
	if err := json.Unmarshal(msg.Data, &vData); err != nil {
		return err
	}

	// return the validationCode as json
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)

	// echo the validation code back to event grid
	fmt.Fprintf(w, "{ \"validationResponse\": \"%v\" }", vData.ValidationCode)
	log.Println("EventGridValidation: Success")

	return nil
}
