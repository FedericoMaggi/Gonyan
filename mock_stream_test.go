package gonyan

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
	m.out <- string(messageBytes)
	return len(messageBytes), nil
}
