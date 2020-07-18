package screen

import (
	"fmt"
	"math/rand"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"

	"github.com/bobtfish/mayhem/fx"
	"github.com/bobtfish/mayhem/logical"
	"github.com/bobtfish/mayhem/movable"
	"github.com/bobtfish/mayhem/render"
)

type EngagedAttack struct {
	*WithBoard
	PlayerIdx       int
	Character       movable.Movable
	MovedCharacters map[movable.Movable]bool
}

func (screen *EngagedAttack) Enter(ss pixel.Picture, win *pixelgl.Window) {
	fmt.Printf("In engaged attack state\n")
}

func (screen *EngagedAttack) Step(ss pixel.Picture, win *pixelgl.Window) GameScreen {
	batch := screen.WithBoard.DrawBoard(ss, win)
	render.NewTextDrawer(ss).DrawText("Engaged to enemy                        ", logical.ZeroVec(), batch)
	batch.Draw(win)

	direction := captureDirectionKey(win)
	if direction != logical.ZeroVec() {
		fmt.Printf("Try to move in v(%d, %d)\n", direction.X, direction.Y)
		currentLocation := screen.Character.GetBoardPosition()
		newLocation := currentLocation.Add(direction)

		attackScreen, _ := DoAttackMaybe(currentLocation, newLocation, screen.PlayerIdx, screen.WithBoard, screen.MovedCharacters)
		if attackScreen != nil {
			fmt.Printf("Can attack in that direction")
			return attackScreen
		}
	}

    if win.JustPressed(pixelgl.Key0) || win.JustPressed(pixelgl.KeyK) {
        return &MoveFindCharacterScreen{
            WithBoard:       screen.WithBoard,
            PlayerIdx:       screen.PlayerIdx,
            MovedCharacters: screen.MovedCharacters,
        }
    }

	return screen
}

type DoAttack struct {
	*WithBoard
	Fx              *fx.Fx
	AttackerV       logical.Vec
	DefenderV       logical.Vec
	PlayerIdx       int
	MovedCharacters map[movable.Movable]bool
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

		// Work out what happened. This is overly simple, but equivalent to what the original game does :)
		defender := screen.WithBoard.Grid.GetGameObject(screen.DefenderV)
		defenceRating := defender.(movable.Attackable).GetDefence() + rand.Intn(9)
		attacker := screen.WithBoard.Grid.GetGameObject(screen.AttackerV)
		attackRating := attacker.(movable.Attackerable).GetCombat() + rand.Intn(9)

		fmt.Printf("Attack rating %d defence rating %d\n", attackRating, defenceRating)
		if attackRating > defenceRating {
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
			WithBoard:       screen.WithBoard,
			PlayerIdx:       screen.PlayerIdx,
			MovedCharacters: screen.MovedCharacters,
		}
	}
	return screen
}
