package gonyan

import (
	"bytes"
	"fmt"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestNewBufferedStream(t *testing.T) {
	b := NewBufferedStream(nil)
	if b.limit != 0 {
		t.Fatalf("Unexpected default limit value. Expected: %d - found: %d", 0, b.limit)
	}
	if len(b.buffer) != DefaultPreallocatedBufferSize {
		t.Fatalf("Unexpected default buffer size. Expected: %d - Found: %d.", DefaultPreallocatedBufferSize, len(b.buffer))
	}
	if b.initialSize != DefaultPreallocatedBufferSize {
		t.Fatalf("Unexpected default initial buffer size. Expected: %d - Found: %d.", DefaultPreallocatedBufferSize, b.initialSize)
	}
	if b.separator != DefaultFlatByteSliceSeparator {
		t.Fatalf("Unexpected default separator. Expected: %c - Found: %c.", DefaultFlatByteSliceSeparator, b.separator)
	}
	if b.routineRunning {
		t.Fatalf("Unexpected boolean flag. Routine Running should be false.")
	}
	if b.scheduleInteval != 0 {
		t.Fatalf("Unexpected default schedule intervale. Expected: %d - Found: %d.", 0, b.scheduleInteval)
	}
}

func TestSetBufferLimit(t *testing.T) {
	b := NewBufferedStream(nil)
	if b.limit != 0 {
		t.Fatalf("Unexpected default limit value. Expected: %d - found: %d", 0, b.limit)
	}

	b.SetBufferLimit(10)
	if b.limit != 10 {
		t.Fatalf("Unexpected limit value. Expected: %d - found: %d", 10, b.limit)
	}

	b.SetBufferLimit(-1)
	if b.limit != 10 {
		t.Fatalf("Unexpected limit value. Expected: %d - found: %d", 10, b.limit)
	}
}

func TestSetStartingSize(t *testing.T) {
	type TestCase struct {
		size                int
		send                bool
		expectedSet         bool
		expectedError       bool
		expectedInitialSize int
	}

	testCases := []TestCase{
		{
			size: -1, send: true,
			expectedSet: false, expectedError: false,
			expectedInitialSize: DefaultPreallocatedBufferSize,
		},
		{
			size: DefaultPreallocatedBufferSize, send: true,
			expectedSet: false, expectedError: false,
			expectedInitialSize: DefaultPreallocatedBufferSize,
		},
		{
			size: 10, send: false,
			expectedSet: true, expectedError: false,
			expectedInitialSize: 10,
		},
		// Flag for transmission is true but the stream is not valid
		// still we do not expect an error since the buffer is empty.
		{
			size: 11, send: true,
			expectedSet: true, expectedError: false,
			expectedInitialSize: 11,
		},
	}

	b := NewBufferedStream(nil)
	if b.initialSize != DefaultPreallocatedBufferSize {
		t.Fatalf("Unexpected default initialSize value. Expected: %d - found: %d", DefaultPreallocatedBufferSize, b.initialSize)
	}

	for i, testCase := range testCases {
		set, err := b.SetStartingSize(testCase.size, testCase.send)
		if set != testCase.expectedSet {
			t.Fatalf("Case#%d - Unexpected set flag should be %t, found: %t.", i, testCase.expectedSet, set)
		}

		if testCase.expectedError && err == nil {
			t.Fatalf("Case#%d - An error was expected, found %v instead", i, err)
		}
		if !testCase.expectedError && err != nil {
			t.Fatalf("Case#%d - Unexpected error found: %s.", i, err.Error())
		}
		if b.initialSize != testCase.expectedInitialSize {
			t.Fatalf("Case#%d - Unexpected initial size found. Expected: %d - Found: %d.", i, testCase.expectedInitialSize, b.initialSize)
		}
	}
}

func TestBufferedStreamSetStartingSizeEdgeCases(t *testing.T) {
	b := NewBufferedStream(nil)
	// Flag for transmission is true but the stream is not valid and
	// we expect an error since we have sent a message to the stream.
	b.Write([]byte(`logging something nasty`))
	set, err := b.SetStartingSize(5, true)
	if !set {
		t.Fatalf("Unexpected set flag, should be true, found: %t.", set)
	}
	if err == nil {
		t.Fatalf("Unexpected nil error.")
	}
	t.Logf("This was an expected error: %s.", err.Error())
	if b.initialSize != 5 {
		t.Fatalf("Unexpected default initial buffer size. Expected: %d - Found: %d.", 11, b.initialSize)
	}
	if len(b.buffer) != 5 {
		t.Fatalf("Unexpected buffer len. Expected: %d - Found: %d.", 5, len(b.buffer))
	}

	s := newMockStream(5)
	b = NewBufferedStream(s)
	// Flag for transmission is true and the stream is now valid so
	// we do not expect an error.
	b.Write([]byte(`logging something nasty, again`))
	set, err = b.SetStartingSize(6, true)
	if !set {
		t.Fatalf("Unexpected set flag, should be true, found: %t.", set)
	}
	if err != nil {
		t.Fatalf("Unexpected error: %s.", err.Error())
	}
	if b.initialSize != 6 {
		t.Fatalf("Unexpected default initial buffer size. Expected: %d - Found: %d.", 11, b.initialSize)
	}
	if len(b.buffer) != 6 {
		t.Fatalf("Unexpected buffer len. Expected: %d - Found: %d.", 5, len(b.buffer))
	}

	received := <-s.out
	if received != "logging something nasty, again" {
		t.Fatalf("Unexpected received message. Expected: `%s` - Found: `%s`.", "logging something nasty, again", received)
	}
}

func TestLogMessageWithTrailingNewLine(t *testing.T) {

	s := newMockStream(5)
	b := NewBufferedStream(s)
	// Flag for transmission is true and the stream is now valid so
	// we do not expect an error.
	b.Write([]byte("logging something with newline\n"))
	set, err := b.SetStartingSize(6, true)
	if !set {
		t.Fatalf("Unexpected set flag, should be true, found: %t.", set)
	}
	if err != nil {
		t.Fatalf("Unexpected error: %s.", err.Error())
	}
	received := <-s.out
	// Note: the newline character here is the one actually inserted in the
	// log, not the separator added when flattening the buffer.
	if received != "logging something with newline\n" {
		t.Fatalf("Unexpected received message. Expected: `%s` - Found: `%s`.", "logging something with newline\n", received)
	}
}

func TestSetSchedulingInterval(t *testing.T) {
	b := NewBufferedStream(nil)

	// Mess up internal values.
	b.scheduleInteval = 100
	b.routineMutex.Lock()
	b.routineRunning = true
	b.routineMutex.Unlock()

	if err := b.SetSchedulingInterval(-100, false); err != nil {
		t.Fatalf("Unexpected error: %s.", err.Error())
	}
	if b.scheduleInteval != 0 {
		t.Fatalf("Unexpected scheduleInterval. Expected: %d - Found: %d.", 0, b.scheduleInteval)
	}
	if b.routineRunning {
		t.Fatalf("Unexpected routineRunning flag. Should be false!")
	}

	// Mess up internal values.
	b.scheduleInteval = 100
	b.routineMutex.Lock()
	b.routineRunning = true
	b.routineMutex.Unlock()
	if err := b.SetSchedulingInterval(0, false); err != nil {
		t.Fatalf("Unexpected error: %s.", err.Error())
	}
	if b.scheduleInteval != 0 {
		t.Fatalf("Unexpected scheduleInterval. Expected: %d - Found: %d.", 0, b.scheduleInteval)
	}
	if b.routineRunning {
		t.Fatalf("Unexpected routineRunning flag. Should be false!")
	}

	// Mess up internal values.
	b.scheduleInteval = 100
	b.routineMutex.Lock()
	b.routineRunning = true
	b.routineMutex.Unlock()
	if err := b.SetSchedulingInterval(0, true); err != nil {
		t.Fatalf("Unexpected error: %s.", err.Error())
	}
	if b.scheduleInteval != 0 {
		t.Fatalf("Unexpected scheduleInterval. Expected: %d - Found: %d.", 0, b.scheduleInteval)
	}
	if b.routineRunning {
		t.Fatalf("Unexpected routineRunning flag. Should be false!")
	}

	// Mess up internal values.
	b.scheduleInteval = 100
	if err := b.SetSchedulingInterval(10*time.Second, false); err != nil {
		t.Fatalf("Unexpected error: %s.", err.Error())
	}
	if b.scheduleInteval != 10*time.Second {
		t.Fatalf("Unexpected scheduleInterval. Expected: %d - Found: %d.", 10*time.Second, b.scheduleInteval)
	}
	if b.routineRunning {
		t.Fatalf("Unexpected routineRunning flag. Should be false!")
	}

	// Mess up internal values.
	b.scheduleInteval = 100
	if err := b.SetSchedulingInterval(1*time.Second, true); err != nil {
		t.Fatalf("Unexpected error: %s.", err.Error())
	}
	defer b.StopAutonomousTransmission()
	if b.scheduleInteval != 1*time.Second {
		t.Fatalf("Unexpected scheduleInterval. Expected: %d - Found: %d.", 1*time.Second, b.scheduleInteval)
	}
	time.Sleep(2 * time.Second) // Wait for the routine to be started.
	b.routineMutex.Lock()
	if !b.routineRunning {
		b.routineMutex.Unlock()
		t.Fatalf("Unexpected routineRunning flag. Should be true!")
	}
	b.routineMutex.Unlock()
}

func TestAutonomousTransmissionRoutine(t *testing.T) {
	mock := newMockStream(3)
	b := NewBufferedStream(mock)

	go func() {
		ticker := time.NewTicker(1 * time.Second)
		if err := b.autonomousTranmissionRoutine(ticker); err != nil {
			t.Fatalf("Unexpected error: %s", err.Error())
		}
	}()

	time.Sleep(2 * time.Second)

	n, err := b.Write([]byte("hej"))
	if err != nil {
		t.Fatalf("Unexpected error: %s", err.Error())
	}
	if n != 1 {
		t.Fatalf("Unexpected buffer count returned. Expected: %d - Found: %d.", 1, n)
	}

	time.Sleep(2 * time.Second)

	// Read from chan
	select {
	case message := <-mock.out:
		if message != string("hej") {
			t.Fatalf("Unexpected read message. Expected: %s - Found: %s.", "hej", string(message))
		}
	default:
		t.Fatalf("Failed read from mock stream.")
	}

	time.Sleep(1 * time.Second)

	n, err = b.Write([]byte("hej"))
	if err != nil {
		t.Fatalf("Unexpected error: %s", err.Error())
	}
	if n != 1 {
		t.Fatalf("Unexpected buffer count returned. Expected: %d - Found: %d.", 1, n)
	}
	n, err = b.Write([]byte("hej"))
	if err != nil {
		t.Fatalf("Unexpected error: %s", err.Error())
	}
	if n != 2 {
		t.Fatalf("Unexpected buffer count returned. Expected: %d - Found: %d.", 2, n)
	}
	n, err = b.Write([]byte("monika"))
	if err != nil {
		t.Fatalf("Unexpected error: %s", err.Error())
	}
	if n != 3 {
		t.Fatalf("Unexpected buffer count returned. Expected: %d - Found: %d.", 3, n)
	}

	time.Sleep(2 * time.Second)

	// Read from chan
	select {
	case messagesBuffer := <-mock.out:
		messages := strings.Split(messagesBuffer, "\n")
		if len(messages) != 3 {
			t.Logf("Messages: %+v.", messages)
			t.Fatalf("Unexpected number of messages. Expected: %d - Found: %d.", 3, len(messages))
		}
		if messages[0] != "hej" {
			t.Fatalf("Unexpected message[0]. Expected: %s - Found: %s.", "hej", messages[0])
		}
		if messages[1] != "hej" {
			t.Fatalf("Unexpected message[0]. Expected: %s - Found: %s.", "hej", messages[1])
		}
		if messages[2] != "monika" {
			t.Fatalf("Unexpected message[0]. Expected: %s - Found: %s.", "monika", messages[2])
		}
	default:
		t.Fatalf("Failed read from mock stream.")
	}
	b.routineMutex.Lock()
	b.routineRunning = false
	b.routineMutex.Unlock()
	time.Sleep(1 * time.Second)
}

func TestAutonomousTransmissionSafeStop(t *testing.T) {
	mock := newMockStream(1)
	b := NewBufferedStream(mock)

	ticker := time.NewTicker(1 * time.Second)

	go func(ticker *time.Ticker) {
		if err := b.autonomousTranmissionRoutine(ticker); err != nil {
			t.Fatalf("Unexpected error on exit: %s.", err.Error())
		}
	}(ticker)

	time.Sleep(2 * time.Second)

	// Force ticker stop from the outside.
	ticker.Stop()
}

func TestAutonomousTransmissionErrors(t *testing.T) {
	failer := newFailerMockStream("fail")
	b := NewBufferedStream(failer)

	ticker := time.NewTicker(1 * time.Second)

	n, err := b.Write([]byte("write something in order to trigger a buffer flush"))
	if err != nil {
		t.Fatalf("Unexpected error: %s.", err.Error())
	}
	if n != 1 {
		t.Fatalf("Unexpected count: Expected: %d - Found: %d.", 1, n)

	}

	err = b.autonomousTranmissionRoutine(ticker)
	if err == nil {
		t.Fatalf("Expected error, found nil.")
	}
	t.Logf("Expected error: %s.", err.Error())
}

func TestStartAutonomousTransmissionErrors(t *testing.T) {
	mock := newMockStream(3)
	b := NewBufferedStream(mock)

	// Falsify routine running flag.
	b.routineRunning = true
	if err := b.StartAutonomousTransmission(); err == nil {
		t.Fatalf("Expected error. Found nil.")
	}

	b.routineRunning = false
	b.scheduleInteval = 0
	if err := b.StartAutonomousTransmission(); err == nil {
		t.Fatalf("Expected error. Found nil.")
	}
}

func TestSetFatalFn(t *testing.T) {
	b := NewBufferedStream(nil)
	mtx := sync.Mutex{}
	invoked := false

	fn := func(err error) {
		mtx.Lock()
		invoked = true
		mtx.Unlock()
	}
	b.SetFatalFn(fn)

	b.fatal(fmt.Errorf("that"))
	time.Sleep(500 * time.Millisecond)

	mtx.Lock()
	hasBeenInvoked := invoked
	mtx.Unlock()

	if !hasBeenInvoked {
		t.Fatalf("The fatal fn should habe been invoked!")
	}
}

func TestSetFlatBufferSeparator(t *testing.T) {
	b := NewBufferedStream(nil)
	if b.separator != '\n' {
		t.Fatalf("Unexpected default separator. Expected: %c - Found: %c", '\n', b.separator)
	}
	b.SetFlatBufferSeparator('c')
	if b.separator != 'c' {
		t.Fatalf("Unexpected default separator. Expected: %c - Found: %c", 'c', b.separator)
	}
}

func TestWriteWithCappedLimitBuffer(t *testing.T) {
	s := newMockStream(1)
	b := NewBufferedStream(s)
	b.SetBufferLimit(10)

	for i := 0; i < 11; i++ {
		n, err := b.Write([]byte(fmt.Sprintf("message_%d", i)))
		if err != nil {
			t.Fatalf("Unexpected error: %s.", err.Error())
		}

		if i < 10 {
			if n != i+1 {
				t.Fatalf("Unexpected number returned. Expected: %d - Found: %d.", i+1, n)
			}
		}
		if i == 10 {
			if n != 1 {
				t.Fatalf("Unexpected number returned. Expected: %d - Found: %d.", 1, n)
			}
		}
	}

	received := <-s.out
	t.Logf("Received bytes: `%s`.", string(received))

	splitted := strings.Split(received, "\n")
	for i := 0; i < 10; i++ {
		if splitted[i] != fmt.Sprintf("message_%d", i) {
			t.Fatalf("Unexpected string for position %d. Expected: %s - Found: %s.", i, fmt.Sprintf("message_%d", i), splitted[i])
		}
	}

	if bytes.Compare([]byte("message_10"), b.buffer[0]) != 0 {
		t.Fatalf("Unexpected value in buffer position 0. Expected: %s - Found: %s.", "message_10", string(b.buffer[0]))
	}
}

func TestBufferAutoResizingNoLimit(t *testing.T) {
	b := NewBufferedStream(nil)
	if set, _ := b.SetStartingSize(5, false); !set {
		t.Fatalf("Starting size set failed.")
	}

	if len(b.buffer) != 5 {
		t.Fatalf("Unexpected buffer size. Expected: %d - Found: %d.", 5, len(b.buffer))
	}

	for i := 0; i < 10; i++ {
		n, err := b.Write([]byte(fmt.Sprintf("message_%d", i)))
		if err != nil {
			t.Fatalf("Unexpected error: %s.", err.Error())
		}
		if n != i+1 {
			t.Fatalf("Unexpected number returned. Expected: %d - Found: %d.", i+1, n)
		}
	}

	if len(b.buffer) != 10 {
		t.Fatalf("Unexpected buffer size. Expected: %d - Found: %d.", 10, len(b.buffer))
	}

	for i := 0; i < 10; i++ {
		if string(b.buffer[i]) != fmt.Sprintf("message_%d", i) {
			t.Fatalf("Unexpected string for position %d. Expected: %s - Found: %s.", i, fmt.Sprintf("message_%d", i), string(b.buffer[i]))
		}
	}
}

func TestFireTransmission(t *testing.T) {
	s := newMockStream(1)
	b := NewBufferedStream(s)

	data := [][]byte{
		[]byte("hey"),
		[]byte("oh"),
		[]byte("let's"),
		[]byte("go"),
	}
	if err := b.fireTransmission(data, 4); err != nil {
		t.Fatalf("Unexpected error: %s", err.Error())
	}

	received := <-s.out
	receivedSliced := strings.Split(received, string('\n'))

	if len(receivedSliced) != 4 {
		t.Fatalf("Unexpected number of messages. Expected: %d - Found: %d.", 4, len(receivedSliced))
	}
	if receivedSliced[0] != "hey" {
		t.Fatalf("Unexpected message in pos [%d]. Expected: %s - Found: %s.", 0, "hey", receivedSliced[0])
	}
	if receivedSliced[1] != "oh" {
		t.Fatalf("Unexpected message in pos [%d]. Expected: %s - Found: %s.", 1, "oh", receivedSliced[1])
	}
	if receivedSliced[2] != "let's" {
		t.Fatalf("Unexpected message in pos [%d]. Expected: %s - Found: %s.", 2, "let's", receivedSliced[2])
	}
	if receivedSliced[3] != "go" {
		t.Fatalf("Unexpected message in pos [%d]. Expected: %s - Found: %s.", 3, "go", receivedSliced[3])
	}
}

func TestFireTransmissionFailure(t *testing.T) {
	// Failure for no is stream set.
	b := NewBufferedStream(nil)
	// Fire the transmission but no Stream has been
	// provided and it will fail immediately.
	err := b.fireTransmission(nil, 0)
	if err == nil {
		t.Fatalf("Unexpected nil error!")
	}
	t.Logf("Expected error: %s.", err.Error())

	// Failure due to Stream.Write errors.
	s := newMockStream(2)
	b = NewBufferedStream(s)

	// Fill up the mock stream.
	s.out <- "hey"
	s.out <- "oh"
	// Prepare data to be fired.
	data := [][]byte{
		[]byte("let's"),
		[]byte("go"),
	}
	// Fire the data, but the chan is already full
	// so this should fail.
	err = b.fireTransmission(data, 2)
	if err == nil {
		t.Fatalf("Unexpected nil error!")
	}
	t.Logf("Expected error: %s.", err.Error())
}

func TestWriteFailureFatalInvocation(t *testing.T) {
	mtx := sync.Mutex{}
	invoked := false

	// Failure due to Stream.Write errors.
	s := newMockStream(1)
	b := NewBufferedStream(s)
	b.SetFatalFn(func(err error) {
		mtx.Lock()
		invoked = true
		mtx.Unlock()
	})
	b.SetBufferLimit(2)
	n, err := b.Write([]byte("hey"))
	if err != nil {
		t.Fatalf("Unexpected error: %s.", err.Error())
	}
	if n != 1 {
		t.Fatalf("Unexpected number returned. Expected: %d - Found: %d", 1, n)
	}
	n, err = b.Write([]byte("oh"))
	if err != nil {
		t.Fatalf("Unexpected error: %s.", err.Error())
	}
	if n != 2 {
		t.Fatalf("Unexpected number returned. Expected: %d - Found: %d", 2, n)
	}
	n, err = b.Write([]byte("let's"))
	if err != nil {
		t.Fatalf("Unexpected error: %s.", err.Error())
	}
	if n != 1 {
		t.Fatalf("Unexpected number returned. Expected: %d - Found: %d", 1, n)
	}
	n, err = b.Write([]byte("go"))
	if err != nil {
		t.Fatalf("Unexpected error: %s.", err.Error())
	}
	if n != 2 {
		t.Fatalf("Unexpected number returned. Expected: %d - Found: %d", 2, n)
	}

	mtx.Lock()
	hasBeenInvoked := invoked
	mtx.Unlock()

	if hasBeenInvoked {
		t.Fatalf("The fatal fn should not have been invoked!")
	}

	n, err = b.Write([]byte("AND NOW THE FATAL"))
	if err != nil {
		t.Fatalf("Unexpected error: %s.", err.Error())
	}
	if n != 1 {
		t.Fatalf("Unexpected number returned. Expected: %d - Found: %d", 1, n)
	}

	time.Sleep(500 * time.Millisecond)

	mtx.Lock()
	hasBeenInvoked = invoked
	mtx.Unlock()

	if !hasBeenInvoked {
		t.Fatalf("The fatal fn should have been invoked!")
	}
}

func TestFlush(t *testing.T) {
	b := NewBufferedStream(nil)

	// Mess up the internal values.
	b.bufferCount = 37
	b.buffer = make([][]byte, 25)

	oldBuff, oldCount := b.flush()
	if len(oldBuff) != 25 {
		t.Fatalf("The length of the old buffer should be %d. Found: %d.", 25, len(oldBuff))
	}
	if oldCount != 37 {
		t.Fatalf("The old count should be %d. Found: %d.", 37, oldCount)
	}

	if len(b.buffer) != DefaultPreallocatedBufferSize {
		t.Fatalf("The length of the new buffer should be %d. Found: %d.", DefaultPreallocatedBufferSize, len(b.buffer))
	}
	if b.bufferCount != 0 {
		t.Fatalf("The new count should be %d. Found: %d.", 0, b.bufferCount)
	}
}

func TestFlatten(t *testing.T) {
	raw := [][]byte{
		[]byte("a"),
		[]byte("b"),
		[]byte("c"),
		[]byte("d"),
	}
	expected := []byte("a|b|c|d")

	flattened := flatten(raw, byte('|'))
	if bytes.Compare(expected, flattened) != 0 {
		t.Fatalf("There where differences between the bytes. Expected: %s - Found: %s.", string(expected), string(flattened))
	}

	expected = []byte("a\nb\nc\nd")
	flattened = flatten(raw, byte('\n'))
	if bytes.Compare(expected, flattened) != 0 {
		t.Fatalf("There where differences between the bytes. Expected: %s - Found: %s.", string(expected), string(flattened))
	}

	rawArr := [5][]byte{
		[]byte("a"),
		[]byte("b"),
		[]byte("c"),
		[]byte("d"),
		nil,
	}
	expected = []byte("a:b:c:d")
	flattened = flatten(rawArr[:], byte(':'))
	if bytes.Compare(expected, flattened) != 0 {
		t.Fatalf("There where differences between the bytes. Expected: %s - Found: %s.", string(expected), string(flattened))
	}
}
