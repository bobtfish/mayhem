package screen

import (
	"fmt"

	"github.com/bobtfish/mayhem/character"
	"github.com/bobtfish/mayhem/fx"
	"github.com/bobtfish/mayhem/game"
	"github.com/bobtfish/mayhem/logical"
	"github.com/bobtfish/mayhem/player"
	"github.com/bobtfish/mayhem/rand"
	screeniface "github.com/bobtfish/mayhem/screen/iface"
)

const GrowChance = 15
const VanishChance = 4

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
	return rand.Intn(100) <= GrowChance
}

func doesItVanish() bool {
	return rand.Intn(100) <= VanishChance
}

type GrowScreen struct {
	*WithBoard
	Consider logical.Vec
	Fx       *fx.Fx
	Grew     bool
}

func (screen *GrowScreen) Enter(ctx screeniface.GameCtx) {
	win := ctx.GetWindow()
	ss := ctx.GetSpriteSheet()
	ClearScreen(ss, win)
}

func (screen *GrowScreen) Step(ctx screeniface.GameCtx) screeniface.GameScreen {
	win := ctx.GetWindow()
	players := ctx.(*game.Window).GetPlayers()
	batch := screen.WithBoard.DrawBoard(ctx)
	batch.Draw(win)

	screen.IterateGrowVanish(ctx)

	firstAlivePlayerIdx := NextPlayerIdx(-1, players)
	fmt.Printf("First alive player index %d\n", firstAlivePlayerIdx)
	nextScreen := &Pause{
		Skip: !screen.Grew,
		NextScreen: &TurnMenuScreen{
			PlayerIdx: firstAlivePlayerIdx,
		},
	}

	return &WaitForFx{
		Fx:         screen.Fx,
		NextScreen: nextScreen,
	}
}

// FIXME - lots of puke worthy type casting in here, should not need this special casing really....
func (screen *GrowScreen) IterateGrowVanish(ctx screeniface.GameCtx) {
	grid := ctx.GetGrid()
	for screen.Consider.X < grid.MaxX() && screen.Consider.Y < grid.MaxY() {
		// If the current tile contains a character
		char, ok := grid.GetGameObject(screen.Consider).(*character.Character)
		if ok {
			for name, replace := range growable {
				// If we're a special growable character (blob or fire)
				if name == char.Name {
					fmt.Printf("v(%d, %d) has %s\n", screen.Consider.X, screen.Consider.Y, char.Name)
					if doesItGrow() {
						adjIdx := 0
						adjNew := false
						adj := grid.AsRect().Adjacents(screen.Consider)
						rand.Shuffle(len(adj), func(i, j int) { adj[i], adj[j] = adj[j], adj[i] })

						// Try to grow into an uncovered square
						for !adjNew && adjIdx < len(adj) {
							adjChar, isChar := grid.GetGameObject(adj[adjIdx]).(*character.Character)
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
						currentObj := grid.GetGameObject(adj[adjIdx])
						_, isPlayer := currentObj.(*player.Player)
						if isPlayer {
							break
						}

						c := char.Clone()
						c.BoardPosition = adj[adjIdx]
						// FIXME - growable objects should always replace each other, so blob doesn't cover fire, it removes it
						if replace { // Fire burns things - remove everything already stacked
							removedObject := grid.GetGameObjectStack(adj[adjIdx]).RemoveTopObject()
							for removedObject != nil {
								removedObject = grid.GetGameObjectStack(adj[adjIdx]).RemoveTopObject()
							}
						}
						grid.PlaceGameObject(adj[adjIdx], c)
						screen.Grew = true
					} else if doesItVanish() {
						grid.GetGameObjectStack(screen.Consider).RemoveTopObject()
						screen.Grew = true
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
						grid.GetGameObjectStack(screen.Consider).RemoveTopObject()
						grid.PlaceGameObject(screen.Consider, char.BelongsTo) // Put the wizard back down
						f := fx.Disbelieve()
						screen.Fx = f
						grid.PlaceGameObject(screen.Consider, f) // Also put a nice animation down
					}
				}
			}
		}

		// Bump tile counter
		screen.Consider.X++
		if screen.Consider.X == grid.MaxX() {
			screen.Consider.X = 0
			screen.Consider.Y++
		}
	}
}
