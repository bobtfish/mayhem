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
	case 4:
		return []logical.Vec{
			logical.V(1, 8),
			logical.V(13, 8),
			logical.V(1, 1),
			logical.V(13, 1),
		}
	case 5:
		return []logical.Vec{
			logical.V(7, 14),
			logical.V(0, 6),
			logical.V(14, 6),
			logical.V(3, 0),
			logical.V(11, 0),
		}
	case 6:
		return []logical.Vec{
			logical.V(7, 9),
			logical.V(0, 8),
			logical.V(14, 8),
			logical.V(0, 1),
			logical.V(7, 0),
			logical.V(14, 1),
		}
	case 7:
		return []logical.Vec{
			logical.V(7, 9),
			logical.V(1, 8),
			logical.V(13, 8),
			logical.V(0, 3),
			logical.V(14, 3),
			logical.V(4, 0),
			logical.V(10, 0),
		}
	case 8:
		return []logical.Vec{
			logical.V(0, 9),
			logical.V(7, 9),
			logical.V(14, 9),
			logical.V(0, 5),
			logical.V(14, 5),
			logical.V(0, 0),
			logical.V(7, 0),
			logical.V(14, 0),
		}
	default:
		panic(fmt.Sprintf("Do not support %d players", number))
	}
}
