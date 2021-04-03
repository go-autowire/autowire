package autowire

import (
	"fmt"
	"github.com/go-autowire/autowire/internal"
	"github.com/go-autowire/autowire/internal/fake"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

const (
	packageName     = "github.com/go-autowire/autowire"
	myFooFieldName  = "myFoo"
	passerFieldName = "passer"
)

func TestInitProd(t *testing.T) {
	currentProfile = _Production
	callCount := 0
	InitProd(func() {
		callCount++
	})
	assert.Equal(t, callCount, 1)
	// set it back to test mode
	currentProfile = _Testing
}

func TestInitProdSkippedInTests(t *testing.T) {
	callCount := 0
	InitProd(func() {
		callCount++
	})
	assert.Equal(t, callCount, 0)
}

func Test_Close(t *testing.T) {
	tmp := fake.Foo{Name: "test"}
	value := reflect.ValueOf(&tmp)
	assert.Equal(t, value.Kind(), reflect.Ptr)
	structType := getStructPtrFullPath(value)
	assert.Equal(t, structType, packageName+"/internal/fake/Foo")
	dependencies[structType] = &tmp
	assert.Equal(t, tmp.CloseCalls, 0)
	errors := Close()
	assert.Equal(t, len(errors), 0)
	assert.Equal(t, tmp.CloseCalls, 1)
	delete(dependencies, structType)
}

type closeError struct {
	name       string
	closeCalls int
}

// Close method of Foo
func (f *closeError) Close() error {
	f.closeCalls++
	return fmt.Errorf("error in the close method")
}

func Test_CloseError(t *testing.T) {
	tmp := closeError{name: "test"}
	value := reflect.ValueOf(&tmp)
	assert.Equal(t, value.Kind(), reflect.Ptr)
	structType := getStructPtrFullPath(value)
	assert.Equal(t, structType, packageName+"/closeError")
	dependencies[structType] = &tmp
	assert.Equal(t, tmp.closeCalls, 0)
	errors := Close()
	assert.Equal(t, len(errors), 1)
	assert.Equal(t, tmp.closeCalls, 1)
	delete(dependencies, structType)
}

func Test_getStructPtrFullPath(t *testing.T) {
	tmp := fake.Foo{Name: "test"}
	value := reflect.ValueOf(&tmp)
	assert.Equal(t, value.Kind(), reflect.Ptr)
	structType := getStructPtrFullPath(value)
	assert.Equal(t, structType, packageName+"/internal/fake/Foo")
}

func Test_findDependency(t *testing.T) {
	tag := "fake/Foo"
	tmp := fake.Foo{Name: "test"}
	value := reflect.ValueOf(&tmp)
	structType := getStructPtrFullPath(value)
	dependencies[structType] = &tmp
	deps := findDependency(tag)
	assert.Equal(t, len(deps), 1)
	dependency := deps[0]
	dependencyType := reflect.TypeOf(dependency)
	assert.Equal(t, dependencyType.String(), reflect.TypeOf(&fake.Foo{}).String())
	delete(dependencies, structType)
}

func TestAutowireUnexportedStruct(t *testing.T) {
	Autowire(&fake.Foo{})
	tmpBar := &fake.Bar{}
	assert.Nil(t, getFieldByName(tmpBar, myFooFieldName))
	Autowire(tmpBar)
	assert.NotNil(t, getFieldByName(tmpBar, myFooFieldName))
	dependencies = make(map[string]interface{})
}

func TestAutowireExportedStruct(t *testing.T) {
	Autowire(&fake.Foo{})
	tmpBaz := &fake.Baz{}
	assert.Nil(t, tmpBaz.MyFoo)
	Autowire(tmpBaz)
	assert.NotNil(t, tmpBaz.MyFoo)
	dependencies = make(map[string]interface{})
}

func TestAutowireExportedInterface(t *testing.T) {
	Autowire(&fake.Foo{})
	tmpBaz := &fake.Qux{}
	assert.Nil(t, tmpBaz.Passer)
	Autowire(tmpBaz)
	assert.NotNil(t, tmpBaz.Passer)
	dependencies = make(map[string]interface{})
}

func TestAutowireUnexportedInterface(t *testing.T) {
	Autowire(&fake.Foo{})
	tmpBaz := &fake.Qus{}
	assert.Nil(t, tmpBaz.Passer())
	Autowire(tmpBaz)
	assert.NotNil(t, tmpBaz.Passer)
	dependencies = make(map[string]interface{})
}

func TestAutowireStructNotImplementingInterface(t *testing.T) {
	Autowire(&fake.Foo{})
	tmp := &fake.Bor{}
	assert.Nil(t, tmp.Passer)
	panicFunc := func() {
		Autowire(tmp)
	}
	assert.Panics(t, panicFunc)
	dependencies = make(map[string]interface{})
}

func TestAutowireUnknownStructOnInterfacePlaceholder(t *testing.T) {
	Autowire(&fake.Foo{})
	tmpDep := &fake.NotFoundTagDependency{}
	assert.Nil(t, getFieldByName(tmpDep, passerFieldName))
	Autowire(tmpDep)
	assert.Nil(t, getFieldByName(tmpDep, passerFieldName))
	dependencies = make(map[string]interface{})
}

func TestAutowireDependencyAlreadyAutowired(t *testing.T) {
	fistFoo := &fake.Foo{}
	Autowire(fistFoo)
	secondFoo := &fake.Foo{}
	Autowire(secondFoo)
	structType := getStructPtrFullPath(reflect.ValueOf(secondFoo))
	assert.Equal(t, dependencies[structType], fistFoo)
	dependencies = make(map[string]interface{})
}

func TestAutowireInvalid(t *testing.T) {
	panicFunc := func() {
		Autowire(nil)
	}
	assert.Panics(t, panicFunc)

	panicFunc = func() {
		Autowire(fake.Foo{})
	}
	assert.Panics(t, panicFunc)
}

func TestAutowired(t *testing.T) {
	Autowire(&fake.Foo{})
	tmpBar := &fake.Bar{}
	assert.Nil(t, getFieldByName(tmpBar, myFooFieldName))
	Autowire(tmpBar)
	assert.NotNil(t, getFieldByName(tmpBar, myFooFieldName))
	resultStruct := Autowired(fake.Bar{}).(*fake.Bar)
	assert.Equal(t, tmpBar, resultStruct)
	resultPtrStruct := Autowired(&fake.Bar{}).(*fake.Bar)
	assert.Equal(t, tmpBar, resultPtrStruct)
	dependencies = make(map[string]interface{})
}

func TestAutowiredNotFound(t *testing.T) {
	result := Autowired(fake.Foo{})
	assert.Nil(t, result)
}

func TestAutowiredInvalid(t *testing.T) {
	panicFunc := func() {
		Autowired(nil)
	}
	assert.Panics(t, panicFunc)
	dependencies = make(map[string]interface{})
}

func getFieldByName(v interface{}, fieldName string) interface{} {
	return internal.GetUnexportedField(reflect.ValueOf(v).Elem().FieldByName(fieldName))
}
