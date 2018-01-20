package gonyan

import (
	"bytes"
	"strings"
	"testing"
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
	b := NewBufferedStream(nil)
	if b.initialSize != DefaultPreallocatedBufferSize {
		t.Fatalf("Unexpected default initialSize value. Expected: %d - found: %d", DefaultPreallocatedBufferSize, b.initialSize)
	}

	set, err := b.SetStartingSize(-1, true)
	if set {
		t.Fatalf("Unexpected set flag should be false, found: %t.", set)
	}
	if err != nil {
		t.Fatalf("Unexpected error: %s", err.Error())
	}

	set, err = b.SetStartingSize(DefaultPreallocatedBufferSize, true)
	if set {
		t.Fatalf("Unexpected set flag, should be false, found: %t.", set)
	}
	if err != nil {
		t.Fatalf("Unexpected error: %s", err.Error())
	}

	set, err = b.SetStartingSize(10, false)
	if !set {
		t.Fatalf("Unexpected set flag, should be true, found: %t.", set)
	}
	if err != nil {
		t.Fatalf("Unexpected error: %s", err.Error())
	}
	if b.initialSize != 10 {
		t.Fatalf("Unexpected initial buffer size. Expected: %d - Found: %d.", 10, b.initialSize)
	}

	// Flag for transmission is true but the stream is not valid
	// still we do not expect an error since the buffer is empty.
	set, err = b.SetStartingSize(11, true)
	if !set {
		t.Fatalf("Unexpected set flag, should be true, found: %t.", set)
	}
	if err != nil {
		t.Fatalf("Unexpected error: %s", err.Error())
	}
	if b.initialSize != 11 {
		t.Fatalf("Unexpected default initial buffer size. Expected: %d - Found: %d.", 11, b.initialSize)
	}

	// Flag for transmission is true but the stream is not valid and
	// we expect an error since we have sent a message to the stream.
	b.Write([]byte(`logging something nasty`))
	set, err = b.SetStartingSize(5, true)
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
	if received != "logging something nasty, again\n" {
		t.Fatalf("Unexpected received message. Expected: `%s` - Found: `%s`.", "logging something nasty, again\n", received)
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
	separator := byte('|')
	expected := []byte("a|b|c|d")

	flattened := flatten(raw, separator)
	if bytes.Compare(expected, flattened) != 0 {
		t.Fatalf("There where differences between the bytes. Expected: %s - Found: %s.", string(expected), string(flattened))
	}
}
