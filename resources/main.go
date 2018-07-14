package main

import "fmt"
import "errors"
import "github.com/edudobay/go-snake/core"

type Sdl struct {}

func (s Sdl) Init() error {
    fmt.Println("SDL init ok")
    return nil
}

func (s Sdl) Dispose() {
    fmt.Println("SDL disposed")
}

func imgResource() core.Resource {
    return core.SimpleResource{
        OnInit: func() error {
            // fmt.Println("IMG is warming up!!")
            fmt.Println("IMG is too hot!!")
            return errors.New("too hot")
        },
        OnDispose: func() {
            fmt.Println("IMG is going home!!")
        },
    }
}

func main() {
    resources := new(core.Resources)
    resources.Init(new(Sdl))
    resources.Init(imgResource())
    resources.Init(core.DisposableResource(func() {
        fmt.Println("just cleaning up the mess...")
    }))

    defer resources.Dispose()
}
