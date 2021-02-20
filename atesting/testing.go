package atesting

import (
	"container/list"
	"github.com/go-autowire/autowire"
	"log"
	"reflect"
	"unicode"
)

func Spy(v interface{}, dependency interface{}) {
	slice := []interface{}{dependency}
	Spies(v, slice)
}

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
							setValue(value, elem, i, dependency)
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

func setValue(value reflect.Value, elem reflect.Value, i int, dependency interface{}) bool {
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
