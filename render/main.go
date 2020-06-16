package render

import (
	_ "image/png"
	"io"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"

	"github.com/bobtfish/mayhem/logical"
)

const WIN_X = 1024
const WIN_Y = 768
const GRID_X = 15
const GRID_Y = 10
const CHAR_PIXELS = 64
const SPRITE_SIZE = 16

type GameWindow struct {
	Window      *pixelgl.Window
	Screen      GameScreen
	SpriteSheet pixel.Picture
}

func (gw *GameWindow) Closed() bool {
	return gw.Window.Closed()
}

func (gw *GameWindow) Update() {
	gw.Screen.Draw(gw.Window)
	gw.Window.Update()
}

func NewGameWindow(spriteReader io.Reader) *GameWindow {
	title := "Mayhem!"

	cfg := pixelgl.WindowConfig{
		Title:  title,
		Bounds: pixel.R(0, 0, WIN_X, WIN_Y),
		VSync:  true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	ss := GetSpriteSheet(spriteReader)

	return &GameWindow{
		Window: win,
		Screen: &MainScreen{
			SpriteSheet: ss,
		},
		SpriteSheet: ss,
		//		SpriteDrawer: NewSpriteDrawer(ss, logical.V(0, CHAR_PIXELS)),
	}
}

type GameScreen interface {
	Draw(*pixelgl.Window)
}

type MainScreen struct {
	SpriteSheet  pixel.Picture
	SpriteDrawer *SpriteDrawer
	TextDrawer   *SpriteDrawer
	Drawn        bool
}

func (screen *MainScreen) Draw(win *pixelgl.Window) {
	if !screen.Drawn {
		win.Clear(colornames.Black)
		sd := NewSpriteDrawer(screen.SpriteSheet, logical.V(0, CHAR_PIXELS))
		drawMainBorder(win, sd)
		sd.WinConverter.Offset = logical.V(CHAR_PIXELS/2, CHAR_PIXELS/2+CHAR_PIXELS)
		screen.SpriteDrawer = sd
		screen.TextDrawer = NewTextDrawer(screen.SpriteSheet, sd.WinConverter.Offset)
		screen.TextDrawer.DrawText("Hello", logical.V(1, 8), win)
	}
	screen.Drawn = true
}

func drawMainBorder(win *pixelgl.Window, sd *SpriteDrawer) {
	batch := sd.GetNewBatch()
	// Bottom left
	sd.DrawSprite(logical.V(6, 20), logical.V(0, 0), batch)
	// Bottom row
	for i := 1; i < 15; i++ {
		sd.DrawSprite(logical.V(7, 20), logical.V(i, 0), batch)
	}
	// Bottom right
	sd.DrawSprite(logical.V(8, 20), logical.V(15, 0), batch)
	// LHS and RHS
	for i := 1; i < 10; i++ {
		sd.DrawSprite(logical.V(5, 20), logical.V(0, i), batch)
		sd.DrawSprite(logical.V(9, 20), logical.V(15, i), batch)
	}
	// Top Left
	sd.DrawSprite(logical.V(2, 20), logical.V(0, 10), batch)
	// Top row
	for i := 1; i < 15; i++ {
		sd.DrawSprite(logical.V(3, 20), logical.V(i, 10), batch)
	}
	// Top right
	sd.DrawSprite(logical.V(4, 20), logical.V(15, 10), batch)
	batch.Draw(win)

	imd := imdraw.New(nil)
	imd.Color = colornames.Blue
	imd.Push(pixel.V(2, 2))
	imd.Push(pixel.V(WIN_X-2, WIN_Y-(CHAR_PIXELS*(GRID_Y+1))-2))
	imd.Rectangle(1)
	imd.Draw(win)
}
