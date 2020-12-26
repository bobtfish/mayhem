package main

import (
	"bytes"
	"encoding/base64"
	"flag"
	"fmt"
	"io"
	"time"

	"github.com/faiface/pixel/pixelgl"

	"github.com/bobtfish/mayhem/character"
	"github.com/bobtfish/mayhem/game"
	"github.com/bobtfish/mayhem/logical"
	"github.com/bobtfish/mayhem/player"
	"github.com/bobtfish/mayhem/render"
	"github.com/bobtfish/mayhem/screen"
	screeniface "github.com/bobtfish/mayhem/screen/iface"

	_ "github.com/bobtfish/mayhem/otherspells"
)

//	_ "github.com/bobtfish/mayhem/spellswithscreen"
//)

func loadSpriteSheet() io.Reader {
	data, err := base64.StdEncoding.DecodeString(sprite_sheet_base64)
	if err != nil {
		panic(err)
	}
	return bytes.NewReader(data)
}

func run() {
	character.LoadCharacterTemplates()

	quickPtr := flag.Bool("quick", false, "Skip Intro questions")
	flag.Parse()

	title := "Mayhem!"

	spriteReader := loadSpriteSheet()
	ss := render.GetSpriteSheet(spriteReader)

	gw := game.NewWindow(ss)
	var s screeniface.GameScreen
	s = &screen.InitialScreen{}
	s.Enter(gw)

	if *quickPtr {
		fred := player.NewPlayer()
		fred.Name = "fred"
		fred.CharacterIcon = logical.V(0, 23)
		fred.Color = render.GetColor(255, 0, 0)
		bob := player.NewPlayer()
		bob.Name = "bob"
		bob.CharacterIcon = logical.V(1, 23)
		bob.Color = render.GetColor(255, 0, 255)
		s = &screen.StartMainGame{
			Players: []*player.Player{&fred, &bob},
		}
		s.Enter(gw)
	}

	//	players := getPlayers()

	//	placePlayers(players, grid)

	//placeCharactersTest(grid, ct)

	frames := 0
	Second := time.Tick(time.Second)

	for !gw.Closed() {
		//sd := render.NewSpriteDrawer(ss).WithOffset(render.GameBoardV())
		//batch := grid.DrawBatch(sd)
		//batch.Draw(gw.Window)

		s = gw.Update(s)

		frames++
		select {
		case <-Second:
			gw.Window.SetTitle(fmt.Sprintf("%s | FPS: %d", title, frames))
			frames = 0
			//blowSomethingUp(grid)
		default:
		}
	}
}

func main() {
	pixelgl.Run(run)
}
