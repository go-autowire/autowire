package example

import (
	"github.com/go-autowire/autowire"
	"github.com/go-autowire/autowire/atesting"
	"github.com/go-autowire/autowire/example/app"
	"log"
	"math/big"
	"testing"
)

type TestPaymentServiceTest struct {
}

func (TestPaymentServiceTest) Balance() *big.Float {
	log.Println("Mocked object...TestPaymentServiceTest...")
	balance, _ := new(big.Float).SetString("300.10")
	return balance
}

type TestAuditClient struct {
}

func (TestAuditClient) Send(event string) {
	log.Printf("Test event delivered")
}

func TestAutowire(t *testing.T) {
	app := autowire.Autowired(app.Application{}).(*app.Application)
	app.UserSvc().PaymentSvc = &TestPaymentServiceTest{}
	atesting.Spies(app, []interface{}{&TestPaymentServiceTest{}, &TestAuditClient{}})
	app.Start()
}
