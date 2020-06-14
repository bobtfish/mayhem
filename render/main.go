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
	Window       *pixelgl.Window
	Screen       GameScreen
	SpriteDrawer SpriteDrawer
}

func (gw *GameWindow) Closed() bool {
	return gw.Window.Closed()
}

func (gw *GameWindow) Update() {
	gw.Screen.Draw(gw.Window, &gw.SpriteDrawer)
	gw.Window.Update()
}

func NewGameWindow(sd SpriteDrawer) *GameWindow {
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
		Window:       win,
		Screen:       &MainScreen{},
		SpriteDrawer: sd,
	}
}

type GameScreen interface {
	Draw(*pixelgl.Window, *SpriteDrawer)
}

type MainScreen struct {
	Drawn bool
}

func (screen *MainScreen) Draw(win *pixelgl.Window, sd *SpriteDrawer) {
	if !screen.Drawn {
		win.Clear(colornames.Black)
		drawMainBorder(win, sd)
	}
	screen.Drawn = true
}

func drawMainBorder(win *pixelgl.Window, sd *SpriteDrawer) {
	batch := sd.GetNewBatch()
	// Bottom left
	sd.DrawSprite(6, 20, 0, 1, batch)
	// Bottom row
	for i := 1; i < 15; i++ {
		sd.DrawSprite(7, 20, i, 1, batch)
	}
	// Bottom right
	sd.DrawSprite(8, 20, 15, 1, batch)
	// LHS and RHS
	for i := 2; i <= 10; i++ {
		sd.DrawSprite(5, 20, 0, i, batch)
		sd.DrawSprite(9, 20, 15, i, batch)
	}
	// Top Left
	sd.DrawSprite(2, 20, 0, 11, batch)
	// Top row
	for i := 1; i < 15; i++ {
		sd.DrawSprite(3, 20, i, 11, batch)
	}
	// Top right
	sd.DrawSprite(4, 20, 15, 11, batch)
	batch.Draw(win)

	imd := imdraw.New(nil)
	imd.Color = colornames.Blue
	imd.Push(pixel.V(2, 2))
	imd.Push(pixel.V(WIN_X-2, WIN_Y-(CHAR_PIXELS*(GRID_Y+1))-2))
	imd.Rectangle(1)
	imd.Draw(win)
}
