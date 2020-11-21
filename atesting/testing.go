package atesting

import (
	"autowire"
	"container/list"
	"log"
	"reflect"
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
		log.Printf("+++ %+v", value)
		var elem reflect.Value
		switch value.Kind() {
		case reflect.Ptr:
			elem = value.Elem()
		case reflect.Struct:
			elem = value
		}
		for i := 0; i < elem.NumField(); i++ {
			field := elem.Type().Field(i)
			log.Printf("[%s] %s:%+v", field.Type.String(), field.Name, elem.Field(i))
			tag, ok := field.Tag.Lookup(autowire.Tag)
			if ok {
				if tag != "" {
					for _, dependency := range dependencies {
						dependValue := reflect.ValueOf(dependency)
						if dependValue.Type().Implements(field.Type) {
							t := reflect.New(dependValue.Type())
							log.Println("found dependency by tag " + tag + " will be used " + t.Type().String())
							elem.Field(i).Set(reflect.ValueOf(dependency))
						}
					}
				}
				if !elem.Field(i).IsNil() {
					queue.PushBack(elem.Field(i).Elem().Interface())
				}
			}
		}
	}
}
