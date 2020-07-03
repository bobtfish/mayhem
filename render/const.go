package render

import "github.com/bobtfish/mayhem/logical"

const CHAR_PIXELS = 64
const SPRITE_SIZE = 16

func MainScreenV() logical.Vec {
	return logical.V(0, CHAR_PIXELS)
}

func GameBoardV() logical.Vec {
	return MainScreenV().Add(logical.V(CHAR_PIXELS/2, CHAR_PIXELS/2))
}
