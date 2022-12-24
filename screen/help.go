package screen

import (
	"github.com/faiface/pixel/pixelgl"

	"github.com/bobtfish/mayhem/logical"
	"github.com/bobtfish/mayhem/render"
	screeniface "github.com/bobtfish/mayhem/screen/iface"
)

type HelpScreenMenu struct{}

func (screen *HelpScreenMenu) Enter(ctx screeniface.GameCtx) {
	win := ctx.GetWindow()
	ss := ctx.GetSpriteSheet()
	ClearScreen(ss, win)
	td := TextDrawer(ss)
	td.DrawText("         Help screen", logical.V(0, 9), render.ColorWhite(), win)
	td.DrawText("1. Keys", logical.V(9, 7), render.ColorWhite(), win)
	td.DrawText("2. Spells", logical.V(9, 6), render.ColorWhite(), win)
	td.DrawText("3. Combat", logical.V(9, 5), render.ColorWhite(), win)
	td.DrawText("4. Undead", logical.V(9, 4), render.ColorWhite(), win)
	td.DrawText("5. Mounts", logical.V(9, 3), render.ColorWhite(), win)
	td.DrawText("6. Victory", logical.V(9, 2), render.ColorWhite(), win)

	textBottomColor("Press Keys 1-6 or 0 to return", render.ColorWhite(), ss, win)
}

func (screen *HelpScreenMenu) Step(ctx screeniface.GameCtx) screeniface.GameScreen {
	win := ctx.GetWindow()
	if win.JustPressed(pixelgl.Key0) {
		return &InitialScreen{}
	}
	if win.JustPressed(pixelgl.Key1) {
		return &HelpScreenKeys{}
	}
	if win.JustPressed(pixelgl.Key2) {
		return &HelpScreenSpells{}
	}
	if win.JustPressed(pixelgl.Key3) {
		return &HelpScreenCombat{}
	}
	if win.JustPressed(pixelgl.Key4) {
		return &HelpScreenUndead{}
	}
	if win.JustPressed(pixelgl.Key5) {
		return &HelpScreenMounts{}
	}
	if win.JustPressed(pixelgl.Key6) {
		return &HelpScreenVictory{}
	}
	return screen
}

type HelpScreenKeys struct{}

// nolint
func (screen *HelpScreenKeys) Enter(ctx screeniface.GameCtx) {
	win := ctx.GetWindow()
	ss := ctx.GetSpriteSheet()
	ClearScreen(ss, win)
	td := TextDrawer(ss)
	td.DrawText("              Keys", logical.V(0, 9), render.ColorWhite(), win)
	td.DrawText("AQWEDCXZ - Move in direction", logical.V(0, 7), render.ColorWhite(), win)
	td.DrawText("S - Select creature/wizard", logical.V(0, 6), render.ColorWhite(), win)
	td.DrawText("K - Cancel movement/attack", logical.V(0, 5), render.ColorWhite(), win)
	td.DrawText("I - Show information on", logical.V(0, 4), render.ColorWhite(), win)
	td.DrawText("    creature", logical.V(0, 3), render.ColorWhite(), win)
	td.DrawText("1-8 - Highlight creations of", logical.V(0, 2), render.ColorWhite(), win)
	td.DrawText("      player # 1-8", logical.V(0, 1), render.ColorWhite(), win)
	td.DrawText("0 - End turn", logical.V(0, 0), render.ColorWhite(), win)
	textBottomColor("   Press any key to continue", render.ColorWhite(), ss, win)
}

func (screen *HelpScreenKeys) Step(ctx screeniface.GameCtx) screeniface.GameScreen {
	win := ctx.GetWindow()
	if win.Typed() != "" {
		return &HelpScreenMenu{}
	}
	return screen
}

type HelpScreenSpells struct{}

