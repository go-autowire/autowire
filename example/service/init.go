// Package service holds all services.
package service

import (
	"github.com/go-autowire/autowire/pkg"
)

// Important Note: First we autowire independent structures
// and the most complex one are at the end of the init function,
// as independent one are injected into others
func init() { //nolint:gochecknoinits
	pkg.Autowire(&AuditService{})
	pkg.InitProd(func() {
		pkg.Autowire(&BankAccountService{})
		pkg.Autowire(&PaypalService{})
	})
	pkg.Autowire(&UserService{})
}
