package fake

// Passer represents interface
type Passer interface {
	Pass()
}

// Foo represents named struct
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

// Baz represents named struct
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

// Passer method is a simple Getter
func (q *Qus) Passer() Passer {
	return q.passer
}

// NotFoundTagDependency represents named struct
type NotFoundTagDependency struct {
	passer Passer `autowire:"fake/FooBaz"`
}

// InvalidInterface represents interface
type InvalidInterface interface {
	Invalid()
}

// Bor represent named struct
type Bor struct {
	Passer InvalidInterface `autowire:"fake/Foo"`
}
