package screen

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

type GameWindow struct {
	Window      *pixelgl.Window
	Screen      GameScreen
	SpriteSheet pixel.Picture
}

func (gw *GameWindow) Closed() bool {
	return gw.Window.Closed()
}

func (gw *GameWindow) Update() {
	newScreen := gw.Screen.Step(gw.SpriteSheet, gw.Window)
	gw.Window.Update()
	if newScreen != gw.Screen {
		gw.Screen = newScreen
		gw.Screen.Enter(gw.SpriteSheet, gw.Window)
	}
}

func NewGameWindow(ss pixel.Picture) *GameWindow {
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

	// FIXME
	screen := &InitialScreen{}
	//screen := NewMainGameScreen([]Player{
	//	Player{
	//		Name:          "Player1",
	//		HumanPlayer:   true,
	//		CharacterIcon: logical.V(0, 23),
	//	},
	//	Player{
	//		Name:          "Player2",
	//		HumanPlayer:   true,
	//		CharacterIcon: logical.V(1, 23),
	//	},
	//})
	screen.Enter(ss, win)

	return &GameWindow{
		Window:      win,
		Screen:      screen,
		SpriteSheet: ss,
	}
}