func (screen *HelpScreenSpells) Enter(ctx screeniface.GameCtx) {
	win := ctx.GetWindow()
	ss := ctx.GetSpriteSheet()
	ClearScreen(ss, win)
	td := TextDrawer(ss)
	td.DrawText("           Spells", logical.V(0, 9), render.ColorWhite(), win)
	//                                        #
	td.DrawText("Select a spell then use", logical.V(0, 7), render.ColorWhite(), win)
	td.DrawText("direction keys to choose", logical.V(0, 6), render.ColorWhite(), win)
	td.DrawText("where to cast it.", logical.V(0, 5), render.ColorWhite(), win)
	td.DrawText("Press S to cast.", logical.V(0, 4), render.ColorWhite(), win)
	td.DrawText("Illusions always succeed but", logical.V(0, 3), render.ColorWhite(), win)
	td.DrawText("can be disbelieved by others.", logical.V(0, 2), render.ColorWhite(), win)
	td.DrawText("   ^=law *=chaos -=neutral", logical.V(0, 0), render.ColorWhite(), win)
	textBottomColor("   Press any key to continue", render.ColorWhite(), ss, win)
}

func (screen *HelpScreenSpells) Step(ctx screeniface.GameCtx) screeniface.GameScreen {
	win := ctx.GetWindow()
	if win.Typed() != "" {
		return &HelpScreenMenu{}
	}
	return screen
}

type HelpScreenCombat struct{}

func (screen *HelpScreenCombat) Enter(ctx screeniface.GameCtx) {
	win := ctx.GetWindow()
	ss := ctx.GetSpriteSheet()
	ClearScreen(ss, win)
	td := TextDrawer(ss)
	td.DrawText("           Combat", logical.V(0, 9), render.ColorWhite(), win)
	//                                        #
	td.DrawText("Move next to another creature", logical.V(0, 7), render.ColorWhite(), win)
	td.DrawText("to engage them in combat.", logical.V(0, 6), render.ColorWhite(), win)
	td.DrawText("Flying creatures can attack", logical.V(0, 5), render.ColorWhite(), win)
	td.DrawText("remotely without engagement.", logical.V(0, 4), render.ColorWhite(), win)
	td.DrawText("If adjacent next turn you may", logical.V(0, 3), render.ColorWhite(), win)
	td.DrawText("remain engaged or may be able", logical.V(0, 2), render.ColorWhite(), win)
	td.DrawText("to break away.", logical.V(0, 1), render.ColorWhite(), win)
	textBottomColor("   Press any key to continue", render.ColorWhite(), ss, win)
}

func (screen *HelpScreenCombat) Step(ctx screeniface.GameCtx) screeniface.GameScreen {
	win := ctx.GetWindow()
	if win.Typed() != "" {
		return &HelpScreenCombatRanged{}
	}
	return screen
}

type HelpScreenCombatRanged struct{}

// nolint
func (screen *HelpScreenCombatRanged) Enter(ctx screeniface.GameCtx) {
	win := ctx.GetWindow()
	ss := ctx.GetSpriteSheet()
	ClearScreen(ss, win)
	td := TextDrawer(ss)
	td.DrawText("        Ranged Combat", logical.V(0, 9), render.ColorWhite(), win)
	//                                        #
	td.DrawText("Some characters have ranged", logical.V(0, 7), render.ColorWhite(), win)
	td.DrawText("combat.", logical.V(0, 6), render.ColorWhite(), win)
	td.DrawText("This always happens after", logical.V(0, 5), render.ColorWhite(), win)
	td.DrawText("movement (K to skip movement)", logical.V(0, 4), render.ColorWhite(), win)
	td.DrawText("Target is selected with", logical.V(0, 3), render.ColorWhite(), win)
	td.DrawText("direction keys, press S to", logical.V(0, 2), render.ColorWhite(), win)
	td.DrawText("fire. Target must be in line", logical.V(0, 1), render.ColorWhite(), win)
	td.DrawText("of sight.", logical.V(0, 0), render.ColorWhite(), win)
	textBottomColor("   Press any key to continue", render.ColorWhite(), ss, win)
}

func (screen *HelpScreenCombatRanged) Step(ctx screeniface.GameCtx) screeniface.GameScreen {
	win := ctx.GetWindow()
	if win.Typed() != "" {
		return &HelpScreenMenu{}
	}
	return screen
}

type HelpScreenUndead struct{}

