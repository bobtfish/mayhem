package main

import (
	"fmt"
	"image"
	"os"
	"time"

	_ "image/png"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"

	"math/rand"

	"github.com/bobtfish/mayhem/render"
)

const GRID_X = 15
const GRID_Y = 10
const SPRITE_SIZE = 16

func loadPicture(path string) (pixel.Picture, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}
	return pixel.PictureDataFromImage(img), nil
}

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
		grid.PlaceCharacter(x, y, ct.NewCharacter(k))
		x++
		if x == 15 {
			x = 0
			y++
		}
	}
}

func run() {
	ss, err := loadPicture("sprite_sheet.png")
	if err != nil {
		panic(err)
	}

	ct := LoadCharacterTemplates("characters.yaml")
	grid := MakeGameGrid(GRID_X, GRID_Y)

	title := "Mayhem!"

	gw := render.NewGameWindow()

	placeCharactersTest(grid, ct)
	grid.Draw(gw.Window, ss)

	QsecondTicks := 0
	frames := 0
	Qsecond := time.Tick(time.Second / 4)

	for !gw.Closed() {
		grid.Draw(gw.Window, ss)

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
