package gonyan

import (
	"fmt"
	"sync"
	"time"
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
	Stream Stream
	// buffer for logs allocation it's emptied everytime logs are transmitted
	// through the provided stream.
	buffer [][]byte
	// maxSize corresponds to the maximum size reachable by the buffer, this to
	// limit greedy memory allocation.
	initialSize int
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
	// scheduleInterval defines the time ticketing for automatic data
	// transmission. This feature is optional and by default is disabled, use
	// SetSchedulingInterval to set the time ticking and activate it.
	scheduleInteval time.Duration
	// routineRunning keeps track of the routine automatic transmission status.
	routineRunning bool
	// routineMutex is used to ensure atomic access on the routineRunning bool.
	routineMutex sync.Mutex
	// separator is used when flattening the buffer into a single dimensional
	// blob of messages, by default is `\n` but can be whatever you expect it
	// to be on the Stream receiver.
	separator byte
	// fatal is an optional function pointer used when something bad appens in
	// the buffered stream.
	fatal func(error)
}

// DefaultPreallocatedBufferSize defines the default buffer size at start.
const DefaultPreallocatedBufferSize = 100

// DefaultFlatByteSliceSeparator defines the default value for the flat buffer
// separator.
const DefaultFlatByteSliceSeparator = '\n'

// NewBufferedStream creates a new BufferedSteam using provided stream.
func NewBufferedStream(stream Stream) *BufferedStream {
	return &BufferedStream{
		limit:           0,
		buffer:          make([][]byte, DefaultPreallocatedBufferSize),
		initialSize:     DefaultPreallocatedBufferSize,
		Stream:          stream,
		scheduleInteval: 0,
		separator:       DefaultFlatByteSliceSeparator,
		fatal: func(err error) {
			fmt.Printf("[Gonyan] [BufferedSteam] [Fatal] %s.\n", err.Error())
		},
	}
}

// SetBufferLimit sets a custom cap limit to the buffer, when this cap is
// reached the buffer gets automatically emptied.
//
// NOTE: Accepted values are non negative integers, negative values are ignored
// while setting the limit to 0 disables the feature.
func (b *BufferedStream) SetBufferLimit(bufferLimit int) {
	if bufferLimit < 0 {
		return
	}
	b.limit = bufferLimit
}

// SetStartingSize sets the initial size of the buffer, note that the buffer
// *will* be appended with new messages if it reaches the set size but when it
// gets flushed and recreated it will be allocated with provided size.
// Please note that this function will trigger a buffer flush a flag can be
// provided in order to decide whether to ignore the flush and lose the old
// buffer content or send it to the stream instead.
// The function returns a boolean flag indicating whether the buffer has been
// resized, and an error that indicates whether there was a failure during data
// transmission through the Stream.
//
// NOTE: Accepted values are non negative integers, negative values are ignored
// also setting the initSize to a value equal to the current one has no effect.
func (b *BufferedStream) SetStartingSize(initSize int, send bool) (bool, error) {
	if initSize < 0 || initSize == b.initialSize {
		return false, nil
	}

	b.initialSize = initSize

	b.bufferMutex.Lock()
	oldBuffer, oldCount := b.flush()
	b.bufferMutex.Unlock()

	if !send {
		return true, nil
	}

	if oldBuffer != nil && oldCount != 0 {
		if err := b.fireTransmission(oldBuffer, oldCount); err != nil {
			return true, err
		}
	}

	return true, nil
}

// SetSchedulingInterval will set the internal scheduling interval for
// automatic data transmission. The function will start the routine if
// provided flag is true, otherwise you will have to start the routine
// by yourself.
// Passing a non-positive value to the function will disable the feature
// and stop the routine.
func (b *BufferedStream) SetSchedulingInterval(tick time.Duration, start bool) error {
	if tick <= 0 {
		b.scheduleInteval = 0
		b.routineMutex.Lock()
		b.routineRunning = false
		b.routineMutex.Unlock()
		return nil
	}
	b.scheduleInteval = tick

	if !start {
		return nil
	}

	// TODO: Consider directly returning since this error will never be
	// returned because we previously make sure that used data are correct.
	if err := b.StartAutonomousTransmission(); err != nil {
		return fmt.Errorf("routine start failed: %s", err.Error())
	}
	return nil
}

