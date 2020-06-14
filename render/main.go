package render

import (
	_ "image/png"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

const WIN_X = 1024
const WIN_Y = 768
const GRID_X = 15
const GRID_Y = 10
const CHAR_PIXELS = 64
const SPRITE_SIZE = 16

type GameWindow struct {
	Window *pixelgl.Window
	Screen GameScreen
}

func (gw *GameWindow) Closed() bool {
	return gw.Window.Closed()
}

func (gw *GameWindow) Update() {
	gw.Screen.Draw(gw.Window)
	gw.Window.Update()
}

type GameScreen interface {
	Draw(*pixelgl.Window)
}

type MainScreen struct {
	Drawn bool
}

func (screen *MainScreen) Draw(win *pixelgl.Window) {
	if !screen.Drawn {
		win.Clear(colornames.Black)
		drawMainBorder(win)
	}
	screen.Drawn = true
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

func NewGameWindow() *GameWindow {
	title := "Mayhem!"

	cfg := pixelgl.WindowConfig{
		Title:  title,
		Bounds: pixel.R(0, 0, WIN_X, WIN_Y),
		//  VSync:  true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	return &GameWindow{
		Window: win,
		Screen: &MainScreen{},
	}
}
