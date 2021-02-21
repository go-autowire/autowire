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

func (TestAuditClient) Send(_ string) {
	log.Printf("Test event delivered")
}

func TestAutowire(t *testing.T) {
	application := autowire.Autowired(app.Application{}).(*app.Application)
	atesting.Spies(application, []interface{}{&TestPaymentServiceTest{}, &TestAuditClient{}})
	application.Start()
}
