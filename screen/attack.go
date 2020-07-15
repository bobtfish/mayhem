package screen

import (
	"fmt"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"

	"github.com/bobtfish/mayhem/fx"
	"github.com/bobtfish/mayhem/logical"
	"github.com/bobtfish/mayhem/movable"
)

type DoAttack struct {
	*WithBoard
	Fx        *fx.Fx
	AttackerV logical.Vec
	DefenderV logical.Vec
	PlayerIdx int
}

func (screen *DoAttack) Enter(ss pixel.Picture, win *pixelgl.Window) {
	ClearScreen(ss, win)
	fmt.Printf("Enter DoAttack screen, place fx\n")
	fx := fx.FxAttack()
	screen.Fx = fx
	screen.WithBoard.Grid.PlaceGameObject(screen.DefenderV, fx)
}

func (screen *DoAttack) Step(ss pixel.Picture, win *pixelgl.Window) GameScreen {
	batch := screen.WithBoard.DrawBoard(ss, win)
	batch.Draw(win)

	// Run animation till attack is finished
	if screen.Fx.RemoveMe() {
		// Work out what happened :)
		defender := screen.WithBoard.Grid.GetGameObject(screen.DefenderV)
		//        defenceRating := attacked.GetDefence()
		attackSucceeds := true // FIXME
		if attackSucceeds {
			// If the defender can be killed, kill them. Otherwise remove them
			ob, corpsable := defender.(movable.Corpseable)
			fmt.Printf("Defender is %T corpsable %v ob %T(%v)\n", defender, corpsable, ob, ob)
			makesCorpse := corpsable && ob.CanMakeCorpse()
			if makesCorpse {
				fmt.Printf("make corpse\n")
				ob.MakeCorpse()
			} else {
				fmt.Printf("remove defender as no corpse\n")
				screen.WithBoard.Grid.GetGameObjectStack(screen.DefenderV).RemoveTopObject()
			}
			doCharacterMove(screen.AttackerV, screen.DefenderV, screen.WithBoard.Grid)
		}
		return &MoveFindCharacterScreen{
			WithBoard: screen.WithBoard,
			PlayerIdx: screen.PlayerIdx,
		}
	}
	return screen
}