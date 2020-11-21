package example_application

import (
	"autowire"
	"autowire/atesting"
	"autowire/example_application/service"
	"log"
	"testing"
)

type FilipPaymentServiceTest struct {
}

func (FilipPaymentServiceTest) Status() {
	log.Println("FilipPaymentServiceTest...")
}

type TestClient struct {
	T string
}

func (TestClient) Connect() {
	log.Println("TestConnected")
}

func TestAutowire(t *testing.T) {
	var service = autowire.Autowired(service.UserService{}).(*service.UserService)
	atesting.Spies(service, []interface{}{&FilipPaymentServiceTest{}, &TestClient{T: "filip"}})
	log.Println(service.Do())
}
