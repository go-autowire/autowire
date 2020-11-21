package service

import (
	"autowire"
	"log"
)

func init() {
	log.Println("Initializing PaymentService")
	autowire.RunOmitTest(func() {
		autowire.Autowire(&BankAccountService{})
		autowire.Autowire(&PaypalService{})
	})
}

type PaymentService interface {
	Status()
}

type BankAccountService struct {
}

func (BankAccountService) Status() {
	log.Println("BankAccountService...")
}

type PaypalService struct {
}

func (PaypalService) Status() {
	log.Println("PaypalService...")
}
