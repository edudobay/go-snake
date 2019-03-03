package main

import (
	"fmt"
	"github.com/alexflint/go-arg"
	"github.com/edudobay/go-snake/core"
	"github.com/edudobay/go-snake/display"
	"github.com/edudobay/go-snake/sdlutil"
	"github.com/edudobay/go-snake/snake"
	"github.com/veandco/go-sdl2/sdl"
	"runtime"
)

const DefaultLevel = 7
const DefaultMap = "data/square.map"

type quitSignal struct{}

type commandLineArgs struct {
	Level int    `arg:"-l,help:start at this level number"`
	Map   string `arg:"-m,help:set a custom map"`
}

type application struct {
	Game       *snake.Game
	Display    *display.Display
	Map        snake.GameMap
	Board      *snake.Board
	Quit       chan quitSignal
	Events     chan sdl.Event
	KeyPresses chan sdl.Keysym
}

func getArgs() commandLineArgs {
	var args commandLineArgs
	args.Level = DefaultLevel
	args.Map = DefaultMap

	arg.MustParse(&args)

	return args
}

func (app application) readKeys() {
	for key := range app.KeyPresses {
		switch key.Sym {
		case sdl.K_q:
			fmt.Println("quit")
			app.Quit <- quitSignal{}

		case sdl.K_LEFT:
			app.Game.Move(snake.Left)
		case sdl.K_RIGHT:
			app.Game.Move(snake.Right)
		case sdl.K_UP:
			app.Game.Move(snake.Up)
		case sdl.K_DOWN:
			app.Game.Move(snake.Down)
		default:
			fmt.Printf("\x1b[1;32mGAME: pressed key %v\x1b[0;39m\n", key)
		}
	}
}

func (app application) handleEvent(event sdl.Event) {
	switch event.(type) {
	case *sdl.QuitEvent:
		println("quit")
		app.Quit <- quitSignal{}

	case *sdl.KeyboardEvent:
		event := event.(*sdl.KeyboardEvent)
		if event.Type == sdl.KEYDOWN {
			app.KeyPresses <- event.Keysym
		}
	}
}

func (app application) handleEvents() {
	for event := range app.Events {
		app.handleEvent(event)
	}
}

func processSdlEvents(events chan<- sdl.Event, quit <-chan quitSignal) {
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
	app.KeyPresses = make(chan sdl.Keysym, 100)

	go app.readKeys()
	go app.handleEvents()

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
	resources.Init(sdlutil.SdlResource())
	resources.Init(sdlutil.ImgResource())

	defer resources.Dispose()

	gameLoop(args)
}
