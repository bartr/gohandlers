package eventgrid

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
)

// TODO - Log / Reject bad EG messages
// TODO - handle upcoming EG change that can send multiple items in a message

// Handler - handle the event grid message
func Handler(next func(w http.ResponseWriter, r *http.Request, env *Envelope)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var env Envelope
		var err error
		var msg []Envelope

		// validate the request
		if r.Body == nil {
			logError(w, errors.New("No request body"))
			return
		}
		defer r.Body.Close()

		// decode the event grid message from the body
		if err = json.NewDecoder(r.Body).Decode(&msg); err != nil {
			logError(w, err)
			return
		}

		// TODO - future versions of event grid may send more than one message in a request
		env = msg[0]

		// validate the event grid envelope
		if err = ValidateEnvelope(&env); err != nil {
			logError(w, err)
			return
		}

		// handle event grid subscription validation events
		if env.EventType == "Microsoft.EventGrid.SubscriptionValidationEvent" {
			r.URL.RawQuery = "validate"

			// handle the event grid validation event
			if err = handleValidate(w, &env); err != nil {
				logError(w, err)
				return
			}
		} else {
			// call the next handler
			if next != nil {
				next(w, r, &env)
			}
		}
	})
}

// log the error and send a 500 status code
func logError(w http.ResponseWriter, err error) {
	// log any error and return 500
	if err != nil {
		log.Println(err)
		w.WriteHeader(500)
	}
}

// ValidateEnvelope - validates a message grid envelope contains required fields
func ValidateEnvelope(env *Envelope) error {
	// verify event grid ID
	if env.ID == "" {
		return errors.New("Event Grid Envelope: missing ID")
	}

	// verify event grid has data
	// TODO - should we do this? are empty data messages possible?
	if env.Data == nil {
		return errors.New("Event Grid Envelope: missing Data")
	}

	// TODO - add more validations?
	return nil
}

// handle the event grid webhook validation request
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
