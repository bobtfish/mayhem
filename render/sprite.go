package render

import (
	"image"
	_ "image/png"
	"io"

	"github.com/faiface/pixel"

	"github.com/bobtfish/mayhem/logical"
)

func loadPicture(io io.Reader) (pixel.Picture, error) {
	img, _, err := image.Decode(io)
	if err != nil {
		return nil, err
	}
	return pixel.PictureDataFromImage(img), nil
}

func NewSpriteDrawer(io io.Reader) SpriteDrawer {
	ss, err := loadPicture(io)
	if err != nil {
		panic(err)
	}
	return SpriteDrawer{SpriteSheet: ss}
}

type SpriteDrawer struct {
	SpriteSheet pixel.Picture
}

func (sd *SpriteDrawer) GetSprite(v logical.Vec) *pixel.Sprite {
	return pixel.NewSprite(sd.SpriteSheet, pixel.R(float64(v.X*SPRITE_SIZE), float64(v.Y*SPRITE_SIZE), float64(v.X*SPRITE_SIZE+SPRITE_SIZE), float64(v.Y*SPRITE_SIZE+SPRITE_SIZE)))
}

func (sd *SpriteDrawer) GetSpriteMatrix(win logical.Vec) pixel.Matrix {
	mat := pixel.IM
	v := pixel.V(float64(win.X*CHAR_PIXELS), float64(win.Y*CHAR_PIXELS))
	mat = mat.Moved(v)
	mat = mat.ScaledXY(v, pixel.V(CHAR_PIXELS/SPRITE_SIZE, CHAR_PIXELS/SPRITE_SIZE))
	return mat.Moved(pixel.V(CHAR_PIXELS/2-1, CHAR_PIXELS/2-1))
}

func (sd *SpriteDrawer) DrawSprite(ss, win logical.Vec, target pixel.Target) {
	sd.GetSprite(ss).Draw(target, sd.GetSpriteMatrix(win))
}

func (sd *SpriteDrawer) GetNewBatch() *pixel.Batch {
	batch := pixel.NewBatch(&pixel.TrianglesData{}, sd.SpriteSheet)
	batch.Clear()
	return batch
}
