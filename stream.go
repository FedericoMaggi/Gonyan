package gonyan

// Stream interface holds the protocol to allow custom streams definition.
type Stream interface {
	Write([]byte) (int, error)
}
