package spellswithscreen

import (
	"fmt"
	"sort"

	"github.com/bobtfish/mayhem/rand"

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
	LaidTiles   int
}

func (screen *MagicWoodSpellScreen) Enter(ctx screeniface.GameCtx) {
	screen.MaybeTiles = getMagicWoodTiles(ctx, screen.PlayerIdx)
}

func (screen *MagicWoodSpellScreen) Step(ctx screeniface.GameCtx) screeniface.GameScreen {
	if len(screen.MaybeTiles) == 0 || screen.LaidTiles == 8 {
		screen.CleanupFunc()
		return &screens.Pause{
			NextScreen: screen.NextScreen,
		}
	}

	batch := DrawBoard(ctx)
	grid := ctx.GetGrid()
	win := ctx.GetWindow()
	players := ctx.(*game.Window).GetPlayers()
	player := players[screen.PlayerIdx]

	placeAt := screen.MaybeTiles[0]
	fmt.Printf("Place at X%d, Y%d", placeAt.X, placeAt.Y)
	grid.PlaceGameObject(placeAt, &character.Character{
		Name:            "Magic Wood",
		Sprite:          magicWoodSpriteVec(),
		Color:           render.GetColor(0, 252, 0),
		Defence:         5,
		MagicResistance: 0,
		Mountable:       true,
		BelongsTo:       player,
	})
	screen.LaidTiles++

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
	for x := -8; x <= 8; x++ {
		for y := 8; y >= -8; y-- {
			vecAdjust := logical.V(x, y)
			possibleVec := fromPosition.Add(vecAdjust)
			//fmt.Printf("Consider possible Vec X%d Y%d\n", possibleVec.X, possibleVec.Y)
			if !grid.AsRect().Contains(possibleVec) {
				continue
			}
			if vecAdjust.Equals(logical.ZeroVec()) {
				continue
			}
			if possibleVec.Distance(fromPosition) > 8 {
				continue
			}
			if !grid.GetGameObject(possibleVec).IsEmpty() {
				continue
			}
			// FIXME - adjacent to the player is OK
			if adjacentsAreNotMagicWood(ctx, possibleVec) {
				//fmt.Printf("Add possible Vec X%d Y%d\n", possibleVec.X, possibleVec.Y)
				maybeTiles = append(maybeTiles, possibleVec)
			}
		}
	}

	// Sort them, nearest to the wizard first
	d := logical.DistanceSortedVecs{
		From: fromPosition,
		List: maybeTiles,
	}
	sort.Sort(d)

	// Produce a list of the equsl closest possible tiles
	minDistance := d.List[0].Distance(fromPosition)
	actualList := make([]logical.Vec, 0)
	if len(d.List) == 0 {
		return d.List
	}
	for i := 0; i < len(d.List); i++ {
		if d.List[i].Distance(fromPosition) > minDistance {
			break
		}
		actualList = append(actualList, d.List[i])
	}

	rand.Shuffle(len(actualList), func(i, j int) { actualList[i], actualList[j] = actualList[j], actualList[i] })

	return actualList
}

func adjacentsAreNotMagicWood(ctx screeniface.GameCtx, location logical.Vec) bool {
	grid := ctx.GetGrid()
	adjS := grid.AsRect().Adjacents(location)
	adjEmpty := true
	for _, adj := range adjS {
		ob := grid.GetGameObject(adj)
		char, isCharacter := ob.(*character.Character)
		if isCharacter && char.Name == "Magic Wood" {
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
			LawRating:     1,
			CastingChance: 80,
			CastRange:     0, // Lol this is a hack
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
