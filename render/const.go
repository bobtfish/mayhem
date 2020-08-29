package render

import "github.com/bobtfish/mayhem/logical"

const CharPixels = 64
const SpriteSize = 16
const InverseVideoOffset = 10

func MainScreenV() logical.Vec {
	return logical.V(0, CharPixels)
}

func GameBoardV() logical.Vec {
	return MainScreenV().Add(logical.V(CharPixels/2, CharPixels/2))
}
