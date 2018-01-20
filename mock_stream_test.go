package gonyan

import (
	"fmt"
)

// mockStream is a simple Stream implementation used for testing purposes.
type mockStream struct {
	out chan string
}

func newMockStream(size int) *mockStream {
	return &mockStream{
		out: make(chan string, size),
	}
}

func (m *mockStream) Write(messageBytes []byte) (int, error) {
	if len(m.out) == cap(m.out) {
		return 0, fmt.Errorf("chan is full")
	}
	m.out <- string(messageBytes)
	return len(messageBytes), nil
}
