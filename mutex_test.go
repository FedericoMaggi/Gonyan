package gonyan

import "testing"

// TestMutexDisable verifies that the disabled flag is properly set when
// calling Disable function.
func TestMutexDisable(t *testing.T) {
	m := mutex{}
	if m.disabled != false {
		t.Fatalf("Disabled flag should be false!")
	}

	m.Disable()
	if m.disabled == false {
		t.Fatalf("Disabled flag should be true!")
	}
}

// TestMutexDisable verifies that the disabled flag is properly set when
// calling Enable function.
func TestMutexEnable(t *testing.T) {
	m := mutex{}
	if m.disabled != false {
		t.Fatalf("Disabled flag should be false!")
	}

	m.Disable()
	m.Enable()
	if m.disabled != false {
		t.Fatalf("Disabled flag should be false!")
	}
}
