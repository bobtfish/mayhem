package screen

import (
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

type Pause struct {
	NextScreen GameScreen
	Started    time.Time
}

func (screen *Pause) Enter(ss pixel.Picture, win *pixelgl.Window) {
	screen.Started = time.Now()
}

func (screen *Pause) Step(ss pixel.Picture, win *pixelgl.Window) GameScreen {
	if screen.Started.Add(time.Second).Before(time.Now()) {
		return screen.NextScreen
	}
	return screen
}
