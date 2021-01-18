package spells

import (
	"fmt"
	"image/color"

	"github.com/bobtfish/mayhem/fx"
	"github.com/bobtfish/mayhem/grid"
	"github.com/bobtfish/mayhem/logical"
	"github.com/bobtfish/mayhem/rand"
	"github.com/bobtfish/mayhem/render"
	screeniface "github.com/bobtfish/mayhem/screen/iface"
)

type ASpell struct {
	Name                string
	LawRating           int
	Reuseable           bool
	CastingChance       int
	CastRange           int
	NoCastFx            bool
	Tries               int
	NoLineOfSightNeeded bool
}

func (s ASpell) TakeOverScreen(ctx screeniface.GameCtx, cleanupFunc func(), nextScreen screeniface.GameScreen, playeerIdx int, target logical.Vec) screeniface.GameScreen {
	return nil
}
func (s ASpell) CastFx() *fx.Fx {
	if s.NoCastFx {
		return nil
	}
	return fx.SpellCast()
}
func (s ASpell) GetName() string {
	return s.Name
}
func (s ASpell) GetDescriptionArray(lawRating int) []string {
	desc := make([]string, 0)
	desc = append(desc, "")
	desc = append(desc, fmt.Sprintf("Casting chance=%d%%", s.GetCastingChance(lawRating)))
	desc = append(desc, "")
	desc = append(desc, fmt.Sprintf("Range=%d", s.GetCastRange()))

	return desc
}
func (s ASpell) GetLawRating() int {
	return s.LawRating
}
func (s ASpell) IsReuseable() bool {
	return s.Reuseable
}
func (s ASpell) CastQuantity() int {
	if s.Tries == 0 {
		return 1
	}
	return s.Tries
}
func (s ASpell) GetCastingChance(gameLawRating int) int {
	// Neutral spells or neutral games give base chance.
	// Spells never get harder, so a lawful game + chaos spell or chaotic game + lawful spell both give base chance
	if s.LawRating == 0 || gameLawRating == 0 || (s.LawRating > 0 && gameLawRating < 0) || (s.LawRating < 0 && gameLawRating > 0) {
		return s.CastingChance
	}
	absLaw := gameLawRating
	if absLaw < 0 {
		absLaw = -absLaw
	}
	cc := s.CastingChance + (absLaw/4)*10 // Every 4 points makes it 10% easier
	if cc > 100 {
		return 100
	}
	return cc
}

func (s ASpell) GetCastRange() int {
	return s.CastRange
}

func (s ASpell) CanCast(target grid.GameObject) bool {
	return true
}
func (s ASpell) CanCastAsIllusion() bool {
	return false
}
func (s ASpell) CastSucceeds(illusion bool, gameLawRating int) bool {
	if illusion && !s.CanCastAsIllusion() {
		panic(fmt.Sprintf("Spell %s (type %T) cannot be illusion, but was cast as one anyway", s.Name, s))
	}
	if illusion || rand.Intn(100) <= s.GetCastingChance(gameLawRating) {
		return true
	}
	return false
}
func (s ASpell) NeedsLineOfSight() bool {
	return !s.NoLineOfSightNeeded
}

func CastingChanceColor(chance int) color.Color {
	switch {
	case chance >= 100:
		return render.GetColor(255, 255, 255)
	case chance >= 80:
		return render.GetColor(255, 255, 0)
	case chance >= 60:
		return render.GetColor(0, 255, 255)
	case chance >= 40:
		return render.GetColor(0, 255, 0)
	default:
		return render.GetColor(255, 0, 255)
	}
}
