package service

import (
	"fmt"
	"github.com/go-autowire/autowire/example/repository"
	"log"
	"math/big"
)

// A UserService represents a named struct
type UserService struct {
	PaymentSvc         PaymentService                `autowire:"service/BankAccountService"`
	auditClient        EventSender                   `autowire:"service/AuditService"`
	userRoleRepository repository.UserRoleRepository `autowire:"repository/InMemoryUserRoleRepository"`
}

// Balance is function returning current balance of the user
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

// SetAuditClient is a Setter for auditClient field
func (u *UserService) SetAuditClient(auditClient EventSender) {
	u.auditClient = auditClient
}

// SetUserRoleRepository is a Setter for userRoleRepository field
func (u *UserService) SetUserRoleRepository(userRoleRepository repository.UserRoleRepository) {
	u.userRoleRepository = userRoleRepository
}
