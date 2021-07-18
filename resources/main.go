package main

import "fmt"
import "errors"
import "github.com/edudobay/go-snake/core"

type Sdl struct{}

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
	if err := resources.Init(new(Sdl)); err != nil {
		fmt.Printf("error initializing SDL: %v\n", err)
	}
	if err := resources.Init(imgResource()); err != nil {
		fmt.Printf("error initializing image: %v\n", err)
	}
	if err := resources.Init(core.DisposableResource(func() {
		fmt.Println("just cleaning up the mess...")
	})); err != nil {
		fmt.Printf("error initializing disposable: %v\n", err)
	}

	defer resources.Dispose()
}
