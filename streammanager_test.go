package gonyan

import (
	"testing"
	"time"
)

func TestStreamManagerRegister(t *testing.T) {
	manager := NewStreamManager()
	stream := newMockStream(3)

	if err := manager.Register(LogLevel(99999), stream); err == nil {
		t.Fatalf("Expected error for invalid loglevel. Found nil instead.")
	}


	if err := manager.Register(Debug, stream); err != nil {
		t.Fatalf("Unexpected error: %s", err.Error())
	}
	if len(manager.debugStreams) != 1 {
		t.Fatalf("Unexpected debugStreams len. Expected: %d - Found: %d.", 1, len(manager.debugStreams))
	}
	if len(manager.verboseStreams) != 0 {
		t.Fatalf("Unexpected verboseStreams len. Expected: %d - Found: %d.", 0, len(manager.verboseStreams))
	}
	if len(manager.infoStreams) != 0 {
		t.Fatalf("Unexpected infoStreams len. Expected: %d - Found: %d.", 0, len(manager.infoStreams))
	}
	if len(manager.warningStreams) != 0 {
		t.Fatalf("Unexpected warningStreams len. Expected: %d - Found: %d.", 0, len(manager.warningStreams))
	}
	if len(manager.errorStreams) != 0 {
		t.Fatalf("Unexpected errorStreams len. Expected: %d - Found: %d.", 0, len(manager.errorStreams))
	}
	if len(manager.fatalStreams) != 0 {
		t.Fatalf("Unexpected fatalStreams len. Expected: %d - Found: %d.", 0, len(manager.warningStreams))
	}

	if err := manager.Register(Verbose, stream); err != nil {
		t.Fatalf("Unexpected error: %s", err.Error())
	}
	if len(manager.debugStreams) != 1 {
		t.Fatalf("Unexpected debugStreams len. Expected: %d - Found: %d.", 1, len(manager.debugStreams))
	}
	if len(manager.verboseStreams) != 1 {
		t.Fatalf("Unexpected verboseStreams len. Expected: %d - Found: %d.", 1, len(manager.verboseStreams))
	}
	if len(manager.infoStreams) != 0 {
		t.Fatalf("Unexpected infoStreams len. Expected: %d - Found: %d.", 0, len(manager.infoStreams))
	}
	if len(manager.warningStreams) != 0 {
		t.Fatalf("Unexpected warningStreams len. Expected: %d - Found: %d.", 0, len(manager.warningStreams))
	}
	if len(manager.errorStreams) != 0 {
		t.Fatalf("Unexpected errorStreams len. Expected: %d - Found: %d.", 0, len(manager.errorStreams))
	}
	if len(manager.fatalStreams) != 0 {
		t.Fatalf("Unexpected fatalStreams len. Expected: %d - Found: %d.", 0, len(manager.warningStreams))
	}

	if err := manager.Register(Info, stream); err != nil {
		t.Fatalf("Unexpected error: %s", err.Error())
	}
	if len(manager.debugStreams) != 1 {
		t.Fatalf("Unexpected debugStreams len. Expected: %d - Found: %d.", 1, len(manager.debugStreams))
	}
	if len(manager.verboseStreams) != 1 {
		t.Fatalf("Unexpected verboseStreams len. Expected: %d - Found: %d.", 1, len(manager.verboseStreams))
	}
	if len(manager.infoStreams) != 1 {
		t.Fatalf("Unexpected infoStreams len. Expected: %d - Found: %d.", 1, len(manager.infoStreams))
	}
	if len(manager.warningStreams) != 0 {
		t.Fatalf("Unexpected warningStreams len. Expected: %d - Found: %d.", 0, len(manager.warningStreams))
	}
	if len(manager.errorStreams) != 0 {
		t.Fatalf("Unexpected errorStreams len. Expected: %d - Found: %d.", 0, len(manager.errorStreams))
	}
	if len(manager.fatalStreams) != 0 {
		t.Fatalf("Unexpected fatalStreams len. Expected: %d - Found: %d.", 0, len(manager.warningStreams))
	}

	if err := manager.Register(Warning, stream); err != nil {
		t.Fatalf("Unexpected error: %s", err.Error())
	}
	if len(manager.debugStreams) != 1 {
		t.Fatalf("Unexpected debugStreams len. Expected: %d - Found: %d.", 1, len(manager.debugStreams))
	}
	if len(manager.verboseStreams) != 1 {
		t.Fatalf("Unexpected verboseStreams len. Expected: %d - Found: %d.", 1, len(manager.verboseStreams))
	}
	if len(manager.infoStreams) != 1 {
		t.Fatalf("Unexpected infoStreams len. Expected: %d - Found: %d.", 1, len(manager.infoStreams))
	}
	if len(manager.warningStreams) != 1 {
		t.Fatalf("Unexpected warningStreams len. Expected: %d - Found: %d.", 1, len(manager.warningStreams))
	}
	if len(manager.errorStreams) != 0 {
		t.Fatalf("Unexpected errorStreams len. Expected: %d - Found: %d.", 0, len(manager.errorStreams))
	}
	if len(manager.fatalStreams) != 0 {
		t.Fatalf("Unexpected fatalStreams len. Expected: %d - Found: %d.", 0, len(manager.warningStreams))
	}

	if err := manager.Register(Error, stream); err != nil {
		t.Fatalf("Unexpected error: %s", err.Error())
	}
	if len(manager.debugStreams) != 1 {
		t.Fatalf("Unexpected debugStreams len. Expected: %d - Found: %d.", 1, len(manager.debugStreams))
	}
	if len(manager.verboseStreams) != 1 {
		t.Fatalf("Unexpected verboseStreams len. Expected: %d - Found: %d.", 1, len(manager.verboseStreams))
	}
	if len(manager.infoStreams) != 1 {
		t.Fatalf("Unexpected infoStreams len. Expected: %d - Found: %d.", 1, len(manager.infoStreams))
	}
	if len(manager.warningStreams) != 1 {
		t.Fatalf("Unexpected warningStreams len. Expected: %d - Found: %d.", 1, len(manager.warningStreams))
	}
	if len(manager.errorStreams) != 1 {
		t.Fatalf("Unexpected errorStreams len. Expected: %d - Found: %d.", 1, len(manager.errorStreams))
	}
	if len(manager.fatalStreams) != 0 {
		t.Fatalf("Unexpected fatalStreams len. Expected: %d - Found: %d.", 0, len(manager.warningStreams))
	}
	
	if err := manager.Register(Fatal, stream); err != nil {
		t.Fatalf("Unexpected error: %s", err.Error())
	}
	if len(manager.debugStreams) != 1 {
		t.Fatalf("Unexpected debugStreams len. Expected: %d - Found: %d.", 1, len(manager.debugStreams))
	}
	if len(manager.verboseStreams) != 1 {
		t.Fatalf("Unexpected verboseStreams len. Expected: %d - Found: %d.", 1, len(manager.verboseStreams))
	}
	if len(manager.infoStreams) != 1 {
		t.Fatalf("Unexpected infoStreams len. Expected: %d - Found: %d.", 1, len(manager.infoStreams))
	}
	if len(manager.warningStreams) != 1 {
		t.Fatalf("Unexpected warningStreams len. Expected: %d - Found: %d.", 1, len(manager.warningStreams))
	}
	if len(manager.errorStreams) != 1 {
		t.Fatalf("Unexpected errorStreams len. Expected: %d - Found: %d.", 1, len(manager.errorStreams))
	}
	if len(manager.fatalStreams) != 1 {
		t.Fatalf("Unexpected fatalStreams len. Expected: %d - Found: %d.", 1, len(manager.warningStreams))
	}
}

func TestStreamManagerSend(t *testing.T) {
	manager := NewStreamManager()
	stream := newMockStream(1)

	sampleMessage := NewLogMessage("TestSend", time.Time{}, "the-message")

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