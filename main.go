package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"time"

	"github.com/faiface/pixel/pixelgl"

	"github.com/bobtfish/mayhem/render"
	"github.com/bobtfish/mayhem/screen"
)

const GRID_X = 15
const GRID_Y = 10

func loadSpriteSheet() io.Reader {
	data, err := base64.StdEncoding.DecodeString(sprite_sheet_base64)
	if err != nil {
		panic(err)
	}
	return bytes.NewReader(data)
}

func run() {
	//ct := LoadCharacterTemplates()
	//grid := MakeGameGrid(logical.V(GRID_X, GRID_Y))

	title := "Mayhem!"

	spriteReader := loadSpriteSheet()
	ss := render.GetSpriteSheet(spriteReader)

	gw := screen.NewGameWindow(ss)

	//	players := getPlayers()

	//	placePlayers(players, grid)

	//placeCharactersTest(grid, ct)

	QsecondTicks := 0
	frames := 0
	Qsecond := time.Tick(time.Second / 4)

	for !gw.Closed() {
		//sd := render.NewSpriteDrawer(ss).WithOffset(render.GameBoardV())
		//batch := grid.DrawBatch(sd)
		//batch.Draw(gw.Window)

		gw.Update()

		frames++
		select {
		case <-Qsecond:
			//grid.AnimationTick()
			QsecondTicks++
			if QsecondTicks == 4 {
				gw.Window.SetTitle(fmt.Sprintf("%s | FPS: %d", title, frames*4))
				frames = 0
				QsecondTicks = 0
				gw.MaybeNextScreen()
				//blowSomethingUp(grid)
			}
		default:
		}
	}
}

func main() {
	pixelgl.Run(run)
}
