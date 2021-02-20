package autowire

import (
	"os"
	"regexp"
	"strings"
)

type profile uint64

const (
	_Production profile = iota
	_Testing
)

func getProfile() profile {
	args := os.Args
	programName := args[0][strings.LastIndex(args[0], "/"):]
	if result, _ := regexp.MatchString("/.*[Tt]est", programName); result {
		return _Testing
	}
	return _Production
}
