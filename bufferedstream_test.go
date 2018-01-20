package gonyan

import "testing"

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
