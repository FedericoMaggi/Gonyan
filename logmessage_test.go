package gonyan

import (
	"bytes"
	"testing"
	"time"
)

// TestSerialise verifies proper serialisation without
// custom metadata fields.
func TestSerialise(t *testing.T) {
	date := time.Date(2017, time.January, 3, 10, 23, 34, 200, time.UTC).UnixNano()
	logMessage := NewLogMessage("Test", date, "messagestring", nil)
	serialised, err := logMessage.Serialise()
	if err != nil {
		t.Fatalf("Unexpected error: %s", err.Error())
	}
	if serialised == nil {
		t.Fatalf("Serialised is nil")
	}

	if bytes.Compare(serialised, []byte(`{"tag":"Test","timestamp":1483439014000000200,"message":"messagestring"}`)) != 0 {
		t.Fatalf("Serialisation error, unexpected serialised log: %s", serialised)
	}
}

// TestSerialiseWithMetadata verifies proper serialisation with
// custom metadata fields.
func TestSerialiseWithMetadata(t *testing.T) {
	metadata := map[string]string{
		"custom": "field",
	}

	date := time.Date(2017, time.January, 3, 10, 23, 34, 200, time.UTC).UnixNano()
	logMessage := NewLogMessage("Test", date, "messagestring", metadata)
	serialised, err := logMessage.Serialise()
	if err != nil {
		t.Fatalf("Unexpected error: %s", err.Error())
	}
	if serialised == nil {
		t.Fatalf("Serialised is nil")
	}

	expected := `{"tag":"Test","timestamp":1483439014000000200,"message":"messagestring","metadata":{"custom":"field"}}`
	if bytes.Compare(serialised, []byte(expected)) != 0 {
		t.Fatalf("Serialisation error, unexpected serialised log: %s", serialised)
	}
}

// TestDeserialise verifies proper LogMessage deserialisation.
func TestDeserialise(t *testing.T) {
	serialised := []byte(`{"tag":"Test","timestamp":1483439014000000200,"message":"messagestring"}`)
	logMessage, err := Deserialise(serialised)
	if err != nil {
		t.Fatalf("Unexpected deserialisation error: %s", err.Error())
	}

	if logMessage.Tag != "Test" {
		t.Fatalf("Unexpected Tag found. Expected: %s - Found: %s", "Test", logMessage.Tag)
	}

	if logMessage.Message != "messagestring" {
		t.Fatalf("Unexpected Message found. Expected: %s - Found: %s", "messagestring", logMessage.Message)
	}
	if logMessage.Timestamp != 1483439014000000200 {
		t.Fatalf("Unexpected Timestamp found. Expected: %d - Found: %d", 1483439014000000200, logMessage.Timestamp)
	}

	date := time.Unix(0, logMessage.Timestamp).UTC()
	if date.Day() != 3 {
		t.Fatalf("Unexpected Timestamp Day found. Expected: %d - Found: %d", 3, date.Day())
	}
	if date.Month() != time.January {
		t.Fatalf("Unexpected Timestamp Month found. Expected: %s - Found: %s", time.January.String(), date.Month().String())
	}
	if date.Year() != 2017 {
		t.Fatalf("Unexpected Timestamp Year found. Expected: %d - Found: %d", 2017, date.Year())
	}
	if date.Hour() != 10 {
		t.Fatalf("Unexpected Timestamp Hour found. Expected: %d - Found: %d", 10, date.Hour())
	}
	if date.Minute() != 23 {
		t.Fatalf("Unexpected Timestamp Minute found. Expected: %d - Found: %d", 23, date.Minute())
	}
	if date.Second() != 34 {
		t.Fatalf("Unexpected Timestamp Second found. Expected: %d - Found: %d", 34, date.Second())
	}
	if date.Nanosecond() != 200 {
		t.Fatalf("Unexpected Timestamp Nanosecond found. Expected: %d - Found: %d", 200, date.Nanosecond())
	}
}

// TestDeserialisationFailure verifies that, if provided an invalid
// LogMessage, the Deserialisation function correctly fails.
func TestDeserialisationFailure(t *testing.T) {
	// Invalid JSON format.
	input := "{this-is-not-a-json}"
	message, err := Deserialise([]byte(input))
	if message != nil {
		t.Fatalf("Deserialisation should have returned nil => %+v", message)
	}
	if err == nil {
		t.Fatalf("Deserialisation should have failed for input: `%s`", input)
	}
}
