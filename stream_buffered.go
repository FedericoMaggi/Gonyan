package gonyan

// BufferedStream represents a wrapper over a standard Stream which is intended
// to allow logs buffering prior transmission. There're two conditions that can
// met a data transmission:
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
	// limit is used to perform an automatic data transmission when the
	// internal buffer reaches the desired limit, this feature is optional
	// and by default is disabled. Use SetBufferLimit to set the limit and
	// activate it.
	limit int
}

// PreallocatedBufferSize defines the default buffer size at start.
const PreallocatedBufferSize = 100

// NewBufferedStream creates a new buffer for provided stream.
func NewBufferedStream(stream *Stream) *BufferedStream {
	return &BufferedStream{
		limit:  PreallocatedBufferSize,
		buffer: make([][]byte, PreallocatedBufferSize),
		Stream: stream,
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

func (b *BufferedStream) Write([]byte) (int, error) {
	return 0, nil
}
