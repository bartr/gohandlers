package eventgrid

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGet(t *testing.T) {

	r, err := http.NewRequest("GET", "https://www.eventgrid.com/person", nil)

	if err != nil {
		t.Error("NewRequest: ", err)
	}

	w := httptest.NewRecorder()

	h := Handler(testHandler)
	h.ServeHTTP(w, r)

	// this should fail
	if w.Code != 500 {
		t.Error("Error Code: ", w.Code)
	}

}

func TestEmptyBody(t *testing.T) {

	r, err := http.NewRequest("POST", "https://www.eventgrid.com/person", nil)

	if err != nil {
		t.Error("NewRequest: ", err)
	}

	w := httptest.NewRecorder()

	h := Handler(testHandler)
	h.ServeHTTP(w, r)

	// this should fail
	if w.Code != 500 {
		t.Error("Error Code: ", w.Code)
	}

}

func TestGoodBody(t *testing.T) {

	s := `[{"id":"1001","topic":"","subject":"person","eventType":"person","eventTime":"2018-08-20T18:04:26Z","dataVersion":"1.0","data":{"firstName": "John", "LastName": "Doe"}}]`

	r, err := http.NewRequest("POST", "https://www.eventgrid.com/person", bytes.NewBufferString(s))

	if err != nil {
		t.Error("NewRequest: ", err)
	}

	w := httptest.NewRecorder()

	h := Handler(testHandler)
	h.ServeHTTP(w, r)

	if w.Code != 200 {
		t.Error("Error Code: ", w.Code)
	}

}

func TestBadJson(t *testing.T) {

	s := `[{"id":"1001","topic":"","subject":"person","eventType":"person","eventTime":"2018-08-20T18:04:26Z","data":{"firstName: "John", "LastName": "Doe"}}]`

	r, err := http.NewRequest("POST", "https://www.eventgrid.com/person", bytes.NewBufferString(s))

	if err != nil {
		t.Error("NewRequest: ", err)
	}

	w := httptest.NewRecorder()

	h := Handler(testHandler)
	h.ServeHTTP(w, r)

	// this should fail
	if w.Code != 500 {
		t.Error("Error Code: ", w.Code)
	}

}

func TestMissingId(t *testing.T) {

	s := `[{"topic":"","subject":"person","eventType":"person","eventTime":"2018-08-20T18:04:26Z","dataVersion":"1.0","data":{"firstName": "John", "LastName": "Doe"}}]`

	r, err := http.NewRequest("POST", "https://www.eventgrid.com/person", bytes.NewBufferString(s))

	if err != nil {
		t.Error("NewRequest: ", err)
	}

	w := httptest.NewRecorder()

	h := Handler(testHandler)
	h.ServeHTTP(w, r)

	// this should fail
	if w.Code != 500 {
		t.Error("Error Code: ", w.Code)
	}

}

func TestMissingData(t *testing.T) {

	s := `[{"id":"1001","topic":"","subject":"person","eventType":"person","eventTime":"2018-08-20T18:04:26Z","dataVersion":"1.0"}]`

	r, err := http.NewRequest("POST", "https://www.eventgrid.com/person", bytes.NewBufferString(s))

	if err != nil {
		t.Error("NewRequest: ", err)
	}

	w := httptest.NewRecorder()

	h := Handler(testHandler)
	h.ServeHTTP(w, r)

	// this should fail
	if w.Code != 500 {
		t.Error("Error Code: ", w.Code)
	}

}

func TestEventGridValidation(t *testing.T) {
	s := `[ { "id": "1234", "subject": "", "data": { "validationCode": "Validation-Code", "validationUrl": "https://foo" }, "eventType": "Microsoft.EventGrid.SubscriptionValidationEvent", "eventTime": "2018-01-25T22:12:19.4556811Z" } ]`

	r, err := http.NewRequest("POST", "https://www.eventgrid.com/person", bytes.NewBufferString(s))

	if err != nil {
		t.Error("NewRequest: ", err)
	}

	w := httptest.NewRecorder()

	h := Handler(testHandler)
	h.ServeHTTP(w, r)

	if w.Code != 200 {
		t.Error("Error Code: ", w.Code)
	}
}

func testHandler(w http.ResponseWriter, r *http.Request, env *Envelope) {
	w.WriteHeader(200)
}
