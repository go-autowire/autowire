package autowire

import (
	"testing"
)

func Test_Profile(t *testing.T) {
	profile := getProfile()
	if profile != _Testing {
		t.Errorf("Expected profile testing found active")
	}
}
