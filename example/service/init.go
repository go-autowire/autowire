// Package service holds all services.
package service

import (
	"github.com/go-autowire/autowire/pkg"
)

func init() { //nolint:gochecknoinits
	pkg.Autowire(&UserService{})
	pkg.RunProd(func() {
		pkg.Autowire(&BankAccountService{})
		pkg.Autowire(&PaypalService{})
	})
	pkg.Autowire(&AuditService{})
}
