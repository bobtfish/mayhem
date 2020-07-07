package render

import (
	"image/color"

	"github.com/bobtfish/mayhem/logical"
	"github.com/faiface/pixel"
)

var textMap map[string]logical.Vec

func (sd SpriteDrawer) DrawTextColor(text string, win logical.Vec, color color.Color, target pixel.Target) {
	for _, c := range text {
		sd.GetSprite(textMap[string(c)]).DrawColorMask(target, sd.GetSpriteMatrix(win), color)
		win.X++
	}
}

func (sd SpriteDrawer) DrawText(text string, win logical.Vec, target pixel.Target) {
	sd.DrawTextColor(text, win, GetColor(255, 255, 255), target)
}

func init() {
	textMap = map[string]logical.Vec{
		" ":  logical.V(1, 40),
		"!":  logical.V(2, 40),
		"\"": logical.V(3, 40),
		"#":  logical.V(4, 40),
		"$":  logical.V(5, 40),
		"%":  logical.V(6, 40),
		"&":  logical.V(7, 40),
		"'":  logical.V(8, 40),
		"(":  logical.V(9, 40),
		")":  logical.V(10, 40),
		"*":  logical.V(11, 40),
		"+":  logical.V(12, 40),
		",":  logical.V(13, 40),
		"-":  logical.V(14, 40),
		".":  logical.V(15, 40),
		"/":  logical.V(16, 40),
		"0":  logical.V(17, 40),
		"1":  logical.V(18, 40),
		"2":  logical.V(19, 40),
		"3":  logical.V(0, 39),
		"4":  logical.V(1, 39),
		"5":  logical.V(2, 39),
		"6":  logical.V(3, 39),
		"7":  logical.V(4, 39),
		"8":  logical.V(5, 39),
		"9":  logical.V(6, 39),
		":":  logical.V(7, 39),
		";":  logical.V(8, 39),
		"<":  logical.V(9, 39),
		"=":  logical.V(10, 39),
		">":  logical.V(11, 39),
		"?":  logical.V(12, 39),
		"@":  logical.V(13, 39),
		"A":  logical.V(14, 39),
		"B":  logical.V(15, 39),
		"C":  logical.V(16, 39),
		"D":  logical.V(17, 39),
		"E":  logical.V(18, 39),
		"F":  logical.V(19, 39),
		"G":  logical.V(0, 38),
		"H":  logical.V(1, 38),
		"I":  logical.V(2, 38),
		"J":  logical.V(3, 38),
		"K":  logical.V(4, 38),
		"L":  logical.V(5, 38),
		"M":  logical.V(6, 38),
		"N":  logical.V(7, 38),
		"O":  logical.V(8, 38),
		"P":  logical.V(9, 38),
		"Q":  logical.V(10, 38),
		"R":  logical.V(11, 38),
		"S":  logical.V(12, 38),
		"T":  logical.V(13, 38),
		"U":  logical.V(14, 38),
		"V":  logical.V(15, 38),
		"W":  logical.V(16, 38),
		"X":  logical.V(17, 38),
		"Y":  logical.V(18, 38),
		"Z":  logical.V(19, 38),
		"[":  logical.V(0, 37),
		"\\": logical.V(1, 37),
		"]":  logical.V(2, 37),
		"^":  logical.V(3, 37),
		"_":  logical.V(4, 37),
		"`":  logical.V(5, 37),
		"a":  logical.V(6, 37),
		"b":  logical.V(7, 37),
		"c":  logical.V(8, 37),
		"d":  logical.V(9, 37),
		"e":  logical.V(10, 37),
		"f":  logical.V(11, 37),
		"g":  logical.V(12, 37),
		"h":  logical.V(13, 37),
		"i":  logical.V(14, 37),
		"j":  logical.V(15, 37),
		"k":  logical.V(16, 37),
		"l":  logical.V(17, 37),
		"m":  logical.V(18, 37),
		"n":  logical.V(19, 37),
		"o":  logical.V(0, 36),
		"p":  logical.V(1, 36),
		"q":  logical.V(2, 36),
		"r":  logical.V(3, 36),
		"s":  logical.V(4, 36),
		"t":  logical.V(5, 36),
		"u":  logical.V(6, 36),
		"v":  logical.V(7, 36),
		"w":  logical.V(8, 36),
		"x":  logical.V(9, 36),
		"y":  logical.V(10, 36),
		"z":  logical.V(11, 36),
		"{":  logical.V(12, 36),
		"|":  logical.V(13, 36),
		"}":  logical.V(14, 36),
		"~":  logical.V(15, 36),
	}
}
