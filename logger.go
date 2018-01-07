// Package gonyan is a stream-based logging library.
package gonyan

import (
	"fmt"
	"time"
)

// Logger represents an instance of a thread safe logger. It can hold different
// streams to log to and a few useful settings to customise the logs such as
// custom tag, metadata and timestamp.
type Logger struct {
	tag       string
	timestamp bool
	streams   *StreamManager
	metadata  map[string]string
}

// NewLogger creates a new logger instance with provided configuration.
func NewLogger(tag string, timestamp bool) *Logger {
	logger := &Logger{
		tag:       tag,
		timestamp: timestamp,
		streams:   NewStreamManager(),
	}
	return logger
}

// SetMetadata sets the optional metadata values for this logger.
// Metadata will be added to each log streamed from the logger instace.
func (l *Logger) SetMetadata(metadata map[string]string) {
	if metadata != nil {
		l.metadata = metadata
	}
}

// ClearMetadata clears the optional metadata from the logger instance.
func (l *Logger) ClearMetadata() {
	l.metadata = nil
}

// RegisterStream register provided stream associating it with provided level
// inside the interal StreamManager instance.
func (l *Logger) RegisterStream(level LogLevel, stream Stream) {
	l.streams.Register(level, stream)
}

// Debugf logs provided message into registered debug level streams.
// The function accepts a format and a variadic number of arguments
// to compose the final log data.
func (l *Logger) Debugf(format string, args ...interface{}) {
	l.Debug(fmt.Sprintf(format, args))
}

// Debug logs provided message into registered debug level streams.
func (l *Logger) Debug(message string) {
	l.Log(Debug, message)
}

// Verbosef logs provided message into registered verbose level streams.
// The function accepts a format and a variadic number of arguments
// to compose the final log data.
func (l *Logger) Verbosef(format string, args ...interface{}) {
	l.Verbose(fmt.Sprintf(format, args))
}

// Verbose logs provided message into registered verbose level streams.
func (l *Logger) Verbose(message string) {
	l.Log(Verbose, message)
}

// Infof logs provided message into info level streams.
// The function accepts a format and a variadic number of arguments
// to compose the final log data.
func (l *Logger) Infof(format string, args ...interface{}) {
	l.Info(fmt.Sprintf(format, args))
}

// Info logs provided message into info level streams.
func (l *Logger) Info(message string) {
	l.Log(Info, message)
}

// Warningf logs provided message into warning level streams.
// The function accepts a format and a variadic number of arguments
// to compose the final log data.
func (l *Logger) Warningf(format string, args ...interface{}) {
	l.Warning(fmt.Sprintf(format, args))
}

// Warning logs provided message into warning level streams.
func (l *Logger) Warning(message string) {
	l.Log(Warning, message)
}

// Errorf logs provided message into error level streams.
// The function accepts a format and a variadic number of arguments
// to compose the final log data.
func (l *Logger) Errorf(format string, args ...interface{}) {
	l.Error(fmt.Sprintf(format, args))
}

// Error logs provided message into error level streams.
func (l *Logger) Error(message string) {
	l.Log(Warning, message)
}

// Fatalf logs provided message into fatal level streams.
// The function accepts a format and a variadic number of arguments
// to compose the final log data.
func (l *Logger) Fatalf(format string, args ...interface{}) {
	l.Fatal(fmt.Sprintf(format, args))
}

// Fatal logs provided message into fatal level streams.
func (l *Logger) Fatal(message string) {
	l.Log(Fatal, message)
}

// Panicf logs provided message into panic level streams.
// The function accepts a format and a variadic number of arguments
// to compose the final log data.
// Note: When log is performed panic() is invoked.
func (l *Logger) Panicf(format string, args ...interface{}) {
	l.Fatal(fmt.Sprintf(format, args))
}

// Panic logs provided message into panic level streams.
// Note: When log is performed panic() is invoked.
func (l *Logger) Panic(message string) {
	l.Log(Fatal, message)
	panic(message)
}

// Logf logs provided message into the streams corresponding to provided level.
// The function accepts a format and a variadic number of arguments
// to compose the final log data.
func (l *Logger) Logf(level LogLevel, format string, args ...interface{}) {
	l.Logf(level, fmt.Sprintf(format, args))
}

// Log function builds the final JSON message and sends it to the correct streams.
func (l *Logger) Log(level LogLevel, message string) {
	var t int64
	if l.timestamp {
		t = time.Now().UTC().UnixNano()
	}

	m := NewLogMessage(l.tag, t, message, l.metadata)

	// Send message to streams via the StreamManager.
	if err := l.streams.Send(level, m); err != nil {
		fmt.Printf("[FATAL] [gonyan] Can't send log `%s` to stream `%s`", message, GetLevelLabel(level))
	}
}
