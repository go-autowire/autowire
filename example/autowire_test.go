package example

import (
	"github.com/go-autowire/autowire/example/app"
	"github.com/go-autowire/autowire/pkg"
	"github.com/go-autowire/autowire/pkg/atesting"
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
	defer pkg.Close()
	application := pkg.Autowired(app.Application{}).(*app.Application)
	atesting.Spy(application, &TestPaymentServiceTest{}, &TestAuditClient{})
	application.Start()
}
