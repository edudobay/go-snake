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

//

type SimpleResource struct {
	OnInit    func() error
	OnDispose func()
}

func (d SimpleResource) Init() error {
	return d.OnInit()
}

func (d SimpleResource) Dispose() {
	d.OnDispose()
}

//

func DisposableResource(dispose func()) Resource {
	return SimpleResource{
		OnInit: func() error {
			return nil
		},
		OnDispose: dispose,
	}
}
