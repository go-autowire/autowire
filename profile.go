package autowire

import (
	"os"
	"regexp"
	"strings"
)

type Profile uint64

const (
	_Default Profile = iota
	_Testing
)

func getProfile() Profile {
	args := os.Args
	programName := args[0][strings.LastIndex(args[0], "/"):]
	if result, _ := regexp.MatchString("/.*[Tt]est", programName); result {
		return _Testing
	}
	return _Default
}
