package atesting

import (
	"container/list"
	"github.com/go-autowire/autowire"
	"log"
	"reflect"
	"unicode"
)

// Spy Function is replacing object field with the one provided in the function as a second argument.
// Spy Function detects automatically which field could be replaced with the provided one.
// Important note: In order to traverse fields of the unexported fields we need to implement Getters.
// As shown inside example package, we are replacing AuditClient with our mock implementation and in order to reach this
// field we need Getter.
// Example:
//   atesting.Spy(application, &TestAuditClient{})
// Or this is equivalent of
//   application.UserSvc().SetAuditClient(&TestAuditClient{})
// Getter UserSvc() is used to access userSvc field, which is unexported. For more information take a look at example package.
//Parameters of Spy function:
//   - `v`          : pointer to structure inside which spy object will be applied
//   - `dependency` : pointer to structure which will be injected
func Spy(v interface{}, dependency interface{}) {
	slice := []interface{}{dependency}
	Spies(v, slice)
}

// Spies Function is replacing object fields with the list of provided dependencies in the function as a second argument.
// Spy Function detects automatically which field could be replaced with the provided one in the list of dependencies.
// Important note: In order to traverse fields of the unexported fields we need to implement Getters.
// As shown inside example package, we are replacing AuditClient with our mock implementation and in order to reach this
// field we need Getter.
// Example:
//   	atesting.Spies(application, []interface{}{&TestPaymentServiceTest{}, &TestAuditClient{}})
// Or this is equivalent of
//   application.UserSvc().SetAuditClient(&TestAuditClient{})
//   application.UserSvc().PaymentSvc = &TestPaymentServiceTest{}
// Getter UserSvc() is used to access userSvc field, which is unexported. In case of PaymentSvc it is not required as field PaymentSvc is exported.
// Parameters of Spies function:
//   - `v`            : structure inside which spy objects will be applied
//   - `dependencies` : list of dependencies which will be injected
// For more information take a look at example package.
func Spies(v interface{}, dependencies []interface{}) {
	queue := list.New()
	queue.PushBack(v)
	for queue.Len() > 0 {
		elemQueue := queue.Front()
		value := reflect.ValueOf(elemQueue.Value)
		queue.Remove(elemQueue)
		var elem reflect.Value
		switch value.Kind() {
		case reflect.Ptr:
			elem = value.Elem()
		case reflect.Struct:
			elem = value
		}
		for i := 0; i < elem.NumField(); i++ {
			field := elem.Type().Field(i)
			tag, ok := field.Tag.Lookup(autowire.Tag)
			if ok {
				if tag != "" {
					for _, dependency := range dependencies {
						dependValue := reflect.ValueOf(dependency)
						if dependValue.Type().Implements(field.Type) {
							t := reflect.New(dependValue.Type())
							log.Println("Injecting Spy on dependency by tag " + tag + " will be used " + t.Type().String())
							autowire.Autowire(dependency)
							setFieldValue(value, elem, i, dependency)
						}
					}
				}
				if !elem.Field(i).IsNil() {
					if elem.Field(i).Elem().CanInterface() {
						queue.PushBack(autowire.Autowired(elem.Field(i).Elem().Interface()))
					} else {
						runeName := []rune(elem.Type().Field(i).Name)
						runeName[0] = unicode.ToUpper(runeName[0])
						methodName := string(runeName)
						method := value.MethodByName(methodName)
						if method.IsValid() {
							result := method.Call([]reflect.Value{})[0]
							queue.PushBack(autowire.Autowired(result.Interface()))
						}
					}
				}
			}
		}
	}
}

func setFieldValue(value reflect.Value, elem reflect.Value, i int, dependency interface{}) bool {
	runeName := []rune(elem.Type().Field(i).Name)
	exported := unicode.IsUpper(runeName[0])
	if exported {
		elem.Field(i).Set(reflect.ValueOf(dependency))
	} else {
		runeName[0] = unicode.ToUpper(runeName[0])
		methodName := "Set" + string(runeName)
		method := value.MethodByName(methodName)
		if method.IsValid() {
			method.Call([]reflect.Value{reflect.ValueOf(dependency)})
			return true
		}
	}
	return false
}
