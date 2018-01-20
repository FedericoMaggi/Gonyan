package gonyan

import (
	"fmt"
	"sync"
)

// BufferedStream represents a wrapper over a standard Stream which is intended
// to allow log buffering before transmission. There're two conditions that can
// trigger a buffer data transmission:
//
// 	- time interval: the buffer is emptied (by default) once every minute; The
// 		interval is configurable using the `SetInterval` function; You can
// 		disable this type of transmission by passing 0 to the function.
// 	- limit reached: the buffer memory is capped (by default) to 100 elements;
// 		reached the limit all stored logs are transmitted. The limit is
// 		configurable by using the `SetBufferLimit` function. You can disable
// 		this type of transmission by passing 0 to the function.
type BufferedStream struct {
	// Stream is the target, for this BufferedStream, where logs are sent.
	Stream *Stream
	// buffer for logs allocation it's emptied everytime logs are transmitted
	// through the provided stream.
	buffer [][]byte
	// bufferCount keeps track of the current number of items in the buffer.
	bufferCount int
	// bufferMutex is a mutex used for accessing the buffer and its capacity
	// counter.
	bufferMutex sync.Mutex
	// limit is used to perform an automatic data transmission when the
	// internal buffer reaches the desired limit, this feature is optional
	// and by default is disabled. Use SetBufferLimit to set the limit and
	// activate it.
	limit int
	// separator is used when flattening the buffer into a single dimensional
	// blob of messages, by default is `\n` but can be whatever you expect it
	// to be on the Stream receiver.
	separator byte
}

// DefaultPreallocatedBufferSize defines the default buffer size at start.
const DefaultPreallocatedBufferSize = 100

// DefaultFlatByteSliceSeparator defines the default value for the flat buffer
// separator.
const DefaultFlatByteSliceSeparator = '\n'

// NewBufferedStream creates a new BufferedSteam using provided stream.
func NewBufferedStream(stream *Stream) *BufferedStream {
	return &BufferedStream{
		limit:     0,
		buffer:    make([][]byte, DefaultPreallocatedBufferSize),
		Stream:    stream,
		separator: DefaultFlatByteSliceSeparator,
	}
}

// SetBufferLimit sets a custom cap limit to the buffer, when this cap is
// reached the buffer gets automatically emptied.
// Accepted values are non negative integers, negative values are ignored
// while setting the limit to 0 disables the feature.
func (b *BufferedStream) SetBufferLimit(bufferLimit int) {
	if bufferLimit >= 0 {
		b.limit = bufferLimit
	}
}

// SetFlatBufferSeparator allows you to define a custom flat buffer separator
// byte, bu default it is set to `\n` but can be whatever you expect it to be
// on your Stream receiver.
// Be careful to choose a proper byte since it will be used to split messages
// when flattening the buffer into a one-dimensional blob.
func (b *BufferedStream) SetFlatBufferSeparator(separator byte) {
	b.separator = separator
}

// Write will store provided log into the buffer prior transmission. If the log
// makes the buffer full it will fire the log transmission to the stream.
func (b *BufferedStream) Write(message []byte) (int, error) {
	var bufferForTransmission [][]byte
	var numberOfMessages int

	b.bufferMutex.Lock()
	// Limit has been reached, substitute the buffer and prepare the old one
	// for transmission.
	if b.limit > 0 && b.bufferCount >= b.limit {
		bufferForTransmission = b.buffer
		numberOfMessages = b.bufferCount

		// Create new buffer the same size of the old one.
		b.buffer = make([][]byte, len(b.buffer))
		b.bufferCount = 0
	}

	// Set the message in the buffer and then increment the position counter.
	b.buffer[b.bufferCount] = message
	b.bufferCount++
	b.bufferMutex.Unlock()

	// If the buffer was full fire a transmission with provided data.
	if bufferForTransmission != nil {
		go b.fireTransmission(bufferForTransmission, numberOfMessages)
	}

	return b.bufferCount, nil
}

// fireTransmission receives the messages slice to be transmitted on the Stream
// and writes it after flattening operation with provided optional separator
// byte (by default: `\n`).
func (b *BufferedStream) fireTransmission(messages [][]byte, amount int) {
	if b.Stream == nil {
		return
	}

	flatBuffer := flatten(messages, b.separator)
	if _, err := (*b.Stream).Write(flatBuffer); err != nil {
		fmt.Printf("[Gonyan]")
	}
}

// flatten utility function can be used to flatten a two dimensional byte slice
// in a single dimensional, to divide the original slices a separator byte can
// be provided.
func flatten(matrix [][]byte, separator byte) []byte {
	flat := []byte{}
	for _, row := range matrix {
		flat = append(flat, row...)
		flat = append(flat, separator)
	}
	return flat
}
