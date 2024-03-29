package game

import (
	"fmt"

	"github.com/bobtfish/mayhem/grid"
	"github.com/bobtfish/mayhem/player"
	screeniface "github.com/bobtfish/mayhem/screen/iface"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

type Window struct {
	Window      *pixelgl.Window
	SpriteSheet pixel.Picture
	Grid        *grid.GameGrid
	Players     []*player.Player
	LawRating   int
}

func (gw *Window) GetWindow() *pixelgl.Window {
	return gw.Window
}

func (gw *Window) GetSpriteSheet() pixel.Picture {
	return gw.SpriteSheet
}

func (gw *Window) AddPlayer(p player.Player) {
	gw.Players = append(gw.Players, &p)
}

func (gw *Window) GetPlayers() []*player.Player {
	return gw.Players
}

func (gw *Window) ResetPlayers() {
	gw.Players = make([]*player.Player, 0)
}

func (gw *Window) GetGrid() *grid.GameGrid {
	return gw.Grid
}

func (gw *Window) GetLawRating() int {
	return gw.LawRating
}

func (gw *Window) AdjustLawRating(a int) {
	gw.LawRating += a
}

func (gw *Window) Closed() bool {
	return gw.Window.Closed()
}

func (gw *Window) Update(screen screeniface.GameScreen) screeniface.GameScreen {
	newScreen := screen.Step(gw)
	gw.Window.Update()
	if newScreen != screen {
		fmt.Printf("New screen %T %+v\n", newScreen, newScreen)
		newScreen.Enter(gw)
		return newScreen
	}
	return screen
}

func NewWindow(ss pixel.Picture) *Window {
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
	grid := grid.MakeGameGrid(GridWidth, GridHeight)

	return &Window{
		Window:      win,
		SpriteSheet: ss,
		Grid:        grid,
		Players:     make([]*player.Player, 0),
	}
}
