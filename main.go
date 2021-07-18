package main

import (
	"fmt"
	"github.com/alexflint/go-arg"
	"github.com/edudobay/go-snake/core"
	"github.com/edudobay/go-snake/display"
	"github.com/edudobay/go-snake/sdlutil"
	"github.com/edudobay/go-snake/snake"
	"github.com/veandco/go-sdl2/sdl"
	"log"
	"runtime"
	"time"
)

const DefaultLevel = 7
const DefaultMap = "data/square.map"

type quitSignal struct{}

type commandLineArgs struct {
	Level int    `arg:"-l,help:start at this level number"`
	Map   string `arg:"-m,help:set a custom map"`
}

type controller interface {
	HandleEvent(sdl.Event)
	OnTick(uint32)
}

type eventToDispatch struct {
	event      sdl.Event
	controller controller
}

type application struct {
	Game             *snake.Game
	Display          *display.Display
	Map              snake.GameMap
	Board            *snake.Board
	Quit             chan quitSignal
	Events           chan sdl.Event
	DispatchQueue    chan eventToDispatch
	ActiveController controller
	LastTick         uint32
}

func getArgs() commandLineArgs {
	var args commandLineArgs
	args.Level = DefaultLevel
	args.Map = DefaultMap

	arg.MustParse(&args)

	return args
}

func (app application) quit() {
	println("quit")
	app.Quit <- quitSignal{}
	close(app.Quit)
}

func (app application) GetController() controller {
	controller := app.ActiveController
	if controller == nil {
		panic("no active controller to dispatch to")
	}
	return controller
}

func (app application) handleEvent(event sdl.Event) {
	controller := app.GetController()
	app.DispatchQueue <- eventToDispatch{event, controller}
}

func (app application) dispatchEvents() {
	for {
		select {
		case e := <-app.DispatchQueue:
			e.controller.HandleEvent(e.event)

		case <-app.Quit:
			fmt.Println("QUIT: detected from dispatch queue")
			return

		case <-time.After(20 * time.Millisecond):
			app.GetController().OnTick(sdl.GetTicks())
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

	app.DispatchQueue = make(chan eventToDispatch, 64)

	app.Map = snake.ReadMap(args.Map)
	app.Board = snake.CreateBoard(app.Map)

	snakeX, snakeY := app.Board.Center()
	app.Board.PutSnake(snakeX, snakeY, 4, snake.Left)

	app.Game = snake.NewGame(args.Level, app.Board)
	fmt.Printf("Game: %v\n", app.Game)

	var err error
	app.Display, err = display.InitDisplay(resources)
	if err != nil {
		return
	}

	app.Display.ClearWindow()

	app.Display.DrawBoard(app.Board)
	app.Display.Update()

	app.Quit = make(chan quitSignal)
	app.Events = make(chan sdl.Event, 100)

	app.LastTick = sdl.GetTicks()
	app.ActiveController = &mainController{app: app}
	app.Game.Move(snake.Right)

	go app.dispatchEvents()
	go app.handleEvents()
	go func() {
		for {
			select {
			case <-app.Board.Updates():
				app.Display.DrawBoard(app.Board)
				app.Display.Update()
			case <-app.Quit:
				return
			}
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
	if err := resources.Init(sdlutil.SdlResource()); err != nil {
		log.Panicln(err)
	}
	if err := resources.Init(sdlutil.ImgResource()); err != nil {
		log.Panicln(err)
	}

	defer resources.Dispose()

	gameLoop(args)
}
