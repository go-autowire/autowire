// repository package
package repository

import "github.com/go-autowire/autowire"

func init() {
	autowire.Autowire(&InMemoryUserRoleRepository{})
}

// UserRole type
type UserRole string

const (
	// OwnerRole literal
	OwnerRole UserRole = "owner"
)

// A String returns UserRole as a string
func (u UserRole) String() string {
	return string(u)
}

// A UserRoleRepository represents interface containing roles related function: GetAllRoles
type UserRoleRepository interface {
	GetAllRoles(userId string) ([]UserRole, error)
}

// A InMemoryUserRoleRepository represents struct, which implements UserRoleRepository interface
type InMemoryUserRoleRepository struct {
}

// GetAllRoles returns all roles of the user
func (i InMemoryUserRoleRepository) GetAllRoles(_ string) ([]UserRole, error) {
	return []UserRole{OwnerRole}, nil
}
