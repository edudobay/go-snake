package core

type Disposable interface {
	Dispose()
}

type Allocatable interface {
	Init() error
}

type Resource interface {
	Allocatable
	Disposable
}

type resourcesStack struct {
	top  Disposable
	rest *resourcesStack
}

func (stack *resourcesStack) Append(d Disposable) *resourcesStack {
	return &resourcesStack{d, stack}
}

func (stack *resourcesStack) Pop() (Disposable, *resourcesStack) {
	return stack.top, stack.rest
}

type Resources struct {
	stack *resourcesStack
}

func (rs *Resources) AddDisposable(d Disposable) {
	rs.stack = rs.stack.Append(d)
}

func (rs *Resources) Init(r Resource) error {
	err := r.Init()
	if err != nil {
		return err
	}
	rs.stack = rs.stack.Append(r)
	return nil
}

func (rs *Resources) Dispose() {
	var d Disposable
	for stack := rs.stack; stack != nil; {
		d, stack = stack.Pop()
		d.Dispose()
	}
	rs.stack = nil
}

type SimpleResource struct {
	OnInit    func() error
	OnDispose func()
}

func DisposableResource(dispose func()) Resource {
	return SimpleResource{
		OnInit: func() error {
			return nil
		},
		OnDispose: dispose,
	}
}

func (d SimpleResource) Init() error {
	return d.OnInit()
}

func (d SimpleResource) Dispose() {
	d.OnDispose()
}
