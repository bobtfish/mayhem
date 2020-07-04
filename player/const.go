package player

import (
	"fmt"
	"github.com/bobtfish/mayhem/logical"
)

func GetStartPositions(number int) []logical.Vec {
	switch number {
	case 2:
		return []logical.Vec{
			logical.V(1, 5),
			logical.V(13, 5),
		}
	case 3:
		return []logical.Vec{
			logical.V(7, 8),
			logical.V(1, 1),
			logical.V(13, 1),
		}
	default:
		panic(fmt.Sprintf("Do not support %d players", number))
	}
}
