package screen

import (
	"io"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"

	"github.com/bobtfish/mayhem/logical"
	"github.com/bobtfish/mayhem/render"
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

	ss := render.GetSpriteSheet(spriteReader)

	// FIXME
	//screen := &InitialScreen{}
	screen := NewMainGameScreen([]Player{
		Player{
			Name:          "Player1",
			HumanPlayer:   true,
			CharacterIcon: logical.V(0, 23),
		},
		Player{
			Name:          "Player2",
			HumanPlayer:   true,
			CharacterIcon: logical.V(1, 23),
		},
	})
	screen.Setup(ss, win)

	return &GameWindow{
		Window:      win,
		Screen:      screen,
		SpriteSheet: ss,
	}
}
