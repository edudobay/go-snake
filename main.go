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
const DefaultMap = "data/square.map"

type commandLineArgs struct {
	Level int    `arg:"-l,help:start at this level number"`
	Map   string `arg:"-m,help:set a custom map"`
}

func getArgs() commandLineArgs {
	var args commandLineArgs
	args.Level = DefaultLevel
	args.Map = DefaultMap

	arg.MustParse(&args)

	return args
}

func readKeys(game snake.Game) {
	for key := range game.KeyPresses() {
		fmt.Println("GAME: pressed key %v", key)
	}
}

func gameLoop(args commandLineArgs) {
	resources := new(core.Resources)
	defer resources.Dispose()

	game := snake.NewGame(args.Level)
	fmt.Printf("Game: %v\n", game)

	d, err := display.InitDisplay(resources)
	if err != nil {
		return
	}

	map_ := snake.ReadMap(args.Map)
	board := snake.CreateBoard(map_)

	d.DrawBoard(board)
	d.Update()

	go readKeys(game)

	quit := false
	for !quit {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				println("quit")
				quit = true

			case *sdl.KeyboardEvent:
				event := event.(*sdl.KeyboardEvent)
				if event.Type == sdl.KEYDOWN {
					game.OnKeyPressed(event.Keysym)
				}
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
