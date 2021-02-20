package service

import (
	"github.com/go-autowire/autowire"
	"log"
	"math/big"
)

func init() {
	autowire.InitProd(func() {
		autowire.Autowire(&BankAccountService{})
		autowire.Autowire(&PaypalService{})
	})
}

type PaymentService interface {
	Balance() *big.Float
}

type BankAccountService struct {
}

func (BankAccountService) Balance() *big.Float {
	log.Println("BankAccountService...")
	balance, _ := new(big.Float).SetString("600.10")
	return balance
}

type PaypalService struct {
}

func (PaypalService) Balance() *big.Float {
	log.Println("PaypalService...")
	balance, _ := new(big.Float).SetString("100.10")
	return balance
}
