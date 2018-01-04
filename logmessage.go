package gonyan

import (
	"encoding/json"
	"fmt"
)

// LogMessage structure defines the basic standard object containing a message
// to be logged via a stream implementation.
type LogMessage struct {
	Tag       string `json:"tag"`
	Timestamp int64  `json:"timestamp,omitempty"`
	Message   string `json:"message"`
}

// NewLogMessage builds a new LogMessage and returns its reference.
func NewLogMessage(tag string, timestamp int64, message string) *LogMessage {
	return &LogMessage{
		Tag:       tag,
		Timestamp: timestamp,
		Message:   message,
	}
}

// Serialise uses caller LogMessage data to generate a valid JSON string
// serialised log.
func (m *LogMessage) Serialise() ([]byte, error) {
	messageBytes, err := json.Marshal(m)
	if err != nil {
		return nil, fmt.Errorf("%s", err.Error())
	}
	return messageBytes, nil
}

// Deserialise uses provided data to generate a LogMessage structure.
func Deserialise(messageBytes []byte) (*LogMessage, error) {
	logMessage := &LogMessage{}
	if err := json.Unmarshal(messageBytes, &logMessage); err != nil {
		return nil, fmt.Errorf("%s", err.Error())
	}

	return logMessage, nil
}
