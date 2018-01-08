package gonyan

import (
	"testing"
)

func TestStreamManagerRegisterForInvalidLogLevel(t *testing.T) {
	manager := NewStreamManager()
	stream := newMockStream(1)

	if err := manager.Register(LogLevel(99999), stream); err == nil {
		t.Fatalf("Expected error for invalid loglevel. Found nil instead.")
	}
}

// TestStreamManagerRegisterForDebugLevel will try to register a stream
// for the specific log level and verify it has been properly registered.
func TestStreamManagerRegisterForDebugLevel(t *testing.T) {
	manager := NewStreamManager()
	stream := newMockStream(1)

	if err := manager.Register(Debug, stream); err != nil {
		t.Fatalf("Unexpected error: %s", err.Error())
	}
	if len(manager.streams[Debug]) != 1 {
		t.Fatalf("Unexpected Debug streams len. Expected: %d - Found: %d.", 1, len(manager.streams[Debug]))
	}
	if len(manager.streams[Verbose]) != 0 {
		t.Fatalf("Unexpected Verbose streams len. Expected: %d - Found: %d.", 0, len(manager.streams[Verbose]))
	}
	if len(manager.streams[Info]) != 0 {
		t.Fatalf("Unexpected Info streams len. Expected: %d - Found: %d.", 0, len(manager.streams[Info]))
	}
	if len(manager.streams[Warning]) != 0 {
		t.Fatalf("Unexpected Warning streams len. Expected: %d - Found: %d.", 0, len(manager.streams[Warning]))
	}
	if len(manager.streams[Error]) != 0 {
		t.Fatalf("Unexpected Error streams len. Expected: %d - Found: %d.", 0, len(manager.streams[Error]))
	}
	if len(manager.streams[Fatal]) != 0 {
		t.Fatalf("Unexpected Fatal streams len. Expected: %d - Found: %d.", 0, len(manager.streams[Fatal]))
	}
	if len(manager.streams[Panic]) != 0 {
		t.Fatalf("Unexpected Panic streams len. Expected: %d - Found: %d.", 1, len(manager.streams[Panic]))
	}
}

// TestStreamManagerRegisterForVerboseLevel will try to register a stream
// for the specific log level and verify it has been properly registered.
func TestStreamManagerRegisterForVerboseLevel(t *testing.T) {
	manager := NewStreamManager()
	stream := newMockStream(1)

	if err := manager.Register(Verbose, stream); err != nil {
		t.Fatalf("Unexpected error: %s", err.Error())
	}
	if len(manager.streams[Debug]) != 0 {
		t.Fatalf("Unexpected Debug streams len. Expected: %d - Found: %d.", 0, len(manager.streams[Debug]))
	}
	if len(manager.streams[Verbose]) != 1 {
		t.Fatalf("Unexpected Verbose streams len. Expected: %d - Found: %d.", 1, len(manager.streams[Verbose]))
	}
	if len(manager.streams[Info]) != 0 {
		t.Fatalf("Unexpected Info streams len. Expected: %d - Found: %d.", 0, len(manager.streams[Info]))
	}
	if len(manager.streams[Warning]) != 0 {
		t.Fatalf("Unexpected Warning streams len. Expected: %d - Found: %d.", 0, len(manager.streams[Warning]))
	}
	if len(manager.streams[Error]) != 0 {
		t.Fatalf("Unexpected Error streams len. Expected: %d - Found: %d.", 0, len(manager.streams[Error]))
	}
	if len(manager.streams[Fatal]) != 0 {
		t.Fatalf("Unexpected Fatal streams len. Expected: %d - Found: %d.", 0, len(manager.streams[Fatal]))
	}
	if len(manager.streams[Panic]) != 0 {
		t.Fatalf("Unexpected Panic streams len. Expected: %d - Found: %d.", 1, len(manager.streams[Panic]))
	}
}

