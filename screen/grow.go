package screen

import (
	"fmt"
	"math/rand"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"

	"github.com/bobtfish/mayhem/character"
	"github.com/bobtfish/mayhem/logical"
	"github.com/bobtfish/mayhem/player"
)

const GROW_CHANCE = 10
const VANISH_CHANCE = 2

var growable map[string]bool

func init() {
	growable = map[string]bool{
		"Gooey Blob": false,
		"Fire":       true,
	}
}

func doesItGrow() bool {
	return rand.Intn(100) <= GROW_CHANCE
}

func doesItVanish() bool {
	return rand.Intn(100) <= VANISH_CHANCE
}

type GrowScreen struct {
	*WithBoard
	Consider logical.Vec
}

func (screen *GrowScreen) Enter(ss pixel.Picture, win *pixelgl.Window) {
	ClearScreen(ss, win)
}

func (screen *GrowScreen) Step(ss pixel.Picture, win *pixelgl.Window) GameScreen {
	batch := screen.WithBoard.DrawBoard(ss, win)
	batch.Draw(win)

	if screen.IterateGrowVanish() {
		// We just grew or vanished a tile, pause for a sec
		return &Pause{
			Grid:       screen.WithBoard.Grid,
			NextScreen: screen,
		}
	}

	return &Pause{
		Grid: screen.WithBoard.Grid,
		NextScreen: &TurnMenuScreen{
			Players: screen.WithBoard.Players,
			Grid:    screen.WithBoard.Grid,
		},
	}
}

func (screen *GrowScreen) IterateGrowVanish() bool {
	for screen.Consider.X < screen.WithBoard.Grid.MaxX() && screen.Consider.Y < screen.WithBoard.Grid.MaxY() {
		// Store current tile we're working on
		thisTileVec := screen.Consider

		// Bump tile counter as we may return early
		screen.Consider.X = screen.Consider.X + 1
		if screen.Consider.X == screen.WithBoard.Grid.MaxX() {
			screen.Consider.X = 0
			screen.Consider.Y = screen.Consider.Y + 1
		}

		// If the current tile contains a character
		character, ok := screen.WithBoard.Grid.GetGameObject(thisTileVec).(*character.Character)
		if ok {
			for name, replace := range growable {
				// If we're a special growable character (blob or fire)
				if name == character.Name {
					fmt.Printf("v(%d, %d) has %s\n", thisTileVec.X, thisTileVec.Y, character.Name)
					if doesItGrow() {
						adj := screen.WithBoard.Grid.AsRect().Adjacents(thisTileVec)
						rand.Shuffle(len(adj), func(i, j int) { adj[i], adj[j] = adj[j], adj[i] })
						fmt.Printf("v(%d, %d) is growing to v(%d, %d)\n", thisTileVec.X, thisTileVec.Y, adj[0].X, adj[0].Y)

						// Never grow to cover a player, if we try to do that just skip the grow
						currentObj := screen.WithBoard.Grid.GetGameObject(adj[0])
						_, isPlayer := currentObj.(*player.Player)
						if isPlayer {
							break
						}

						c := character.Clone()
						c.BoardPosition = adj[0]
						if replace { // Fire burns things - remove everything already stacked
							removedObject := screen.WithBoard.Grid.GetGameObjectStack(adj[0]).RemoveTopObject()
							for removedObject != nil {
								removedObject = screen.WithBoard.Grid.GetGameObjectStack(adj[0]).RemoveTopObject()
							}
						}
						screen.WithBoard.Grid.PlaceGameObject(adj[0], c)
						return true
					} else {
						if doesItVanish() {
							fmt.Printf("v(%d, %d) has vanished\n")
							screen.WithBoard.Grid.GetGameObjectStack(thisTileVec).RemoveTopObject()
							return true
						}
					}
					// Don't bother to check if we're another character type, we already matched
					break
				}
			}
		}
	}
	// We covered every tile already
	return false
}
