package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"

	"math/rand"

	"github.com/bobtfish/mayhem/logical"
	"github.com/bobtfish/mayhem/screen"
)

const GRID_X = 15
const GRID_Y = 10

func pickColour() pixel.RGBA {
	return pixel.RGB(rand.Float64(), rand.Float64(), rand.Float64())
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
		grid.PlaceGameObject(logical.V(x, y), ct.NewCharacter(k))
		x++
		if x == 15 {
			x = 0
			y++
		}
	}
}

func getPlayers() PlayerList {
	l := make([]Player, 2)
	l[0] = NewHumanPlayer("Player1", 0)
	l[1] = NewHumanPlayer("Player2", 1)
	return l
}

func placePlayers(players PlayerList, grid *GameGrid) {
	y := 9
	for i, player := range players {
		x := 0
		if (i % 2) != 0 {
			x = 14
		}
		grid.PlaceGameObject(logical.V(x, y), player)
		if (i % 2) != 0 {
			y--
		}
	}
}

func blowSomethingUp(grid *GameGrid) {
	x := rand.Intn(14)
	y := rand.Intn(9)
	fxA := []*Fx{FxWarp(), FxBlam(), FxFire(), FxBoom(), FxPop()}
	fxn := rand.Intn(len(fxA))
	grid.PlaceGameObject(logical.V(x, y), fxA[fxn])
}

func run() {
	data, err := base64.StdEncoding.DecodeString(sprite_sheet_base64)
	if err != nil {
		panic(err)
	}
	r := bytes.NewReader(data)

	//ct := LoadCharacterTemplates()
	//	grid := MakeGameGrid(logical.V(GRID_X, GRID_Y))

	title := "Mayhem!"

	gw := screen.NewGameWindow(r)

	//	players := getPlayers()

	//	placePlayers(players, grid)

	//	placeCharactersTest(grid, ct)

	QsecondTicks := 0
	frames := 0
	Qsecond := time.Tick(time.Second / 4)

	for !gw.Closed() {
		//		batch := grid.DrawBatch(&gw.SpriteDrawer)
		//		batch.Draw(gw.Window)

		gw.Update()

		frames++
		select {
		case <-Qsecond:
			//			grid.AnimationTick()
			QsecondTicks++
			if QsecondTicks == 4 {
				gw.Window.SetTitle(fmt.Sprintf("%s | FPS: %d", title, frames*4))
				frames = 0
				QsecondTicks = 0
				gw.MaybeNextScreen()
				//				blowSomethingUp(grid)
			}
		default:
		}
	}
}

func main() {
	pixelgl.Run(run)
}
