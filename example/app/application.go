// Package app holds primary application code
package app

import (
	"github.com/go-autowire/autowire/pkg"
	"log"

	"github.com/go-autowire/autowire/example/configuration"
	"github.com/go-autowire/autowire/example/service"
)

func init() {
	pkg.Autowire(&Application{})
}

// A Application represents named struct
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