// TestStreamManagerRegisterForInfoLevel will try to register a stream
// for the specific log level and verify it has been properly registered.
func TestStreamManagerRegisterForInfoLevel(t *testing.T) {
	manager := NewStreamManager()
	stream := newMockStream(1)

	if err := manager.Register(Info, stream); err != nil {
		t.Fatalf("Unexpected error: %s", err.Error())
	}
	if len(manager.streams[Debug]) != 0 {
		t.Fatalf("Unexpected Debug streams len. Expected: %d - Found: %d.", 0, len(manager.streams[Debug]))
	}
	if len(manager.streams[Verbose]) != 0 {
		t.Fatalf("Unexpected Verbose streams len. Expected: %d - Found: %d.", 0, len(manager.streams[Verbose]))
	}
	if len(manager.streams[Info]) != 1 {
		t.Fatalf("Unexpected Info streams len. Expected: %d - Found: %d.", 1, len(manager.streams[Info]))
	}
	if len(manager.streams[Warning]) != 0 {
		t.Fatalf("Unexpected Warning streams len. Expected: %d - Found: %d.", 0, len(manager.streams[Warning]))
	}
	if len(manager.streams[Error]) != 0 {
		t.Fatalf("Unexpected Error streams len. Expected: %d - Found: %d.", 0, len(manager.streams[Error]))
	}
	if len(manager.streams[Fatal]) != 0 {
		t.Fatalf("Unexpected Fatal streams len. Expected: %d - Found: %d.", 0, len(manager.streams[Fatal]))
	}
	if len(manager.streams[Panic]) != 0 {
		t.Fatalf("Unexpected Panic streams len. Expected: %d - Found: %d.", 1, len(manager.streams[Panic]))
	}
}

// TestStreamManagerRegisterForWarningLevel will try to register a stream
// for the specific log level and verify it has been properly registered.
func TestStreamManagerRegisterForWarningLevel(t *testing.T) {
	manager := NewStreamManager()
	stream := newMockStream(1)

	if err := manager.Register(Warning, stream); err != nil {
		t.Fatalf("Unexpected error: %s", err.Error())
	}
	if len(manager.streams[Debug]) != 0 {
		t.Fatalf("Unexpected Debug streams len. Expected: %d - Found: %d.", 0, len(manager.streams[Debug]))
	}
	if len(manager.streams[Verbose]) != 0 {
		t.Fatalf("Unexpected Verbose streams len. Expected: %d - Found: %d.", 0, len(manager.streams[Verbose]))
	}
	if len(manager.streams[Info]) != 0 {
		t.Fatalf("Unexpected Info streams len. Expected: %d - Found: %d.", 0, len(manager.streams[Info]))
	}
	if len(manager.streams[Warning]) != 1 {
		t.Fatalf("Unexpected Warning streams len. Expected: %d - Found: %d.", 1, len(manager.streams[Warning]))
	}
	if len(manager.streams[Error]) != 0 {
		t.Fatalf("Unexpected Error streams len. Expected: %d - Found: %d.", 0, len(manager.streams[Error]))
	}
	if len(manager.streams[Fatal]) != 0 {
		t.Fatalf("Unexpected Fatal streams len. Expected: %d - Found: %d.", 0, len(manager.streams[Fatal]))
	}
	if len(manager.streams[Panic]) != 0 {
		t.Fatalf("Unexpected Panic streams len. Expected: %d - Found: %d.", 1, len(manager.streams[Panic]))
	}
}

// TestStreamManagerRegisterForErrorLevel will try to register a stream
// for the specific log level and verify it has been properly registered.
func TestStreamManagerRegisterForErrorLevel(t *testing.T) {
	manager := NewStreamManager()
	stream := newMockStream(1)

	if err := manager.Register(Error, stream); err != nil {
		t.Fatalf("Unexpected error: %s", err.Error())
	}
	if len(manager.streams[Debug]) != 0 {
		t.Fatalf("Unexpected Debug streams len. Expected: %d - Found: %d.", 0, len(manager.streams[Debug]))
	}
	if len(manager.streams[Verbose]) != 0 {
		t.Fatalf("Unexpected Verbose streams len. Expected: %d - Found: %d.", 0, len(manager.streams[Verbose]))
	}
	if len(manager.streams[Info]) != 0 {
		t.Fatalf("Unexpected Info streams len. Expected: %d - Found: %d.", 0, len(manager.streams[Info]))
	}
	if len(manager.streams[Warning]) != 0 {
		t.Fatalf("Unexpected Warning streams len. Expected: %d - Found: %d.", 0, len(manager.streams[Warning]))
	}
	if len(manager.streams[Error]) != 1 {
		t.Fatalf("Unexpected Error streams len. Expected: %d - Found: %d.", 1, len(manager.streams[Error]))
	}
	if len(manager.streams[Fatal]) != 0 {
		t.Fatalf("Unexpected Fatal streams len. Expected: %d - Found: %d.", 0, len(manager.streams[Fatal]))
	}
	if len(manager.streams[Panic]) != 0 {
		t.Fatalf("Unexpected Panic streams len. Expected: %d - Found: %d.", 1, len(manager.streams[Panic]))
	}
}

