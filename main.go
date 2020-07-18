package main

import (
	"bytes"
	"encoding/base64"
	"flag"
	"fmt"
	"io"
	"math/rand"
	"time"

	"github.com/faiface/pixel/pixelgl"

	"github.com/bobtfish/mayhem/character"
	"github.com/bobtfish/mayhem/logical"
	"github.com/bobtfish/mayhem/player"
	"github.com/bobtfish/mayhem/render"
	"github.com/bobtfish/mayhem/screen"
	"github.com/bobtfish/mayhem/spells"
)

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

	gw := screen.NewGameWindow(ss)

	if *quickPtr {
		gw.Screen = &screen.StartMainGame{
			Players: []*player.Player{
				&player.Player{
					Name:          "fred",
					HumanPlayer:   true,
					CharacterIcon: logical.V(0, 23),
					Color:         render.GetColor(255, 0, 0),
					ChosenSpell:   -1,
					Spells:        spells.ChooseSpells(),
					Alive:         true,
				},
				&player.Player{
					Name:          "bob",
					HumanPlayer:   true,
					CharacterIcon: logical.V(1, 23),
					Color:         render.GetColor(255, 0, 255),
					ChosenSpell:   -1,
					Spells:        spells.ChooseSpells(),
					Alive:         true,
				},
			},
		}
		gw.Screen.Enter(gw.SpriteSheet, gw.Window)
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

		gw.Update()

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
	rand.Seed(int64(time.Now().Nanosecond()))
	pixelgl.Run(run)
}
