// service package
package service

import . "github.com/go-autowire/autowire"

// Important Note: First we autowire independent structures
// and the most complex one are at the end of the init function,
// as independent one are injected into others
func init() {
	Autowire(&AuditService{})
	InitProd(func() {
		Autowire(&BankAccountService{})
		Autowire(&PaypalService{})
	})
	Autowire(&UserService{})
}
