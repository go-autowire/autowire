package configuration

import (
	"github.com/go-autowire/autowire"
)

func init() {
	autowire.Autowire(New("default"))
}

type ApplicationConfig struct {
	kind string
}

func New(kind string) *ApplicationConfig {
	return &ApplicationConfig{kind: kind}
}

func (a ApplicationConfig) Kind() string {
	return a.kind
}
