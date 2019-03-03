package sdlutil

import (
	"fmt"
	"github.com/edudobay/go-snake/core"
	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
)

func SdlResource() core.Resource {
	return core.SimpleResource{
		OnInit: func() error {
			fmt.Println("initializing SDL")
			return sdl.Init(sdl.INIT_AUDIO | sdl.INIT_VIDEO)
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
			imgResult := img.Init(img.INIT_PNG)
			if imgResult&img.INIT_PNG != img.INIT_PNG {
				return fmt.Errorf("unable to init image lib: %s", img.GetError())
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
