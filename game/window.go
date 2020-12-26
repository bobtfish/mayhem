package game

import (
	screeniface "github.com/bobtfish/mayhem/screen/iface"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

type GameWindow struct {
	Window      *pixelgl.Window
	SpriteSheet pixel.Picture
}

func (gw *GameWindow) GetWindow() *pixelgl.Window {
	return gw.Window
}

func (gw *GameWindow) GetSpriteSheet() pixel.Picture {
	return gw.SpriteSheet
}

func (gw *GameWindow) Closed() bool {
	return gw.Window.Closed()
}

func (gw *GameWindow) Update(screen screeniface.GameScreen) screeniface.GameScreen {
	newScreen := screen.Step(gw)
	gw.Window.Update()
	if newScreen != screen {
		newScreen.Enter(gw)
		return newScreen
	}
	return screen
}

func NewWindow(ss pixel.Picture) *GameWindow {
	title := "Mayhem!"

	cfg := pixelgl.WindowConfig{
		Title:  title,
		Bounds: pixel.R(0, 0, WinX, WinY),
		VSync:  true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	return &GameWindow{
		Window:      win,
		SpriteSheet: ss,
	}
}
