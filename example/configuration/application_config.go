// Package configuration holds application configs
package configuration

import (
	"github.com/go-autowire/autowire/pkg"
)

//nolint:gochecknoinits
func init() {
	pkg.Autowire(New("default"))
}

// A ApplicationConfig represents name struct, which hold application configuration
type ApplicationConfig struct {
	apiKey string
}

// New returns new ApplicationConfig
func New(apiKey string) *ApplicationConfig {
	return &ApplicationConfig{apiKey: apiKey}
}

// ApiKey is a Getter, which returns apiKey value
func (a ApplicationConfig) ApiKey() string { //nolint:revive,stylecheck
	return a.apiKey
}
