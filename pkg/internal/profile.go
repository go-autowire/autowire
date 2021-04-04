package internal

import (
	"os"
	"regexp"
	"strings"
)

// Profile represents enum type
type Profile uint64

const (
	// Production profile
	Production Profile = iota
	// Testing profile
	Testing
)

// GetProfile function returns current profile
func GetProfile() Profile {
	args := os.Args
	programName := args[0][strings.LastIndex(args[0], "/"):]
	if result, _ := regexp.MatchString("/.*[Tt]est", programName); result {
		return Testing
	}
	return Production
}
