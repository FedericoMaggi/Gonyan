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

// TestLoggerStreamsProperLogDataForDebug verifies that the logger
// sends all information with and without metadata for Debug stream.
func TestLoggerStreamsProperLogDataForDebug(t *testing.T) {
	l := NewLogger("TestLoggerStreamsProperLogDataForDebug", false)

	stream := newMockStream(1)
	l.RegisterStream(Debug, stream)

	l.Debugf("Hi %s", "there")

	message := <-stream.out
	expected := `{"tag":"TestLoggerStreamsProperLogDataForDebug","message":"Hi there"}`
	if message != expected {
		t.Fatalf("Unexpected message received from stream. Expected: `%s`, found: `%s`", expected, message)
	}

	l.SetMetadata(map[string]string{"custom": "field"})
	l.Debugf("this log should have metadata")

	message = <-stream.out
	expected = `{"tag":"TestLoggerStreamsProperLogDataForDebug","message":"this log should have metadata","metadata":{"custom":"field"}}`
	if message != expected {
		t.Fatalf("Unexpected message received from stream. Expected: `%s`, found: `%s`", expected, message)
	}
}

// TestLoggerStreamsProperLogDataForVerbose verifies that the logger
// sends all information with and without metadata for Verbose stream.
func TestLoggerStreamsProperLogDataForVerbose(t *testing.T) {
	l := NewLogger("TestLoggerStreamsProperLogDataForVerbose", false)

	stream := newMockStream(1)
	l.RegisterStream(Verbose, stream)

	l.Verbose("Hi there")

	message := <-stream.out
	expected := `{"tag":"TestLoggerStreamsProperLogDataForVerbose","message":"Hi there"}`
	if message != expected {
		t.Fatalf("Unexpected message received from stream. Expected: `%s`, found: `%s`", expected, message)
	}

	l.SetMetadata(map[string]string{"custom": "field"})
	l.Verbosef("this log should have metadata")

	message = <-stream.out
	expected = `{"tag":"TestLoggerStreamsProperLogDataForVerbose","message":"this log should have metadata","metadata":{"custom":"field"}}`
	if message != expected {
		t.Fatalf("Unexpected message received from stream. Expected: `%s`, found: `%s`", expected, message)
	}
}

// TestLoggerStreamsProperLogDataForInfo verifies that the logger
// sends all information with and without metadata for Info stream.
func TestLoggerStreamsProperLogDataForInfo(t *testing.T) {
	l := NewLogger("TestLoggerStreamsProperLogDataForInfo", false)

	stream := newMockStream(1)
	l.RegisterStream(Info, stream)

	l.Info("Hi there")

	message := <-stream.out
	expected := `{"tag":"TestLoggerStreamsProperLogDataForInfo","message":"Hi there"}`
	if message != expected {
		t.Fatalf("Unexpected message received from stream. Expected: `%s`, found: `%s`", expected, message)
	}

	l.SetMetadata(map[string]string{"custom": "field"})
	l.Infof("this log should have metadata")

	message = <-stream.out
	expected = `{"tag":"TestLoggerStreamsProperLogDataForInfo","message":"this log should have metadata","metadata":{"custom":"field"}}`
	if message != expected {
		t.Fatalf("Unexpected message received from stream. Expected: `%s`, found: `%s`", expected, message)
	}
}

// TestLoggerStreamsProperLogDataForWarning verifies that the logger
// sends all information with and without metadata for Warning stream.
func TestLoggerStreamsProperLogDataForWarning(t *testing.T) {
	l := NewLogger("TestLoggerStreamsProperLogDataForWarning", false)

	stream := newMockStream(1)
	l.RegisterStream(Warning, stream)

	l.Warning("Hi there")

	message := <-stream.out
	expected := `{"tag":"TestLoggerStreamsProperLogDataForWarning","message":"Hi there"}`
	if message != expected {
		t.Fatalf("Unexpected message received from stream. Expected: `%s`, found: `%s`", expected, message)
	}

	l.SetMetadata(map[string]string{"custom": "field"})
	l.Warningf("this log should have metadata")

	message = <-stream.out
	expected = `{"tag":"TestLoggerStreamsProperLogDataForWarning","message":"this log should have metadata","metadata":{"custom":"field"}}`
	if message != expected {
		t.Fatalf("Unexpected message received from stream. Expected: `%s`, found: `%s`", expected, message)
	}
}

// TestLoggerStreamsProperLogDataForError verifies that the logger
// sends all information with and without metadata for Error stream.
func TestLoggerStreamsProperLogDataForError(t *testing.T) {
	l := NewLogger("TestLoggerStreamsProperLogDataForError", false)

	stream := newMockStream(1)
	l.RegisterStream(Error, stream)

	l.Error("Hi there")

	message := <-stream.out
	expected := `{"tag":"TestLoggerStreamsProperLogDataForError","message":"Hi there"}`
	if message != expected {
		t.Fatalf("Unexpected message received from stream. Expected: `%s`, found: `%s`", expected, message)
	}

	l.SetMetadata(map[string]string{"custom": "field"})
	l.Errorf("this log should have metadata")

	message = <-stream.out
	expected = `{"tag":"TestLoggerStreamsProperLogDataForError","message":"this log should have metadata","metadata":{"custom":"field"}}`
	if message != expected {
		t.Fatalf("Unexpected message received from stream. Expected: `%s`, found: `%s`", expected, message)
	}
}

// TestLoggerStreamsProperLogDataForFatal verifies that the logger
// sends all information with and without metadata for Fatal stream.
func TestLoggerStreamsProperLogDataForFatal(t *testing.T) {
	l := NewLogger("TestLoggerStreamsProperLogDataForFatal", false)

	stream := newMockStream(1)
	l.RegisterStream(Fatal, stream)

	l.Fatal("Hi there")

	message := <-stream.out
	expected := `{"tag":"TestLoggerStreamsProperLogDataForFatal","message":"Hi there"}`
	if message != expected {
		t.Fatalf("Unexpected message received from stream. Expected: `%s`, found: `%s`", expected, message)
	}

	l.SetMetadata(map[string]string{"custom": "field"})
	l.Fatalf("this log should have metadata")

	message = <-stream.out
	expected = `{"tag":"TestLoggerStreamsProperLogDataForFatal","message":"this log should have metadata","metadata":{"custom":"field"}}`
	if message != expected {
		t.Fatalf("Unexpected message received from stream. Expected: `%s`, found: `%s`", expected, message)
	}
}
