package gonyan

import (
	"fmt"
)

// StreamManager wraps up all supported stream types.
type StreamManager struct {
	streams map[LogLevel][]Stream
}

// NewStreamManager creates a new, properly initialised, StreamManager instance.
func NewStreamManager() *StreamManager {
	s := &StreamManager{}
	s.streams = make(map[LogLevel][]Stream)
	s.streams[Debug] = make([]Stream, 0)
	s.streams[Verbose] = make([]Stream, 0)
	s.streams[Info] = make([]Stream, 0)
	s.streams[Warning] = make([]Stream, 0)
	s.streams[Error] = make([]Stream, 0)
	s.streams[Fatal] = make([]Stream, 0)
	s.streams[Panic] = make([]Stream, 0)

	return s
}

// Register internally saves provided stream into proper container.
func (s *StreamManager) Register(level LogLevel, stream Stream) error {
	registeredStreams, ok := s.streams[level]
	if !ok {
		return fmt.Errorf("invalid log level provided")
	}

	registeredStreams = append(registeredStreams, stream)
	s.streams[level] = registeredStreams
	return nil
}

// Send fires stream writes for provided LogMessage into proper streams.
func (s *StreamManager) Send(level LogLevel, message *LogMessage) error {
	if message == nil {
		return fmt.Errorf("invalid nil message")
	}

	registeredStreams, ok := s.streams[level]
	if !ok {
		return fmt.Errorf("invalid log level provided")
	}

	messageString, err := message.Serialise()
	if err != nil {
		return fmt.Errorf("serialisation error: %s", err.Error())
	}

	for i := 0; i < len(registeredStreams); i++ {
		registeredStreams[i].Write(messageString)
	}
	return nil
}
