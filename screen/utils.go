package screen

import (
	"github.com/faiface/pixel/pixelgl"

	"github.com/bobtfish/mayhem/logical"
	"github.com/bobtfish/mayhem/render"
)

func NextPlayerIdx(playerIdx int, players []*player.Player) int {
	playerIdx++
	if playerIdx == len(players) {
		return playerIdx
	}
	if !players[playerIdx].Alive {
		return NextPlayerIdx(playerIdx, players)
	}
	return playerIdx
}

var numKeyMap map[pixelgl.Button]int
var spellKeyMap map[pixelgl.Button]int
var directionKeyMap map[pixelgl.Button]logical.Vec

func captureNumKey(win *pixelgl.Window) int {
	return captureKey(win, numKeyMap)
}

func captureSpellKey(win *pixelgl.Window) int {
	return captureKey(win, spellKeyMap)
}

func captureKey(win *pixelgl.Window, keyMap map[pixelgl.Button]int) int {
	for button, r := range keyMap {
		if win.JustPressed(button) {
			return r
		}
	}
	return -1
}

func captureDirectionKey(win *pixelgl.Window) logical.Vec {
	for button, r := range directionKeyMap {
		if win.JustPressed(button) {
			return r
		}
	}
	return logical.ZeroVec()
}

func init() {
	numKeyMap = map[pixelgl.Button]int{
		pixelgl.Key0: 0,
		pixelgl.Key1: 1,
		pixelgl.Key2: 2,
		pixelgl.Key3: 3,
		pixelgl.Key4: 4,
		pixelgl.Key5: 5,
		pixelgl.Key6: 6,
		pixelgl.Key7: 7,
		pixelgl.Key8: 8,
		pixelgl.Key9: 9,
	}
	spellKeyMap = map[pixelgl.Button]int{
		pixelgl.KeyA: 0,
		pixelgl.KeyB: 1,
		pixelgl.KeyC: 2,
		pixelgl.KeyD: 3,
		pixelgl.KeyE: 4,
		pixelgl.KeyF: 5,
		pixelgl.KeyG: 6,
		pixelgl.KeyH: 7,
		pixelgl.KeyI: 8,
		pixelgl.KeyJ: 9,
		pixelgl.KeyK: 10,
		pixelgl.KeyL: 11,
		pixelgl.KeyM: 12,
		pixelgl.KeyN: 13,
	}
	directionKeyMap = map[pixelgl.Button]logical.Vec{
		pixelgl.KeyA: logical.V(-1, 0),
		pixelgl.KeyD: logical.V(1, 0),
		pixelgl.KeyQ: logical.V(-1, 1),
		pixelgl.KeyW: logical.V(0, 1),
		pixelgl.KeyE: logical.V(1, 1),
		pixelgl.KeyZ: logical.V(-1, -1),
		pixelgl.KeyX: logical.V(0, -1),
		pixelgl.KeyC: logical.V(1, -1),
	}
}

func intToChar(i int) string {
	return string('A' + i)
}

func drawMainBorder(win *pixelgl.Window, sd render.SpriteDrawer) {
	batch := sd.GetNewBatch()
	color := render.GetColor(0, 0, 255)
	// Bottom left
	sd.DrawSpriteColor(logical.V(6, 20), logical.V(0, 0), color, batch)
	// Bottom row
	for i := 1; i < 15; i++ {
		sd.DrawSpriteColor(logical.V(7, 20), logical.V(i, 0), color, batch)
	}
	// Bottom right
	sd.DrawSpriteColor(logical.V(8, 20), logical.V(15, 0), color, batch)
	// LHS and RHS
	for i := 1; i < 10; i++ {
		sd.DrawSpriteColor(logical.V(5, 20), logical.V(0, i), color, batch)
		sd.DrawSpriteColor(logical.V(9, 20), logical.V(15, i), color, batch)
	}
	// Top Left
	sd.DrawSpriteColor(logical.V(2, 20), logical.V(0, 10), color, batch)
	// Top row
	for i := 1; i < 15; i++ {
		sd.DrawSpriteColor(logical.V(3, 20), logical.V(i, 10), color, batch)
	}
	// Top right
	sd.DrawSpriteColor(logical.V(4, 20), logical.V(15, 10), color, batch)
	batch.Draw(win)
}
