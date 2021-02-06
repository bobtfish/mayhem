package screen

import (
	"fmt"

	"github.com/faiface/pixel/pixelgl"

	"github.com/bobtfish/mayhem/fx"
	"github.com/bobtfish/mayhem/game"
	"github.com/bobtfish/mayhem/logical"
	"github.com/bobtfish/mayhem/render"
	screeniface "github.com/bobtfish/mayhem/screen/iface"
)

// Move state onto next player spell cast (if there are players left)
// or onto the movement phase if all spells have been cast
func NextSpellCastOrMove(playerIdx int, ctx screeniface.GameCtx, skipPause bool) screeniface.GameScreen {
	players := ctx.(*game.Window).GetPlayers()
	var nextScreen screeniface.GameScreen
	nextIdx := NextPlayerIdx(playerIdx, players)
	nextScreen = &DisplaySpellCastScreen{
		PlayerIdx: nextIdx,
	}

	if nextIdx == len(players) {
		// All players have cast their spells, movement comes next
		nextScreen = &MoveAnnounceScreen{}
	}
	return &Pause{
		Skip:       skipPause,
		NextScreen: nextScreen,
	}
}

// Display the spell name in the bottom bar until player presses a direction key or s

type DisplaySpellCastScreen struct {
	PlayerIdx int
}

func (screen *DisplaySpellCastScreen) Enter(ctx screeniface.GameCtx) {
	win := ctx.GetWindow()
	ss := ctx.GetSpriteSheet()
	players := ctx.(*game.Window).GetPlayers()
	ClearScreen(ss, win)
	thisPlayer := players[screen.PlayerIdx]
	if thisPlayer.ChosenSpell >= 0 {
		spell := thisPlayer.Spells[thisPlayer.ChosenSpell]
		batch := DrawBoard(ctx)
		textBottomMulti([]TextWithColor{
			TextWithColor{Text: fmt.Sprintf("%s ", thisPlayer.Name), Color: render.GetColor(255, 255, 0)},
			TextWithColor{Text: fmt.Sprintf("%s ", spell.GetName()), Color: render.GetColor(0, 242, 0)},
			TextWithColor{Text: fmt.Sprintf("%d", spell.GetCastRange()), Color: render.GetColor(244, 244, 244)},
		}, ss, batch)
		batch.Draw(win)
	}
}

func (screen *DisplaySpellCastScreen) Step(ctx screeniface.GameCtx) screeniface.GameScreen {
	win := ctx.GetWindow()
	players := ctx.(*game.Window).GetPlayers()
	thisPlayer := players[screen.PlayerIdx]
	batch := DrawBoard(ctx)
	batch.Draw(win)
	if (thisPlayer.ChosenSpell < 0) || win.JustPressed(pixelgl.Key0) {
		thisPlayer.ChosenSpell = -1 // Make sure to un-choose spell
		return NextSpellCastOrMove(screen.PlayerIdx, ctx, true)
	}
	if win.JustPressed(pixelgl.KeyS) || !captureDirectionKey(win).Equals(logical.ZeroVec()) {
		return &TargetSpellScreen{
			WithCursor: &WithCursor{
				CursorPosition: thisPlayer.BoardPosition,
			},
			PlayerIdx:      screen.PlayerIdx,
			CastsRemaining: thisPlayer.Spells[thisPlayer.ChosenSpell].CastQuantity(),
		}
	}
	return screen
}

// If range 0 then cast the spell straight away (on the wizard
// If range > 0 then move cursor around to find a target until S is pressed
type TargetSpellScreen struct {
	*WithCursor
	PlayerIdx      int
	CastsRemaining int
	MessageShown   bool
	CastBefore     bool
}

func (screen *TargetSpellScreen) Enter(ctx screeniface.GameCtx) {
	win := ctx.GetWindow()
	ss := ctx.GetSpriteSheet()
	textBottom("", ss, win) // clear bottom bar
}

