package main

import (
	"fmt"
	"github.com/alexflint/go-arg"
	"github.com/edudobay/go-snake/core"
	"github.com/edudobay/go-snake/display"
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
	resources := new(core.Resources)
	defer resources.Dispose()

	game := snake.NewGame(args.Level)
	fmt.Printf("Game: %v\n", game)

	if err := display.InitDisplay(resources); err != nil {
		return
	}

	quit := false
	for !quit {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				println("quit")
				quit = true
				break
			}
		}
	}

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
