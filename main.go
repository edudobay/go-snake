package main

import "os"
import "io"
import "bufio"
import "unicode"

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

	readMap(args.Map)

	d, err := display.InitDisplay(resources)
	if err != nil {
		return
	}

	d.DrawSprite(display.SpriteFood, 80, 80)
	d.Update()

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

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func checkN(n, actualN int, e error) {
	check(e)
	if n != actualN {
		panic(fmt.Errorf("%d != expected %d", actualN, n))
	}
}

func assert(condition bool, msg string) {
	if !condition {
		panic(fmt.Errorf("assertion failed: %s", msg))
	}
}

func readMap(filename string) []int8 {

	f, err := os.Open(filename)
	check(err)

	var width, height int
	n, err := fmt.Fscanf(f, "%d %d", &width, &height)
	check(err)
	assert(n == 2, "invalid size header")

	size := width * height
	assert(size > 0, "invalid width/height")

	map_ := make([]int8, size)

	reader := bufio.NewReader(f)

	for i := 0; i < size; {
		b, _, err := reader.ReadRune()
		if err == io.EOF {
			panic("premature end of file")
		}

		if unicode.IsSpace(b) {
			continue
		}

		switch b {
		case '.':
			map_[i] = (int8)(display.SpriteNone)
		case '#':
			map_[i] = (int8)(display.SpriteWall)
		case 'x':
			map_[i] = (int8)(display.SpriteInvalid)
		default:
			panic("invalid char in map")
		}

		i++
	}

	return map_
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
