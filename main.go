package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"

	"math/rand"

	"github.com/bobtfish/mayhem/logical"
	"github.com/bobtfish/mayhem/render"
)

const GRID_X = 15
const GRID_Y = 10
const SPRITE_SIZE = 16

func pickColour() pixel.RGBA {
	return pixel.RGB(rand.Float64(), rand.Float64(), rand.Float64())
}

func drawHydra(ss pixel.Picture, win *pixelgl.Window) {
	rect := pixel.R(0, 16, 16, 32)
	sprite := pixel.NewSprite(ss, rect)
	mat := pixel.IM
	mat = mat.Moved(win.Bounds().Center())
	mat = mat.ScaledXY(win.Bounds().Center(), pixel.V(4, 4))
	sprite.Draw(win, mat)
}

func placeCharactersTest(grid *GameGrid, ct CharacterTypes) {
	x := 0
	y := 0
	for k := range ct {
		grid.PlaceCharacter(logical.V(x, y), ct.NewCharacter(k))
		x++
		if x == 15 {
			x = 0
			y++
		}
	}
}

func run() {
	data, err := base64.StdEncoding.DecodeString(sprite_sheet_base64)
	if err != nil {
		panic(err)
	}
	r := bytes.NewReader(data)

	sd := render.NewSpriteDrawer(r)

	ct := LoadCharacterTemplates()
	grid := MakeGameGrid(logical.V(GRID_X, GRID_Y))

	title := "Mayhem!"

	gw := render.NewGameWindow(sd)
	placeCharactersTest(grid, ct)

	QsecondTicks := 0
	frames := 0
	Qsecond := time.Tick(time.Second / 4)

	for !gw.Closed() {
		batch := grid.DrawBatch(&sd)
		batch.Draw(gw.Window)

		gw.Update()

		frames++
		select {
		case <-Qsecond:
			grid.AnimationTick()
			QsecondTicks++
			if QsecondTicks == 4 {
				gw.Window.SetTitle(fmt.Sprintf("%s | FPS: %d", title, frames*4))
				frames = 0
				QsecondTicks = 0
			}
		default:
		}
	}
}

func main() {
	pixelgl.Run(run)
}
