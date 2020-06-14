package main

import (
	"fmt"
	"image"
	"os"
	"time"

	_ "image/png"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"

	"math/rand"

	"github.com/bobtfish/mayhem/render"
)

const WIN_X = 1024
const WIN_Y = 768
const GRID_X = 15
const GRID_Y = 10
const CHAR_PIXELS = 64
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

func drawMainBorderOne(inset, width int, win *pixelgl.Window) {
	imd := imdraw.New(nil)
	imd.Color = colornames.Blue
	// Bottom left
	imd.Push(pixel.V(float64(inset), float64(WIN_Y-(CHAR_PIXELS*(GRID_Y+1))+inset)))
	// Top right
	imd.Push(pixel.V(float64(WIN_X-inset), float64(WIN_Y-inset)))
	imd.Rectangle(float64(width))
	imd.Draw(win)
}

func drawMainBorder(win *pixelgl.Window) {
	drawMainBorderOne(2, 1, win)
	drawMainBorderOne(8, 2, win)
	drawMainBorderOne(16, 4, win)
	drawMainBorderOne(24, 6, win)

	imd := imdraw.New(nil)
	imd.Color = colornames.Blue
	imd.Push(pixel.V(2, 2))
	imd.Push(pixel.V(WIN_X-2, WIN_Y-(CHAR_PIXELS*(GRID_Y+1))-2))
	imd.Rectangle(1)
	imd.Draw(win)
}

func drawMainWindow(win *pixelgl.Window, grid *GameGrid, ss pixel.Picture) {
	win.Clear(colornames.Black)
	drawMainBorder(win)
	grid.Draw(win, ss)
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
	drawMainWindow(gw.Window, grid, ss)

	QsecondTicks := 0
	frames := 0
	Qsecond := time.Tick(time.Second / 4)

	for !gw.Closed() {
		drawMainWindow(gw.Window, grid, ss)

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