func (screen *TargetSpellScreen) Step(ctx screeniface.GameCtx) screeniface.GameScreen {
	win := ctx.GetWindow()
	ss := ctx.GetSpriteSheet()
	grid := ctx.GetGrid()
	players := ctx.(*game.Window).GetPlayers()
	thisPlayer := players[screen.PlayerIdx]
	spell := thisPlayer.Spells[thisPlayer.ChosenSpell]
	batch := DrawBoard(ctx)

	if spell.GetCastRange() == 0 {
		target := screen.WithCursor.CursorPosition
		fmt.Printf("Cast spell %s (%d) on V(%d, %d)\n", spell.GetName(), spell.GetCastRange(), target.X, target.Y)
		return screen.AnimateAndCast(ctx)
	}
	if screen.WithCursor.MoveCursor(ctx) || !screen.MessageShown {
		screen.MessageShown = false
		screen.WithCursor.DrawCursor(ctx, batch)
	}
	if win.JustPressed(pixelgl.KeyS) {
		target := screen.WithCursor.CursorPosition
		if spell.GetCastRange() < target.Distance(thisPlayer.BoardPosition) {
			textBottomColor("Out of range", render.GetColor(0, 249, 249), ss, batch)
			fmt.Printf("Out of range! Spell cast range %d but distance to target is %d\n", spell.GetCastRange(), target.Distance(thisPlayer.BoardPosition))
			screen.MessageShown = true
		} else {
			if spell.NeedsLineOfSight() && !grid.HaveLineOfSight(thisPlayer.BoardPosition, target) {
				textBottom("No line of sight", ss, batch)
				screen.MessageShown = true
			} else {
				grid := ctx.GetGrid()
				if spell.CanCast(grid.GetGameObject(target)) {
					fmt.Printf("Cast spell %s (%d) on V(%d, %d)\n", spell.GetName(), spell.GetCastRange(), target.X, target.Y)
					return screen.AnimateAndCast(ctx)
				}
				fmt.Printf("Cannot cast '%s' on non-empty square\n", spell.GetName())
			}
		}
	}
	batch.Draw(win)
	if win.JustPressed(pixelgl.Key0) {
		return NextSpellCastOrMove(screen.PlayerIdx, ctx, true)
	}
	return screen
}

func (screen *TargetSpellScreen) AnimateAndCast(ctx screeniface.GameCtx) screeniface.GameScreen {
	players := ctx.(*game.Window).GetPlayers()
	grid := ctx.GetGrid()
	thisPlayer := players[screen.PlayerIdx]
	spell := thisPlayer.Spells[thisPlayer.ChosenSpell]
	target := screen.WithCursor.CursorPosition
	anim := spell.CastFx()
	if anim != nil {
		grid.PlaceGameObject(target, anim)
	}
	return &WaitForFx{
		NextScreen: &DoSpellCast{
			Target:         target,
			PlayerIdx:      screen.PlayerIdx,
			CastsRemaining: screen.CastsRemaining,
			CastBefore:     screen.CastBefore,
		},
		Fx: anim,
	}
}

// This screen does the actual mechanics of the animation
// and then casting the spell once animation is finished

type DoSpellCast struct {
	Target         logical.Vec
	PlayerIdx      int
	CastsRemaining int
	CastBefore     bool
}

func (screen *DoSpellCast) Enter(ctx screeniface.GameCtx) {
}

func (screen *DoSpellCast) Step(ctx screeniface.GameCtx) screeniface.GameScreen {
	win := ctx.GetWindow()
	ss := ctx.GetSpriteSheet()
	grid := ctx.GetGrid()
	players := ctx.(*game.Window).GetPlayers()
	castsRemaining := screen.CastsRemaining - 1
	batch := DrawBoard(ctx)
	// Fx for spell cast finished
	// Work out what happened :)
	fmt.Printf("About to call player CastSpell method\n")
	p := players[screen.PlayerIdx]

	fmt.Printf("IN Player spell cast\n")
	spell := p.Spells[p.ChosenSpell]
	fmt.Printf("Player spell %T cast on %T\n", spell, screen.Target)

	cleanupFunc := func() {
		if !spell.IsReuseable() {
			p.Spells = append(p.Spells[:p.ChosenSpell], p.Spells[p.ChosenSpell+1:]...)
		}
		p.ChosenSpell = -1
	}
	nextScreen := NextSpellCastOrMove(screen.PlayerIdx, ctx, false)

	takeOver := spell.TakeOverScreen(ctx, cleanupFunc, nextScreen, screen.PlayerIdx, screen.Target)
	if takeOver != nil {
		return takeOver
	}

	// Do the spell cast here
	var success bool
	var canCastMore bool
	var anim *fx.Fx
	if screen.CastBefore || spell.CastSucceeds(p.CastIllusion, ctx.GetLawRating()) {
		canCastMore = true
		success, anim = spell.DoCast(p.CastIllusion, screen.Target, grid, p)
	}

	fmt.Printf("Finished player CastSpell method\n")
	if !screen.CastBefore {
		if (canCastMore && castsRemaining > 0) || success {
			fmt.Printf("Spell Succeeds\n")
			textBottom("Spell Succeeds", ss, batch)
			ctx.AdjustLawRating(spell.GetLawRating())
		} else {
			fmt.Printf("Spell failed\n")
			textBottomColor("Spell Failed", render.GetColor(255, 0, 255), ss, batch)
		}
	}
	batch.Draw(win)
	if castsRemaining > 0 && canCastMore {
		nextScreen = &TargetSpellScreen{
			WithCursor: &WithCursor{
				CursorPosition: screen.Target,
			},
			PlayerIdx:      screen.PlayerIdx,
			CastsRemaining: castsRemaining,
			CastBefore:     true,
		}
	} else {
		cleanupFunc()
	}
	return &WaitForFx{
		NextScreen: nextScreen,
		Fx:         anim,
	}

}