func (screen *HelpScreenUndead) Enter(ctx screeniface.GameCtx) {
	win := ctx.GetWindow()
	ss := ctx.GetSpriteSheet()
	ClearScreen(ss, win)
	td := TextDrawer(ss)
	td.DrawText("           Undead", logical.V(0, 9), render.ColorWhite(), win)
	//                                        #
	td.DrawText("Some characters are undead.", logical.V(0, 7), render.ColorWhite(), win)
	td.DrawText("They can only be attacked by", logical.V(0, 6), render.ColorWhite(), win)
	td.DrawText("other undead characters or", logical.V(0, 5), render.ColorWhite(), win)
	td.DrawText("magic weapons.", logical.V(0, 4), render.ColorWhite(), win)
	td.DrawText("The raise dead spell will", logical.V(0, 2), render.ColorWhite(), win)
	td.DrawText("turn a corpse into an undead", logical.V(0, 1), render.ColorWhite(), win)
	td.DrawText("creature.", logical.V(0, 0), render.ColorWhite(), win)
	textBottomColor("   Press any key to continue", render.ColorWhite(), ss, win)
}

func (screen *HelpScreenUndead) Step(ctx screeniface.GameCtx) screeniface.GameScreen {
	win := ctx.GetWindow()
	if win.Typed() != "" {
		return &HelpScreenMenu{}
	}
	return screen
}

type HelpScreenMounts struct{}

// nolint
func (screen *HelpScreenMounts) Enter(ctx screeniface.GameCtx) {
	win := ctx.GetWindow()
	ss := ctx.GetSpriteSheet()
	ClearScreen(ss, win)
	td := TextDrawer(ss)
	td.DrawText("           Mounts", logical.V(0, 9), render.ColorWhite(), win)
	//                                        #
	td.DrawText("Some characters can be ridden", logical.V(0, 7), render.ColorWhite(), win)
	td.DrawText("by wizards. Simply move your", logical.V(0, 6), render.ColorWhite(), win)
	td.DrawText("wizard onto the creature to", logical.V(0, 5), render.ColorWhite(), win)
	td.DrawText("mount it.", logical.V(0, 4), render.ColorWhite(), win)
	td.DrawText("This allows faster movement", logical.V(0, 3), render.ColorWhite(), win)
	td.DrawText("and your wizard cannot be", logical.V(0, 2), render.ColorWhite(), win)
	td.DrawText("killed unless their mount is", logical.V(0, 1), render.ColorWhite(), win)
	td.DrawText("killed first.", logical.V(0, 0), render.ColorWhite(), win)
	textBottomColor("   Press any key to continue", render.ColorWhite(), ss, win)
}

func (screen *HelpScreenMounts) Step(ctx screeniface.GameCtx) screeniface.GameScreen {
	win := ctx.GetWindow()
	if win.Typed() != "" {
		return &HelpScreenMenu{}
	}
	return screen
}

type HelpScreenVictory struct{}

func (screen *HelpScreenVictory) Enter(ctx screeniface.GameCtx) {
	win := ctx.GetWindow()
	ss := ctx.GetSpriteSheet()
	ClearScreen(ss, win)
	td := TextDrawer(ss)
	td.DrawText("          Victory", logical.V(0, 9), render.ColorWhite(), win)
	//                                        #
	td.DrawText("To win the game, simply kill", logical.V(0, 7), render.ColorWhite(), win)
	td.DrawText("all the other wizards.", logical.V(0, 6), render.ColorWhite(), win)
	td.DrawText("When a player is killed, all", logical.V(0, 4), render.ColorWhite(), win)
	td.DrawText("of their creations will also", logical.V(0, 3), render.ColorWhite(), win)
	td.DrawText("vanish.", logical.V(0, 2), render.ColorWhite(), win)
	textBottomColor("   Press any key to continue", render.ColorWhite(), ss, win)
}

func (screen *HelpScreenVictory) Step(ctx screeniface.GameCtx) screeniface.GameScreen {
	win := ctx.GetWindow()
	if win.Typed() != "" {
		return &HelpScreenMenu{}
	}
	return screen
}
