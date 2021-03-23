package fake

type Passer interface {
	Pass()
}

// Foo represent named struct
type Foo struct {
	Name       string
	CloseCalls int
}

// Pass method
func (f Foo) Pass() {
}

// Close method of Foo
func (f *Foo) Close() error {
	f.CloseCalls++
	return nil
}

// Bar represent named struct
type Bar struct {
	myFoo *Foo `autowire:""`
}

// SetMyFoo method is a simple Setter
func (b *Bar) SetMyFoo(myFoo *Foo) {
	b.myFoo = myFoo
}

// MyFoo method is a simple Getter
func (b *Bar) MyFoo() *Foo {
	return b.myFoo
}

// Bar represent named struct
type Baz struct {
	MyFoo *Foo `autowire:""`
}

// Qux represent named struct autowiring interface instance into exported field
type Qux struct {
	Passer Passer `autowire:"fake/Foo"`
}

// Qus represent named struct autowiring interface instance into unexported field
type Qus struct {
	passer Passer `autowire:"fake/Foo"`
}

// SetMyFoo method is a simple Setter
func (q *Qus) SetPasser(passer Passer) {
	q.passer = passer
}

// Passer method is a simple Getter
func (q *Qus) Passer() Passer {
	return q.passer
}

type NotFoundTagDependency struct {
	passer Passer `autowire:"fake/FooBaz"`
}

// SetMyFoo method is a simple Setter
func (n *NotFoundTagDependency) SetPasser(passer Passer) {
	n.passer = passer
}

// Passer method is a simple Getter
func (n *NotFoundTagDependency) Passer() Passer {
	return n.passer
}

// InvalidInterface represents interface
type InvalidInterface interface {
	Invalid()
}

// Bor represent named struct
type Bor struct {
	Passer InvalidInterface `autowire:"fake/Foo"`
}
