package internal

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

const unexportedField = "unexported"

type Foo struct {
	unexported string
}

func TestGetUnexportedField(t *testing.T) {
	foo := &Foo{unexported: "x"}
	value := GetUnexportedField(reflect.ValueOf(foo).Elem().FieldByName(unexportedField))
	assert.Equal(t, value, "x")
}

func TestSetUnexportedField(t *testing.T) {
	foo := &Foo{unexported: "old"}
	SetUnexportedField(reflect.ValueOf(foo).Elem().FieldByName(unexportedField), "changed")
	assert.Equal(t, foo.unexported, "changed")
}

func TestSetFieldValue(t *testing.T) {
	foo := &Foo{unexported: "old"}
	SetFieldValue(reflect.ValueOf(foo).Elem(), 0, "changed")
	assert.Equal(t, foo.unexported, "changed")
}
