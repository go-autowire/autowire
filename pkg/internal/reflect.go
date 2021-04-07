package internal

import (
	"reflect"
	"unicode"
	"unsafe"
)

// GetUnexportedField functions returns value of the unexported field
func GetUnexportedField(field reflect.Value) interface{} {
	//nolint:gosec
	return reflect.NewAt(field.Type(), unsafe.Pointer(field.UnsafeAddr())).Elem().Interface()
}

// SetUnexportedField functions sets value of the unexported field
func SetUnexportedField(field reflect.Value, value interface{}) {
	//nolint:gosec
	reflect.NewAt(field.Type(), unsafe.Pointer(field.UnsafeAddr())).
		Elem().
		Set(reflect.ValueOf(value))
}

// SetFieldValue functions injects dependency into a field.
// It supports exported and unexported fields.
func SetFieldValue(elem reflect.Value, i int, dependency interface{}) {
	runeName := []rune(elem.Type().Field(i).Name)
	exported := unicode.IsUpper(runeName[0])
	if exported {
		elem.Field(i).Set(reflect.ValueOf(dependency))
	} else {
		SetUnexportedField(elem.Field(i), dependency)
	}
}
