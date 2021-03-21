// app package
package app

import (
	"log"

	"github.com/go-autowire/autowire"
	"github.com/go-autowire/autowire/example/configuration"
	"github.com/go-autowire/autowire/example/service"
)

func init() {
	autowire.Autowire(&Application{})
}

// Application represents named struct
type Application struct {
	config  *configuration.ApplicationConfig `autowire:""`
	userSvc *service.UserService             `autowire:""`
}

// Start method is starting application
func (a Application) Start() {
	log.Println("Config ApiKey : " + a.config.ApiKey()[:3] + "****")
	userId := "serviceaccount@demo.com"
	balance, err := a.userSvc.Balance(userId)
	if err != nil {
		log.Fatalln(err.Error())
	}
	log.Println("Current balance is " + balance.String())
}

// SetConfig method is a Setter of private field config
func (a *Application) SetConfig(config *configuration.ApplicationConfig) {
	a.config = config
}

// SetUserSvc method is a Setter of private field userSvc
func (a *Application) SetUserSvc(userSvc *service.UserService) {
	a.userSvc = userSvc
}

// UserSvc method is a Getter of private field userSvc
func (a *Application) UserSvc() *service.UserService {
	return a.userSvc
}
