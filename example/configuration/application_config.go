package configuration

import (
	"github.com/go-autowire/autowire"
)

func init() {
	autowire.Autowire(New("default"))
}

// A ApplicationConfig represents name struct, which hold application configuration
type ApplicationConfig struct {
	apiKey string
}

// New returns new ApplicationConfig
func New(apiKey string) *ApplicationConfig {
	return &ApplicationConfig{apiKey: apiKey}
}

// ApiKey is a Setter, which returns apiKey value
func (a ApplicationConfig) ApiKey() string {
	return a.apiKey
}