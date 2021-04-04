package internal

import (
	"testing"
)

func Test_Profile(t *testing.T) {
	profile := GetProfile()
	if profile != Testing {
		t.Errorf("Expected Profile testing found active")
	}
}