// StartAutonomousTransmission starts the routine transmission in a concurrent
// goroutine, please note that if no scheduling interval has been set then no
// routine will be started and an error will be returned.
func (b *BufferedStream) StartAutonomousTransmission() error {
	b.routineMutex.Lock()
	routineRunning := b.routineRunning
	b.routineMutex.Unlock()

	if routineRunning {
		return fmt.Errorf("routine already running")
	}
	if b.scheduleInteval == 0 {
		return fmt.Errorf("no scheduling interval set")
	}

	ticker := time.NewTicker(b.scheduleInteval)
	go func(ticker *time.Ticker) {
		if err := b.autonomousTranmissionRoutine(ticker); err != nil {
			// TODO: Consider using a different communication channel.
			if b.fatal != nil {
				b.fatal(fmt.Errorf("routine has been stopped due to error: %s", err.Error()))
			}
		}
	}(ticker)

	return nil
}

// StopAutonomousTransmission sets the boolean state to false causing the next
// routine iteration to be immediately interrupted.
func (b *BufferedStream) StopAutonomousTransmission() {
	b.routineMutex.Lock()
	b.routineRunning = false
	b.routineMutex.Unlock()
}

func (b *BufferedStream) autonomousTranmissionRoutine(ticker *time.Ticker) error {
	b.routineMutex.Lock()
	b.routineRunning = true
	b.routineMutex.Unlock()

	for range ticker.C {
		b.routineMutex.Lock()
		routineRunning := b.routineRunning
		b.routineMutex.Unlock()
		if routineRunning == false {
			// Stop the ticker and exit.
			ticker.Stop()
			// return fmt.Errorf("routine has been stopped")
			return nil
		}

		// Let's flush the old buffer and send it through the stream and
		// if the buffer was full fire a transmission with retrieved data.
		b.bufferMutex.Lock()
		oldBuffer, oldSize := b.flush()
		b.bufferMutex.Unlock()
		if oldBuffer != nil && oldSize != 0 {
			if err := b.fireTransmission(oldBuffer, oldSize); err != nil {
				return fmt.Errorf("gonyan buffered stream failure during data transmission: %s", err.Error())
			}
		}
	}
	return nil
}

// SetFatalFn sets the optional function for fatal error signals.
func (b *BufferedStream) SetFatalFn(fatalFn func(error)) {
	b.fatal = fatalFn
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
	var oldBuffer [][]byte
	var oldSize int

	b.bufferMutex.Lock()

	// If capped transmission is enabled and the set limit has been reached
	// then substitute the buffer and prepare the old one for transmission.
	if b.limit > 0 && b.bufferCount >= b.limit {
		oldBuffer, oldSize = b.flush()
	}

	// If the buffer count has reached the total length of the buffer then we
	// need a new slot for the received message, allocate it empty.
	if b.bufferCount >= len(b.buffer) {
		b.buffer = append(b.buffer, []byte{})
	}

	// Set the message in the buffer and then increment the position counter.
	b.buffer[b.bufferCount] = message
	b.bufferCount++
	// Concurrently safe read the value to be returned.
	newCount := b.bufferCount
	b.bufferMutex.Unlock()

	// If the buffer was full fire a transmission with provided data.
	if oldBuffer != nil && oldSize != 0 {
		go func(buffer [][]byte, size int) {
			if err := b.fireTransmission(buffer, size); err != nil {
				if b.fatal != nil {
					b.fatal(fmt.Errorf("gonyan buffered stream failure during data transmission: %s", err.Error()))
				}
			}
		}(oldBuffer, oldSize)
	}

	return newCount, nil
}

// fireTransmission receives the messages slice to be transmitted on the Stream
// and writes it after flattening operation with provided optional separator
// byte (by default: `\n`).
func (b *BufferedStream) fireTransmission(messages [][]byte, amount int) error {
	if b.Stream == nil {
		return fmt.Errorf("invalid stream found")
	}

	flatBuffer := flatten(messages, b.separator)
	if n, err := b.Stream.Write(flatBuffer); err != nil {
		return fmt.Errorf("failed write on stream: %s, returned count: %d", err.Error(), n)
	}
	return nil
}

// flush will return the old buffer and its count and allocate a new empty
// buffer to be used.
// This is not cuncurrent safe by its own and *must* be called after a lock
// has been set.
func (b *BufferedStream) flush() ([][]byte, int) {
	oldBuffer := b.buffer
	oldBufferSize := b.bufferCount

	// Create new buffer using set initial size.
	b.buffer = make([][]byte, b.initialSize)
	b.bufferCount = 0

	return oldBuffer, oldBufferSize
}

// flatten utility function can be used to flatten a two dimensional byte slice
// in a single dimensional, to divide the original slices a separator byte can
// be provided.
func flatten(matrix [][]byte, separator byte) []byte {
	flat := []byte{}
	for i, row := range matrix {
		if len(row) == 0 || i+1 > len(matrix) {
			continue
		}
		flat = append(flat, row...)

		// Avoid appending the separator after last message.
		if i+1 != len(matrix) && matrix[i+1] != nil {
			flat = append(flat, separator)
		}
	}
	return flat
}
