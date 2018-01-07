package gonyan

import (
	"testing"
)

func TestGetLevelLabel(t *testing.T) {
	if label := GetLevelLabel(Debug); label != "Debug" {
		t.Fatalf("Invalid level label for Debug level found: %s", label)
	}
	if label := GetLevelLabel(Verbose); label != "Verbose" {
		t.Fatalf("Invalid level label for Verbose level found: %s", label)
	}
	if label := GetLevelLabel(Info); label != "Info" {
		t.Fatalf("Invalid level label for Info level found: %s", label)
	}
	if label := GetLevelLabel(Warning); label != "Warning" {
		t.Fatalf("Invalid level label for Warning level found: %s", label)
	}
	if label := GetLevelLabel(Error); label != "Error" {
		t.Fatalf("Invalid level label for Error level found: %s", label)
	}
	if label := GetLevelLabel(Fatal); label != "Fatal" {
		t.Fatalf("Invalid level label for Fatal level found: %s", label)
	}
	if label := GetLevelLabel(Panic); label != "Panic" {
		t.Fatalf("Invalid level label for Panic level found: %s", label)
	}
}
