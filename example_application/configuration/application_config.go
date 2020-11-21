package configuration

import (
	"autowire"
	"log"
)

func init() {
	log.Println("Initializing ApplicationConfig")
	autowire.Autowire(New("con"))
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
