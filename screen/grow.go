package screen

import (
	"fmt"
	"math/rand"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"

	"github.com/bobtfish/mayhem/character"
	"github.com/bobtfish/mayhem/fx"
	"github.com/bobtfish/mayhem/logical"
	"github.com/bobtfish/mayhem/player"
)

const GROW_CHANCE = 15
const VANISH_CHANCE = 4

var growable map[string]bool
var explodeIfMounted map[string]bool

func init() {
	// FIXME - this should be encoded in the characters themselves in some way, not hard coded here
	growable = map[string]bool{
		"Gooey Blob": false,
		"Magic Fire": true,
	}
	explodeIfMounted = map[string]bool{
		"Magic Castle": true,
		"Dark Citadel": true,
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
	Fx       *fx.Fx
	Grew     bool
}

func (screen *GrowScreen) Enter(ss pixel.Picture, win *pixelgl.Window) {
	ClearScreen(ss, win)
}

func (screen *GrowScreen) Step(ss pixel.Picture, win *pixelgl.Window) GameScreen {
	batch := screen.WithBoard.DrawBoard(ss, win)
	batch.Draw(win)

	screen.IterateGrowVanish()

	firstAlivePlayerIdx := NextPlayerIdx(-1, screen.WithBoard.Players)
	fmt.Printf("First alive player index %d\n", firstAlivePlayerIdx)
	nextScreen := &Pause{
		Skip: !screen.Grew,
		Grid: screen.WithBoard.Grid,
		NextScreen: &TurnMenuScreen{
			Players:   screen.WithBoard.Players,
			Grid:      screen.WithBoard.Grid,
			PlayerIdx: firstAlivePlayerIdx,
		},
	}

	if screen.Fx != nil {
		return &WaitForFx{
			Grid:       screen.WithBoard.Grid,
			Fx:         screen.Fx,
			NextScreen: nextScreen,
		}
	}
	return nextScreen
}

// FIXME - lots of puke worthy type casting in here, should not need this special casing really....
func (screen *GrowScreen) IterateGrowVanish() {
	for screen.Consider.X < screen.WithBoard.Grid.MaxX() && screen.Consider.Y < screen.WithBoard.Grid.MaxY() {
		// If the current tile contains a character
		char, ok := screen.WithBoard.Grid.GetGameObject(screen.Consider).(*character.Character)
		if ok {
			for name, replace := range growable {
				// If we're a special growable character (blob or fire)
				if name == char.Name {
					fmt.Printf("v(%d, %d) has %s\n", screen.Consider.X, screen.Consider.Y, char.Name)
					if doesItGrow() {
						adjIdx := 0
						adjNew := false
						adj := screen.WithBoard.Grid.AsRect().Adjacents(screen.Consider)
						rand.Shuffle(len(adj), func(i, j int) { adj[i], adj[j] = adj[j], adj[i] })

						// Try to grow into an uncovered square
						for !adjNew && adjIdx < len(adj) {
							adjChar, isChar := screen.WithBoard.Grid.GetGameObject(adj[adjIdx]).(*character.Character)
							if isChar && adjChar.Name == char.Name {
								adjIdx++
							} else {
								adjNew = true
							}
						}
						// Cannot grow in any direction
						if adjIdx == len(adj) {
							break
						}
						fmt.Printf("v(%d, %d) is growing to v(%d, %d)\n", screen.Consider.X, screen.Consider.Y, adj[adjIdx].X, adj[adjIdx].Y)

						// Never grow to cover a player, if we try to do that just skip the grow
						// FIXME - should be able to grow onto players who didn't cast it
						currentObj := screen.WithBoard.Grid.GetGameObject(adj[adjIdx])
						_, isPlayer := currentObj.(*player.Player)
						if isPlayer {
							break
						}

						c := char.Clone()
						c.BoardPosition = adj[adjIdx]
						// FIXME - growable objects should always replace each other, so blob doesn't cover fire, it removes it
						if replace { // Fire burns things - remove everything already stacked
							removedObject := screen.WithBoard.Grid.GetGameObjectStack(adj[adjIdx]).RemoveTopObject()
							for removedObject != nil {
								removedObject = screen.WithBoard.Grid.GetGameObjectStack(adj[adjIdx]).RemoveTopObject()
							}
						}
						screen.WithBoard.Grid.PlaceGameObject(adj[adjIdx], c)
						screen.Grew = true
					} else {
						if doesItVanish() {
							screen.WithBoard.Grid.GetGameObjectStack(screen.Consider).RemoveTopObject()
							screen.Grew = true
						}
					}
					// Don't bother to check if we're another character type, we already matched
					break
				}
			}
			for name := range explodeIfMounted {
				// If we're a special explodable character (castle or citadel)
				if name == char.Name {
					if char.CarryingPlayer && rand.Intn(9)+1 <= 2 { // 20% chance
						screen.Grew = true
						screen.WithBoard.Grid.GetGameObjectStack(screen.Consider).RemoveTopObject()
						screen.WithBoard.Grid.PlaceGameObject(screen.Consider, char.BelongsTo) // Put the wizard back down
						f := fx.FxDisbelieve()
						screen.Fx = f
						screen.WithBoard.Grid.PlaceGameObject(screen.Consider, f) // Also put a nice animation down
					}
				}
			}
		}

		// Bump tile counter
		screen.Consider.X = screen.Consider.X + 1
		if screen.Consider.X == screen.WithBoard.Grid.MaxX() {
			screen.Consider.X = 0
			screen.Consider.Y = screen.Consider.Y + 1
		}
	}
}
