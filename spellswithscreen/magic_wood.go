package spellswithscreen

import (
	"fmt"
	"sort"

	"github.com/bobtfish/mayhem/character"
	"github.com/bobtfish/mayhem/game"
	"github.com/bobtfish/mayhem/logical"
	"github.com/bobtfish/mayhem/render"
	screens "github.com/bobtfish/mayhem/screen"
	screeniface "github.com/bobtfish/mayhem/screen/iface"
	"github.com/bobtfish/mayhem/spells"
	spelliface "github.com/bobtfish/mayhem/spells/iface"
)

type MagicWoodSpellScreen struct {
	NextScreen  screeniface.GameScreen
	CleanupFunc func() // This is a closure that removes the spell from the player after casting, called when leaving
	PlayerIdx   int
	MaybeTiles  []logical.Vec
}

func (screen *MagicWoodSpellScreen) Enter(ctx screeniface.GameCtx) {
	screen.MaybeTiles = getMagicWoodTiles(ctx, screen.PlayerIdx)
}

func (screen *MagicWoodSpellScreen) Step(ctx screeniface.GameCtx) screeniface.GameScreen {
	if len(screen.MaybeTiles) == 0 {
		screen.CleanupFunc()
		return &screens.Pause{
			NextScreen: screen.NextScreen,
		}
	}

	batch := DrawBoard(ctx)
	//ss := ctx.GetSpriteSheet()
	grid := ctx.GetGrid()
	win := ctx.GetWindow()
	players := ctx.(*game.Window).GetPlayers()
	player := players[screen.PlayerIdx]

	placeAt := screen.MaybeTiles[0]
	fmt.Printf("Place at X%d, Y%d", placeAt.X, placeAt.Y)
	grid.PlaceGameObject(placeAt, &character.Character{
		Name:            "Magic Wood",
		Sprite:          magicWoodSpriteVec(),
		Color:           render.GetColor(255, 255, 255),
		Defence:         2,
		MagicResistance: 5, // FIXME

		// FIXME - ugh this is gross - would it be better done up a level?
		BelongsTo: player,
	})

	batch.Draw(win)

	return &screens.Pause{
		NextScreen: screen,
	}
}

func getMagicWoodTiles(ctx screeniface.GameCtx, playerIdx int) []logical.Vec {
	grid := ctx.GetGrid()
	players := ctx.(*game.Window).GetPlayers()
	player := players[playerIdx]
	fromPosition := player.BoardPosition

	// Find all tiles we could possibly put something on
	maybeTiles := make([]logical.Vec, 0)
	for x := -4; x <= 4; x++ {
		for y := 4; y >= -4; y-- {
			vecAdjust := logical.V(x, y)
			possibleVec := fromPosition.Add(vecAdjust)
			fmt.Printf("Consider possible Vec X%d Y%d\n", possibleVec.X, possibleVec.Y)
			if grid.AsRect().Contains(possibleVec) && !vecAdjust.Equals(logical.ZeroVec()) && grid.GetGameObject(possibleVec).IsEmpty() && adjacentsAreEmpty(ctx, possibleVec) {
				fmt.Printf("Add possible Vec X%d Y%d\n", possibleVec.X, possibleVec.Y)
				maybeTiles = append(maybeTiles, possibleVec)
			}
		}
	}
	d := logical.DistanceSortedVecs{
		From: fromPosition,
		List: maybeTiles,
	}
	sort.Reverse(d)
	return d.List
}

func adjacentsAreEmpty(ctx screeniface.GameCtx, location logical.Vec) bool {
	grid := ctx.GetGrid()
	adjS := grid.AsRect().Adjacents(location)
	adjEmpty := true
	for _, adj := range adjS {
		if !grid.GetGameObject(adj).IsEmpty() {
			adjEmpty = false
			break
		}
	}
	return adjEmpty
}

func magicWoodSpriteVec() logical.Vec {
	return logical.V(5, 4)
}

func init() {
	spelliface.CreateSpell(ScreenSpell{
		ASpell: spells.ASpell{
			Name:          "Magic Wood",
			CastingChance: 100, // FIXME
			CastRange:     6,   // FIXME
		},
		TakeOverFunc: func(ctx screeniface.GameCtx, cleanupFunc func(), nextScreen screeniface.GameScreen, playerIdx int, target logical.Vec) screeniface.GameScreen {
			return &MagicWoodSpellScreen{
				CleanupFunc: cleanupFunc,
				NextScreen:  nextScreen,
				PlayerIdx:   playerIdx,
			}
		},
	})
}
