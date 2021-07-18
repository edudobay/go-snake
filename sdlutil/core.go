package sdlutil

import (
	"errors"
	"fmt"
	"github.com/edudobay/go-snake/core"
	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
)

func SdlResource() core.Resource {
	return core.SimpleResource{
		OnInit: func() error {
			fmt.Println("initializing SDL")
			if err := sdl.Init(sdl.INIT_AUDIO | sdl.INIT_VIDEO); err != nil {
				return errors.New(fmt.Sprintf("error initializing SDL: %v", err))
			} else {
				return nil
			}
		},
		OnDispose: func() {
			fmt.Println("cleaning up SDL")
			sdl.Quit()
		},
	}
}

func ImgResource() core.Resource {
	return core.SimpleResource{
		OnInit: func() error {
			fmt.Println("initializing SDL Image library")
			if err := img.Init(img.INIT_PNG); err != nil {
				return errors.New(fmt.Sprintf("error initializing SDL_image: %v", err))
			} else {
				return nil
			}
		},
		OnDispose: func() {
			fmt.Println("cleaning up SDL Image library")
			img.Quit()
		},
	}
}
