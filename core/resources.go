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

type Resources struct {
    disposables []Disposable
}

func (rs *Resources) AddDisposable(r Disposable) {
    rs.disposables = append(rs.disposables, r)
}

func (rs *Resources) Init(r Resource) error {
	err := r.Init()
	if err != nil {
		return err
	}
    rs.disposables = append(rs.disposables, r)
	return nil
}

func (rs Resources) Dispose() {
    for i := len(rs.disposables) - 1; i >= 0; i-- {
        rs.disposables[i].Dispose()
    }
    rs.disposables = []Disposable{}
}

type SimpleResource struct {
    OnInit func() error
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

