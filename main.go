package main

import (
	"fmt"
	"github.com/alexflint/go-arg"
	"github.com/edudobay/go-snake/core"
	sdlCore "github.com/edudobay/go-snake/sdl"
	"github.com/edudobay/go-snake/snake"
	"github.com/veandco/go-sdl2/sdl"
)

const DefaultLevel = 7

type commandLineArgs struct {
	Level int    `arg:"-l,help:start at this level number"`
	Map   string `arg:"-m,help:set a custom map"`
}

func getArgs() commandLineArgs {
	var args commandLineArgs
	args.Level = DefaultLevel

	arg.MustParse(&args)

	return args
}

func gameLoop(args commandLineArgs) {
	window, err := sdl.CreateWindow(
		"test",
		sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		800, 600,
		sdl.WINDOW_SHOWN)
	if err != nil {
		panic(err)
	}
	defer window.Destroy()

	surface, err := window.GetSurface()
	if err != nil {
		panic(err)
	}

	surface.FillRect(nil, 0)

	rect := sdl.Rect{0, 0, 200, 200}
	surface.FillRect(&rect, 0xffff0000)
	window.UpdateSurface()

	game := snake.NewGame(args.Level)
	fmt.Printf("Game: %v\n", game)

	quit := false
	for !quit {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				println("quit")
				quit = true
				break
			}

			sdl.Delay(20)
			window.UpdateSurface()
		}
	}

	//	sdl.Delay(100)
	//	window.UpdateSurface()
	//
	//	sdl.Delay(2000)
	//
	//	rect2 := sdl.Rect{0, 0, 100, 100}
	//	surface.FillRect(&rect2, 0xff00ff00)
	//	window.UpdateSurface()
	//
	//	sdl.Delay(1000)

}

func main() {

	args := getArgs()

	fmt.Printf("Starting at level %d\n", args.Level)
	if args.Map != "" {
		fmt.Printf("Map: %v", args.Map)
	} else {
		fmt.Println("Map is nil")
	}

	resources := new(core.Resources)
	resources.Init(sdlCore.SdlResource())
	resources.Init(sdlCore.ImgResource())

	defer resources.Dispose()

	gameLoop(args)
}
