package main

import (
	"fmt"
	"github.com/edudobay/go-snake/snake"
	"github.com/veandco/go-sdl2/sdl"
	"time"
)

type mainController struct {
	app *application
}

func (c *mainController) HandleEvent(event sdl.Event) {
	switch event.(type) {
	case *sdl.QuitEvent:
		c.app.quit()

	case *sdl.KeyboardEvent:
		event := event.(*sdl.KeyboardEvent)
		if event.Type == sdl.KEYDOWN {
			c.keyPressed(event.Keysym)
		}
	}
}

func (c *mainController) OnTick(tick uint32) {
	elapsed := time.Duration(tick-c.app.LastTick) * time.Millisecond
	if elapsed >= 200*time.Millisecond {
		c.app.Game.Tick()
		c.app.LastTick = tick
	}
}

func (c *mainController) keyPressed(key sdl.Keysym) {
	switch key.Sym {
	case sdl.K_q:
		c.app.quit()

	case sdl.K_LEFT:
		c.app.Game.Move(snake.Left)
	case sdl.K_RIGHT:
		c.app.Game.Move(snake.Right)
	case sdl.K_UP:
		c.app.Game.Move(snake.Up)
	case sdl.K_DOWN:
		c.app.Game.Move(snake.Down)
	default:
		fmt.Printf("\x1b[1;32mGAME: pressed key %v\x1b[0;39m\n", key)
	}
}
