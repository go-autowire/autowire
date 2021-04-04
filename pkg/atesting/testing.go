// Package atesting provides Spy function for easy way to mock dependencies
package atesting

import (
	"container/list"
	"log"
	"reflect"

	"github.com/go-autowire/autowire/pkg"
	"github.com/go-autowire/autowire/pkg/internal"
)

// Spy Function is replacing object field with the one provided in the function as a variadic arguments.
// Spy Function will detect fully automatically which field could be replaced
// with the provided one as variadic arguments.
// As shown inside example package, we are replacing AuditClient with our mock implementation.
// Example:
//   Spy(application, &TestAuditClient{})
// Or this is equivalent of doing it manually
//   application.UserSvc().SetAuditClient(&TestAuditClient{})
// When we don't use Spy function, we need to provide Getter method UserSvc() in order to
// access unexported userSvc field.
// For more information take a look at test file in example package.
// Parameters of Spy function:
//   - `v`          : pointer to structure inside which spy object will be injected
//   - `dependencies` : this is variadic argument, pointer to mocked structures which are gonna be injected
func Spy(v interface{}, dependencies ...interface{}) {
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
		}
		for i := 0; i < elem.NumField(); i++ {
			field := elem.Type().Field(i)
			tag, ok := field.Tag.Lookup(pkg.Tag)
			if ok {
				if tag != "" {
					for _, currentDependency := range dependencies {
						dependValue := reflect.ValueOf(currentDependency)
						if dependValue.Type().Implements(field.Type) {
							t := reflect.New(dependValue.Type())
							log.Println("Injecting Spy on currentDependency by tag " + tag + " will be used " + t.Type().String())
							pkg.Autowire(currentDependency)
							internal.SetFieldValue(elem, i, currentDependency)
						}
					}
				} else {
					for _, currentDependency := range dependencies {
						log.Printf("Checking compatibility between %s & %s", reflect.TypeOf(currentDependency), elem.Field(i).Type())
						if reflect.TypeOf(currentDependency) == elem.Field(i).Type() {
							internal.SetFieldValue(elem, i, currentDependency)
						}
					}
				}
				if !elem.Field(i).IsNil() {
					if elem.Field(i).Elem().CanInterface() {
						queue.PushBack(pkg.Autowired(elem.Field(i).Elem().Interface()))
					} else {
						queue.PushBack(pkg.Autowired(internal.GetUnexportedField(elem.Field(i))))
					}
				}
			}
		}
	}
}
