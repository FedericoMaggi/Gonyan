package gonyan

import "testing"
import "os"

func TestNewLogger(t *testing.T) {
	l := NewLogger("X", nil, true)
	l.RegisterStream(Debug, os.Stdout)
}
