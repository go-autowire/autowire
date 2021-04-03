package atesting_test

import (
	. "github.com/go-autowire/autowire"
	. "github.com/go-autowire/autowire/atesting"
	"github.com/stretchr/testify/assert"
	"testing"
)

type FooEr interface {
	Foo()
}

// A Foo represent named struct
type Foo struct {
	Name       string
	CloseCalls int
}

// Pass method
func (f Foo) Foo() {
}

type BarEr interface {
	Baz()
}

// Foo represent named struct
type Bar struct {
	Name       string
	CloseCalls int
}

// Pass method
func (b Bar) Baz() {
}

// A FooBarUnexported represent named struct
type FooBarUnexported struct {
	foo *Foo `autowire:""`
	bar *Bar `autowire:""`
}

// A FooBar represent named struct
type FooBar struct {
	Foo *Foo `autowire:""`
	Bar *Bar `autowire:""`
}

// Baz represent named struct
type Baz struct {
	MyFoo FooEr `autowire:"Foo"`
	MyBaz BarEr `autowire:"Bar"`
}

func TestSpyUnexportedStructPtr(t *testing.T) {
	fooName := "foo"
	Autowire(&Foo{Name: fooName})
	barName := "bar"
	Autowire(&Bar{Name: barName})
	tmpFooBar := &FooBarUnexported{}
	Autowire(tmpFooBar)
	assert.Equal(t, tmpFooBar.foo.Name, fooName)
	assert.Equal(t, tmpFooBar.bar.Name, barName)
	testFooName := "testFoo"
	testBarName := "testBar"
	Spy(tmpFooBar, &Foo{Name: testFooName}, &Bar{Name: testBarName})
	assert.Equal(t, tmpFooBar.foo.Name, testFooName)
	assert.Equal(t, tmpFooBar.bar.Name, testBarName)
	assert.Equal(t, 0, len(Close()))
}

func TestSpyExportedStructPtr(t *testing.T) {
	fooName := "foo"
	Autowire(&Foo{Name: fooName})
	barName := "bar"
	Autowire(&Bar{Name: barName})
	tmpFooBar := &FooBar{}
	Autowire(tmpFooBar)
	assert.Equal(t, tmpFooBar.Foo.Name, fooName)
	assert.Equal(t, tmpFooBar.Bar.Name, barName)
	testFooName := "testFoo"
	testBarName := "testBar"
	Spy(tmpFooBar, &Foo{Name: testFooName}, &Bar{Name: testBarName})
	assert.Equal(t, tmpFooBar.Foo.Name, testFooName)
	assert.Equal(t, tmpFooBar.Bar.Name, testBarName)
	assert.Equal(t, 0, len(Close()))
}

func TestSpyInterface(t *testing.T) {
	fooName := "foo"
	Autowire(&Foo{Name: fooName})
	barName := "bar"
	Autowire(&Bar{Name: barName})
	tmpBaz := &Baz{}
	Autowire(tmpBaz)
	assert.Equal(t, tmpBaz.MyFoo.(*Foo).Name, fooName)
	assert.Equal(t, tmpBaz.MyBaz.(*Bar).Name, barName)
	testFooName := "testFoo"
	testBarName := "testBar"
	Spy(tmpBaz, &Foo{Name: testFooName}, &Bar{Name: testBarName})
	assert.Equal(t, tmpBaz.MyFoo.(*Foo).Name, testFooName)
	assert.Equal(t, tmpBaz.MyBaz.(*Bar).Name, testBarName)
	assert.Equal(t, 0, len(Close()))
}