// TestStreamManagerRegisterFoFatalLevel will try to register a stream
// for the specific log level and verify it has been properly registered.
func TestStreamManagerRegisterFoFatalLevel(t *testing.T) {
	manager := NewStreamManager()
	stream := newMockStream(1)

	if err := manager.Register(Fatal, stream); err != nil {
		t.Fatalf("Unexpected error: %s", err.Error())
	}
	if len(manager.streams[Debug]) != 0 {
		t.Fatalf("Unexpected Debug streams len. Expected: %d - Found: %d.", 0, len(manager.streams[Debug]))
	}
	if len(manager.streams[Verbose]) != 0 {
		t.Fatalf("Unexpected Verbose streams len. Expected: %d - Found: %d.", 0, len(manager.streams[Verbose]))
	}
	if len(manager.streams[Info]) != 0 {
		t.Fatalf("Unexpected Info streams len. Expected: %d - Found: %d.", 0, len(manager.streams[Info]))
	}
	if len(manager.streams[Warning]) != 0 {
		t.Fatalf("Unexpected Warning streams len. Expected: %d - Found: %d.", 0, len(manager.streams[Warning]))
	}
	if len(manager.streams[Error]) != 0 {
		t.Fatalf("Unexpected Error streams len. Expected: %d - Found: %d.", 0, len(manager.streams[Error]))
	}
	if len(manager.streams[Fatal]) != 1 {
		t.Fatalf("Unexpected Fatal streams len. Expected: %d - Found: %d.", 1, len(manager.streams[Fatal]))
	}
	if len(manager.streams[Panic]) != 0 {
		t.Fatalf("Unexpected Panic streams len. Expected: %d - Found: %d.", 1, len(manager.streams[Panic]))
	}
}

// TestStreamManagerRegisterForPanicLevel will try to register a stream
// for the specific log level and verify it has been properly registered.
func TestStreamManagerRegisterForPanicLevel(t *testing.T) {
	manager := NewStreamManager()
	stream := newMockStream(1)

	if err := manager.Register(Panic, stream); err != nil {
		t.Fatalf("Unexpected error: %s", err.Error())
	}
	if len(manager.streams[Debug]) != 0 {
		t.Fatalf("Unexpected Debug streams len. Expected: %d - Found: %d.", 0, len(manager.streams[Debug]))
	}
	if len(manager.streams[Verbose]) != 0 {
		t.Fatalf("Unexpected Verbose streams len. Expected: %d - Found: %d.", 0, len(manager.streams[Verbose]))
	}
	if len(manager.streams[Info]) != 0 {
		t.Fatalf("Unexpected Info streams len. Expected: %d - Found: %d.", 0, len(manager.streams[Info]))
	}
	if len(manager.streams[Warning]) != 0 {
		t.Fatalf("Unexpected Warning streams len. Expected: %d - Found: %d.", 0, len(manager.streams[Warning]))
	}
	if len(manager.streams[Error]) != 0 {
		t.Fatalf("Unexpected Error streams len. Expected: %d - Found: %d.", 0, len(manager.streams[Error]))
	}
	if len(manager.streams[Fatal]) != 0 {
		t.Fatalf("Unexpected Fatal streams len. Expected: %d - Found: %d.", 0, len(manager.streams[Fatal]))
	}
	if len(manager.streams[Panic]) != 1 {
		t.Fatalf("Unexpected Panic streams len. Expected: %d - Found: %d.", 1, len(manager.streams[Panic]))
	}
}

func TestStreamManagerSend(t *testing.T) {
	manager := NewStreamManager()
	stream := newMockStream(1)

	sampleMessage := NewLogMessage("TestSend", 0, "the-message", nil)

	if err := manager.Send(Debug, nil); err == nil {
		t.Fatalf("Expected error for invalid nil message. Found nil instead.")
	}
	if err := manager.Send(LogLevel(9999), sampleMessage); err == nil {
		t.Fatalf("Expected error for invalid log level. Found nil instead.")
	}

	if err := manager.Register(Debug, stream); err != nil {
		t.Fatalf("Unexpected error: %s", err.Error())
	}
	if err := manager.Send(Debug, sampleMessage); err != nil {
		t.Fatalf("Unexpected error: %s", err.Error())
	}
}
