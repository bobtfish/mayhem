package screen

import (
	"fmt"
	"math/rand"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"

	"github.com/bobtfish/mayhem/logical"
	"github.com/bobtfish/mayhem/movable"
	"github.com/bobtfish/mayhem/render"
)

type RangedCombat struct {
	*WithBoard
	PlayerIdx       int
	Character       movable.Movable
	MovedCharacters map[movable.Movable]bool
}

func (screen *RangedCombat) Enter(ss pixel.Picture, win *pixelgl.Window) {
	fmt.Printf("In ranged combat state\n")
}

func (screen *RangedCombat) Step(ss pixel.Picture, win *pixelgl.Window) GameScreen {
	return &MoveFindCharacterScreen{
		WithBoard:       screen.WithBoard,
		PlayerIdx:       screen.PlayerIdx,
		MovedCharacters: screen.MovedCharacters,
	}
}

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
		return &RangedCombat{
			WithBoard:       screen.WithBoard,
			PlayerIdx:       screen.PlayerIdx,
			Charatcer:       screen.Character,
			MovedCharacters: screen.MovedCharacters,
		}
	}

	return screen
}

type DoAttack struct {
	*WithBoard
	AttackerV       logical.Vec
	DefenderV       logical.Vec
	PlayerIdx       int
	MovedCharacters map[movable.Movable]bool
}

func (screen *DoAttack) Enter(ss pixel.Picture, win *pixelgl.Window) {}

func (screen *DoAttack) Step(ss pixel.Picture, win *pixelgl.Window) GameScreen {
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
			died := screen.WithBoard.Grid.GetGameObjectStack(screen.DefenderV).RemoveTopObject()
			if KillIfPlayer(died, screen.WithBoard.Grid) {
				if WeHaveAWinner(screen.WithBoard.Players) {
					return &WinnerScreen{
						WithBoard: screen.WithBoard,
					}
				}
			}
		}

		doCharacterMove(screen.AttackerV, screen.DefenderV, screen.WithBoard.Grid)
	}
	return &RangedCombat{
		WithBoard:       screen.WithBoard,
		PlayerIdx:       screen.PlayerIdx,
		Charatcer:       screen.Character,
		MovedCharacters: screen.MovedCharacters,
	}
}
