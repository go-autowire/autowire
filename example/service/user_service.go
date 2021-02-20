package service

import (
	"fmt"
	"github.com/go-autowire/autowire"
	"github.com/go-autowire/autowire/example/repository"
	"log"
	"math/big"
)

func init() {
	autowire.Autowire(&repository.InMemoryUserRoleRepository{})
	autowire.Autowire(&UserService{})
}

type UserService struct {
	PaymentSvc         PaymentService                `autowire:"service/BankAccountService"`
	auditClient        AuditEventSender              `autowire:"service/AuditClient"`
	userRoleRepository repository.UserRoleRepository `autowire:"repository/InMemoryUserRoleRepository"`
}

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

func (u *UserService) SetAuditClient(auditClient AuditEventSender) {
	u.auditClient = auditClient
}

func (u *UserService) SetUserRoleRepository(userRoleRepository repository.UserRoleRepository) {
	u.userRoleRepository = userRoleRepository
}
