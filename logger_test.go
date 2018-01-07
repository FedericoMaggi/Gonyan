package gonyan

import (
	"os"
	"testing"
)

// TestNewLoggerWithStdout will test that the os.Stdout
// file is compatible with current Stream implementation.
// TODO: Capture stdout content to verify correct logging.
func TestNewLoggerWithStdout(t *testing.T) {
	l := NewLogger("TestNewLoggerWithStdout", false)
	l.RegisterStream(Debug, os.Stdout)

	// Expected log is:
	//  {"tag":"TestNewLoggerWithStdout","message":"this is a debug log and should appear on stdout"}
	l.Debug("this is a debug log and should appear on stdout")
}

// TestNewLoggerWithStderr will test that the os.Stdout
// file is compatible with current Stream implementation.
// TODO: Capture stderr content to verify correct logging.
func TestNewLoggerWithStderr(t *testing.T) {
	l := NewLogger("TestNewLoggerWithStderr", true)
	l.RegisterStream(Error, os.Stderr)

	// Expected log is:
	// 	{"tag":"TestNewLoggerWithStderr","message":"this is an error log and should appear on stderr"}
	l.Error("this is an error log and should appear on stderr")
}

// TestLoggerSetMetadata verifies that metadata are properly set
// to the internal field.
func TestLoggerSetMetadata(t *testing.T) {
	metadata := map[string]string{
		"custom": "field",
	}

	l := NewLogger("TestLoggerSetMetadata", false)

	// Set metadata.
	l.SetMetadata(metadata)
	if l.metadata == nil {
		t.Fatal("Metadata is still nil, should have been set!")
	}

	// Retrieve one metadata field.
	val, ok := l.metadata["custom"]
	if !ok {
		t.Fatal("Metadata with key `custom` not found!")
	}
	if val != "field" {
		t.Fatalf("Unexpected metadata with key `custom` value. Expected: `%s` - Found: `%s`.", "field", val)
	}
}

// TestLoggerClearMetadata verifies that metadata field
// is properly set to nil.
func TestLoggerClearMetadata(t *testing.T) {
	metadata := map[string]string{
		"custom": "field",
	}

	l := NewLogger("TestLoggerSetMetadata", false)
	l.SetMetadata(metadata)
	if l.metadata == nil {
		t.Fatal("Metadata is still nil, should have been set!")
	}

	l.ClearMetadata()
	if l.metadata != nil {
		t.Fatal("Metadata is not nil!")
	}
}

// TestLoggerStreamsProperLogData verifies that the logger sends,
// for all stream types all information with and without metadata.
func TestLoggerStreamsProperLogData(t *testing.T) {
	l := NewLogger("TestLoggerStreamsProperLogData", false)

	stream := newMockStream(1)
	l.RegisterStream(Debug, stream)

	l.Debugf("Hi %s", "there")

	message := <-stream.out
	expected := `{"tag":"TestLoggerStreamsProperLogData","message":"Hi there"}`
	if message != expected {
		t.Fatalf("Unexpected message received from stream. Expected: `%s`, found: `%s`", expected, message)
	}

	l.SetMetadata(map[string]string{"custom": "field"})
	l.Debugf("this log should have metadata")

	message = <-stream.out
	expected = `{"tag":"TestLoggerStreamsProperLogData","message":"this log should have metadata","metadata":{"custom":"field"}}`
	if message != expected {
		t.Fatalf("Unexpected message received from stream. Expected: `%s`, found: `%s`", expected, message)
	}
}
