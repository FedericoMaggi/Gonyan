package gonyan

import (
	"os"
	"testing"
)

// TestNewLoggerWithStdout will test that the os.Stdout
// file is compatible with current Stream implementation.
// TODO: Capture stdout content to verify correct logging.
func TestNewLoggerWithStdout(t *testing.T) {
	l := NewLogger("X", nil, false)
	l.RegisterStream(Debug, os.Stdout)
	l.Debug("this is a debug log and should appear on stdout")
}

// TestNewLoggerWithStderr will test that the os.Stdout
// file is compatible with current Stream implementation.
// TODO: Capture stderr content to verify correct logging.
func TestNewLoggerWithStderr(t *testing.T) {
	l := NewLogger("X", nil, true)
	l.RegisterStream(Error, os.Stderr)
	l.Error("this is an error log")
}
