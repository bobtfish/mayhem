package render

import (
	"fmt"
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

func (gw *GameWindow) MaybeNextScreen() {
	if gw.Screen.Finished() {
		screen := gw.Screen.NextScreen()
		screen.Setup(gw.SpriteSheet, gw.Window)
		gw.Screen = screen
	}
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

	screen := &MainScreen{}
	screen.Setup(ss, win)

	return &GameWindow{
		Window:      win,
		Screen:      screen,
		SpriteSheet: ss,
	}
}

type GameScreen interface {
	Setup(pixel.Picture, *pixelgl.Window)
	Draw(*pixelgl.Window)
	Finished() bool
	NextScreen() GameScreen
}

type ScreenBasics struct {
	SpriteSheet  pixel.Picture
	SpriteDrawer *SpriteDrawer
	TextDrawer   *SpriteDrawer
}

func (screen *ScreenBasics) Setup(ss pixel.Picture, win *pixelgl.Window) {
	screen.SpriteSheet = ss
	sd := NewSpriteDrawer(ss, logical.V(0, CHAR_PIXELS))
	drawMainBorder(win, sd)
	sd.WinConverter.Offset = logical.V(CHAR_PIXELS/2, CHAR_PIXELS/2+CHAR_PIXELS)
	screen.SpriteDrawer = sd
	screen.TextDrawer = NewTextDrawer(ss, sd.WinConverter.Offset)
}

type MainScreen struct {
	ScreenBasics
	DrawnFirst         bool
	WizardCount        int
	DrawnSecond        bool
	ComputerDifficulty int
}

func captureNumKey(win *pixelgl.Window) int {
	if win.JustPressed(pixelgl.Key0) {
		return 0
	}
	if win.JustPressed(pixelgl.Key1) {
		return 1
	}
	if win.JustPressed(pixelgl.Key2) {
		return 2
	}
	if win.JustPressed(pixelgl.Key3) {
		return 3
	}
	if win.JustPressed(pixelgl.Key3) {
		return 4
	}
	if win.JustPressed(pixelgl.Key5) {
		return 5
	}
	if win.JustPressed(pixelgl.Key6) {
		return 6
	}
	if win.JustPressed(pixelgl.Key7) {
		return 7
	}
	if win.JustPressed(pixelgl.Key8) {
		return 8
	}
	if win.JustPressed(pixelgl.Key9) {
		return 9
	}
	return -1
}

func (screen *MainScreen) Draw(win *pixelgl.Window) {
	if !screen.DrawnFirst {
		win.Clear(colornames.Black)

		sd := NewSpriteDrawer(screen.SpriteSheet, logical.V(0, CHAR_PIXELS))
		drawMainBorder(win, sd)
		sd.WinConverter.Offset = logical.V(CHAR_PIXELS/2, CHAR_PIXELS/2+CHAR_PIXELS)
		screen.SpriteDrawer = sd
		screen.TextDrawer = NewTextDrawer(screen.SpriteSheet, sd.WinConverter.Offset)
		screen.TextDrawer.DrawText("  MAYHEM - Remake of Chaos", logical.V(0, 9), win)
		screen.TextDrawer.DrawText("         By bobtfish", logical.V(0, 8), win)
		screen.TextDrawer.DrawText("How many wizards?", logical.V(0, 6), win)
		screen.TextDrawer.DrawText("(Press 2 to 8)", logical.V(0, 5), win)
		screen.DrawnFirst = true
	} else {
		if screen.WizardCount == 0 {
			c := captureNumKey(win)
			if c >= 2 && c <= 8 {
				screen.WizardCount = c
				screen.TextDrawer.DrawText(fmt.Sprintf("%d", c), logical.V(18, 6), win)
			}
		} else {
			if !screen.DrawnSecond {
				screen.TextDrawer.DrawText("Level of computer wizards?", logical.V(0, 3), win)
				screen.TextDrawer.DrawText("(Press 1 to 8)", logical.V(0, 2), win)
				screen.DrawnSecond = true
			} else {
				c := captureNumKey(win)
				if c >= 1 && c <= 8 {
					screen.ComputerDifficulty = c
					screen.TextDrawer.DrawText(fmt.Sprintf("%d", c), logical.V(27, 3), win)
				}
			}
		}
	}
}

func (screen *MainScreen) NextScreen() GameScreen {
	return &MainScreen{}
}

func (screen *MainScreen) Finished() bool {
	if screen.ComputerDifficulty > 0 && screen.WizardCount > 0 {
		return true
	}
	return false
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
