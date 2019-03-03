package main

import (
	"fmt"
	"runtime"
	"github.com/alexflint/go-arg"
	"github.com/edudobay/go-snake/core"
	"github.com/edudobay/go-snake/display"
	sdlCore "github.com/edudobay/go-snake/sdl"
	"github.com/edudobay/go-snake/snake"
	"github.com/veandco/go-sdl2/sdl"
)

const DefaultLevel = 7
const DefaultMap = "data/square.map"

type quitSignal struct{}
type quitReceiver <-chan quitSignal
type quitSender chan<- quitSignal

type commandLineArgs struct {
	Level int    `arg:"-l,help:start at this level number"`
	Map   string `arg:"-m,help:set a custom map"`
}

type application struct {
	Game 		snake.Game
	Display  	*display.Display
	Map			snake.GameMap
	Board		snake.Board
	Quit		chan quitSignal
	Events		chan sdl.Event
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
		fmt.Printf("\x1b[1;32mGAME: pressed key %v\x1b[0;39m\n", key)
	}
}

func handleGameEvent(game snake.Game, event sdl.Event, quit quitSender) {
	switch event.(type) {
	case *sdl.QuitEvent:
		println("quit")
		quit <- quitSignal{}

	case *sdl.KeyboardEvent:
		event := event.(*sdl.KeyboardEvent)
		if event.Type == sdl.KEYDOWN {
			game.OnKeyPressed(event.Keysym)
		}
	}
}

func processSdlEvents(events chan<- sdl.Event, quit quitReceiver) {
	alive := true

	go func() {
		<-quit
		alive = false
	}()

	for alive {
		if event := sdl.PollEvent(); event != nil {
			events <- event
		}

		// Yield execution to other goroutines
		runtime.Gosched()
	}
}

func gameLoop(args commandLineArgs) {
	resources := new(core.Resources)
	defer resources.Dispose()

	app := new(application)

	app.Game = snake.NewGame(args.Level)
	fmt.Printf("Game: %v\n", app.Game)

	var err error
	app.Display, err = display.InitDisplay(resources)
	if err != nil {
		return
	}

	app.Map = snake.ReadMap(args.Map)
	app.Board = snake.CreateBoard(app.Map)

	app.Display.DrawBoard(app.Board)
	app.Display.Update()

	app.Quit = make(chan quitSignal)
	app.Events = make(chan sdl.Event, 100)

	go readKeys(app.Game)

	go func() {
		for event := range app.Events {
			handleGameEvent(app.Game, event, app.Quit)
		}
	}()

	processSdlEvents(app.Events, app.Quit)
}

func main() {

	args := getArgs()

	fmt.Printf("Starting at level %d\n", args.Level)
	if args.Map != "" {
		fmt.Printf("Map: %v\n", args.Map)
	} else {
		fmt.Println("Map is nil")
	}

	resources := new(core.Resources)
	resources.Init(sdlCore.SdlResource())
	resources.Init(sdlCore.ImgResource())

	defer resources.Dispose()

	gameLoop(args)
}
