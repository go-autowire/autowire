package example_test

import (
	"log"
	"math/big"
	"testing"

	"github.com/go-autowire/autowire/example/app"
	"github.com/go-autowire/autowire/pkg"
	"github.com/go-autowire/autowire/pkg/atesting"
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

func TestExampleAutowire(t *testing.T) {
	defer pkg.Close()
	application := pkg.Autowired(app.Application{}).(*app.Application)
	atesting.Spy(application, &TestPaymentServiceTest{}, &TestAuditClient{})
	application.Start()
}
