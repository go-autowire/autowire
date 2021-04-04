package internal

import (
	"os"
	"regexp"
	"strings"
)

type profile uint64

const (
	Production profile = iota
	Testing
)

func GetProfile() profile {
	args := os.Args
	programName := args[0][strings.LastIndex(args[0], "/"):]
	if result, _ := regexp.MatchString("/.*[Tt]est", programName); result {
		return Testing
	}
	return Production
}
