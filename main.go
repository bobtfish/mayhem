package main

import (
	"image"
	"os"

	_ "image/png"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"

	"math/rand"
)

const WIN_X = 1024
const WIN_Y = 768
const GRID_X = 15
const GRID_Y = 10
const CHAR_PIXELS = 64

var (
	ct   CharacterTypes
	grid *GameGrid
)

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

func getRectangle(x float64, y float64, size float64) *imdraw.IMDraw {
	imd := imdraw.New(nil)
	imd.Color = pickColour()
	imd.Push(pixel.V(x, y))
	imd.Push(pixel.V(x+size, y+size))
	imd.Rectangle(1)
	return imd
}

func drawMainBorderOne(inset, width int, win *pixelgl.Window) {
	imd := imdraw.New(nil)
	imd.Color = pickColour()
	// Bottom left
	imd.Push(pixel.V(float64(inset), float64(WIN_Y-(CHAR_PIXELS*(GRID_Y+1))+inset)))
	// Top right
	imd.Push(pixel.V(float64(WIN_X-inset), float64(WIN_Y-inset)))
	imd.Rectangle(float64(width))
	imd.Draw(win)
}

func drawBorder(win *pixelgl.Window) {
	drawMainBorderOne(2, 1, win)
	drawMainBorderOne(8, 2, win)
	drawMainBorderOne(16, 4, win)
	drawMainBorderOne(24, 6, win)

	imd := imdraw.New(nil)
	imd.Color = pickColour()
	imd.Push(pixel.V(2, 2))
	imd.Push(pixel.V(WIN_X-2, WIN_Y-(CHAR_PIXELS*(GRID_Y+1))-2))
	imd.Rectangle(1)
	imd.Draw(win)
}

func run() {
	ct = LoadCharacters("characters.yaml")
	grid = MakeGameGrid(GRID_X, GRID_Y)

	cfg := pixelgl.WindowConfig{
		Title:  "Pixel Rocks!",
		Bounds: pixel.R(0, 0, WIN_X, WIN_Y),
		VSync:  true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}
	/*win.Clear(colornames.Skyblue)

	pic, err := loadPicture("chaosa.png")
	if err != nil {
		panic(err)
	}
	sprite := pixel.NewSprite(pic, pic.Bounds())
	sprite.Draw(win, pixel.IM.Moved(win.Bounds().Center()))
	*/

	drawBorder(win)

	grid.Draw(CHAR_PIXELS/2, WIN_Y-(CHAR_PIXELS*GRID_Y+CHAR_PIXELS/2), win)
	for !win.Closed() {
		win.Update()
	}
}

func main() {
	pixelgl.Run(run)
}
