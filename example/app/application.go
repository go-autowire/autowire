package app

import (
	"fmt"
	"github.com/go-autowire/autowire"
	"github.com/go-autowire/autowire/example/configuration"
	"github.com/go-autowire/autowire/example/service"
	"log"
)

func init() {
	autowire.Autowire(&Application{})
}

type Application struct {
	config  *configuration.ApplicationConfig `autowire:""`
	userSvc *service.UserService             `autowire:""`
}

func (a Application) Start() {
	log.Println("Config Kind : " + a.config.Kind())
	userId := "serviceaccount@demo.com"
	balance, err := a.userSvc.Balance(userId)
	if err != nil {
		log.Fatalln(err.Error())
	}
	fmt.Println("Current balance is " + balance.String())
}

func (a *Application) SetConfig(config *configuration.ApplicationConfig) {
	a.config = config
}

func (a *Application) SetUserSvc(userSvc *service.UserService) {
	a.userSvc = userSvc
}

func (a *Application) UserSvc() *service.UserService {
	return a.userSvc
}
