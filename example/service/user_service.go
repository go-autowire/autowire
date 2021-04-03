package service

import (
	"fmt"
	"log"
	"math/big"

	"github.com/go-autowire/autowire/example/repository"
)

// A UserService represents a named struct
type UserService struct {
	PaymentSvc         PaymentService                `autowire:"service/BankAccountService"`
	auditClient        EventSender                   `autowire:"service/AuditService"`
	userRoleRepository repository.UserRoleRepository `autowire:"repository/InMemoryUserRoleRepository"`
}

// Balance is a method returning current balance of the user.
func (u UserService) Balance(userId string) (*big.Float, error) {
	if u.validateUser(userId) {
		u.auditClient.Send("Balance:check")
		return u.PaymentSvc.Balance(), nil
	}
	return nil, fmt.Errorf("invalid user with id %s", userId)
}

func (u UserService) validateUser(userId string) bool {
	roles, err := u.userRoleRepository.GetAllRoles(userId)
	if err != nil {
		log.Println(err.Error())
		return false
	}
	for _, role := range roles {
		if role == repository.OwnerRole {
			return true
		}
	}
	return false
}
