package repository

import "github.com/go-autowire/autowire"

func init() {
	autowire.Autowire(&InMemoryUserRoleRepository{})
}

type UserRole string

const (
	OwnerRole UserRole = "owner"
)

func (u UserRole) String() string {
	return string(u)
}

type UserRoleRepository interface {
	GetAllRoles(userId string) ([]UserRole, error)
}

type InMemoryUserRoleRepository struct {
}

func (i InMemoryUserRoleRepository) GetAllRoles(_ string) ([]UserRole, error) {
	return []UserRole{OwnerRole}, nil
}
