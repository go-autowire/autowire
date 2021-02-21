package service

import (
	"log"
	"math/big"
)

// A PaymentService represents interface with one function, which returns Account balance
type PaymentService interface {
	// The implementation of Balance returns Account balance
	Balance() *big.Float
}

// A BankAccountService represents named struct
type BankAccountService struct {
}

// A Balance returns Account balance
func (BankAccountService) Balance() *big.Float {
	log.Println("BankAccountService...")
	balance, _ := new(big.Float).SetString("600.10")
	return balance
}

// A PaypalService represents named struct
type PaypalService struct {
}

// A Balance returns Account balance
func (PaypalService) Balance() *big.Float {
	log.Println("PaypalService...")
	balance, _ := new(big.Float).SetString("100.10")
	return balance
}
