package service

import (
	"autowire"
	"autowire/example_application/configuration"
	"log"
)

func init() {
	log.Println("Initializing UserService")
	s := &UserService{}
	autowire.Autowire(s)
}

type UserService struct {
	Config     *configuration.ApplicationConfig `autowire:""`
	PaymentSvc PaymentService                   `autowire:"service/BankAccountService"`
	Client     Client                           `autowire:"service/AppClient"`
}

func (u UserService) Do() string {
	u.Client.Connect()
	u.PaymentSvc.Status()
	return u.Config.Kind()
}
