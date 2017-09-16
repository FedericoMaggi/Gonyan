package gonyan

import (
	"fmt"
)

// StreamManager wraps up all supported stream types.
type StreamManager struct {
	debugStreams []Stream
	verboseStreams []Stream
	infoStreams []Stream
	warningStreams []Stream
	errorStreams []Stream
	fatalStreams []Stream
}

// NewStreamManager creates a new, properly initialised, StreamManager instance.
func NewStreamManager() *StreamManager {
	return &StreamManager{
		debugStreams: make([]Stream,0),
		verboseStreams: make([]Stream,0),
		infoStreams: make([]Stream,0),
		warningStreams: make([]Stream,0),
		errorStreams: make([]Stream,0),
		fatalStreams: make([]Stream,0),
	}
}

// Register internally saves provided stream into proper container.
func (s *StreamManager) Register(level LogLevel, stream Stream) error {
	switch level {
	case Debug:
		s.debugStreams = append(s.debugStreams, stream)
	case Verbose:
		s.verboseStreams = append(s.verboseStreams, stream)
	case Info:
		s.infoStreams = append(s.infoStreams, stream)
	case Warning:
		s.warningStreams = append(s.warningStreams, stream)
	case Error:
		s.errorStreams = append(s.errorStreams, stream)
	case Fatal:
		s.fatalStreams = append(s.fatalStreams, stream)
	default: 
		return fmt.Errorf("invalid log level provided")
	}
	return nil
}

// Send fires stream writes for provided LogMessage into proper streams.
func (s *StreamManager) Send(level LogLevel, message *LogMessage) error {
	if message == nil {
		return fmt.Errorf("invalid nil message")
	}
	
	messageString, err := message.Serialise()
	if err != nil {
		return fmt.Errorf("serialisation error: %s", err.Error())
	}

	switch level {
	case Debug:
		s.sendToDebugStreams(messageString)
	case Verbose:
		s.sendToVerboseStreams(messageString)
	case Info:
		s.sendToInfoStreams(messageString)
	case Warning:
		s.sendToWarningStreams(messageString)
	case Error:
		s.sendToErrorStreams(messageString)
	case Fatal:
		s.sendToFatalStreams(messageString)

	default: 
		return fmt.Errorf("invalid log level provided")
	}
	return nil
}

// sendToDebugStreams fires writes for all registered Debug streams.
func (s *StreamManager) sendToDebugStreams(message string) error {
	for i := 0; i < len(s.debugStreams); i++ {
		s.debugStreams[i].Write(message)
	}
	return nil
}

// sendToVerboseStreams fires writes for all registered Verbose streams.
func (s *StreamManager) sendToVerboseStreams(message string) error {
	for i := 0; i < len(s.verboseStreams); i++ {
		s.verboseStreams[i].Write(message)
	}
	return nil
}

// sendToInfoStreams fires writes for all registered Info streams.
func (s *StreamManager) sendToInfoStreams(message string) error {
	for i := 0; i < len(s.infoStreams); i++ {
		s.infoStreams[i].Write(message)
	}
	return nil
}

// sendToWarningStreams fires writes for all registered Warning streams.
func (s *StreamManager) sendToWarningStreams(message string) error {
	for i := 0; i < len(s.warningStreams); i++ {
		s.warningStreams[i].Write(message)
	}
	return nil
}

// sendToErrorStreams fires writes for all registered Error streams.
func (s *StreamManager) sendToErrorStreams(message string) error {
	for i := 0; i < len(s.errorStreams); i++ {
		s.errorStreams[i].Write(message)
	}
	return nil
}

// sendToFatalStreams fires writes for all registered Fatal streams.
func (s *StreamManager) sendToFatalStreams(message string) error {
	for i := 0; i < len(s.fatalStreams); i++ {
		(s.fatalStreams[i]).Write(message)
	}
	return nil
}
