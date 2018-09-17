package eventgrid

import "encoding/json"

// Envelope - event grid envelope structure
// Custom payload is contained in the Data field
type Envelope struct {
	ID          string `json:"id"`
	Topic       string `json:"topic"`
	Subject     string `json:"subject"`
	EventType   string `json:"eventType"`
	EventTime   string `json:"eventTime"`
	DataVersion string `json:"dataVersion"`

	// Data varies by message so is not processed
	Data json.RawMessage `json:"data"`
}
