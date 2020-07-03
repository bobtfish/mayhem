package render

import (
	"image"
	_ "image/png"
	"io"

	"github.com/faiface/pixel"

	"github.com/bobtfish/mayhem/logical"
)

func GetSpriteSheet(io io.Reader) pixel.Picture {
	img, _, err := image.Decode(io)
	if err != nil {
		panic(err)
	}
	return pixel.PictureDataFromImage(img)
}

func NewSpriteDrawer(ss pixel.Picture) SpriteDrawer {
	return SpriteDrawer{
		SpriteSheet:      ss,
		SpriteSheetSizeV: logical.V(SPRITE_SIZE, SPRITE_SIZE),
		WinSizeV:         logical.V(CHAR_PIXELS, CHAR_PIXELS),
	}
}

func NewTextDrawer(ss pixel.Picture) SpriteDrawer {
	return SpriteDrawer{
		SpriteSheet:      ss,
		SpriteSheetSizeV: logical.V(SPRITE_SIZE/2, SPRITE_SIZE),
		WinSizeV:         logical.V(CHAR_PIXELS/2, CHAR_PIXELS),
	}
}

type SpriteDrawer struct {
	SpriteSheet      pixel.Picture
	SpriteSheetSizeV logical.Vec
	WinSizeV         logical.Vec
	WinOffsetV       logical.Vec
}

func (sd SpriteDrawer) WithOffset(v logical.Vec) SpriteDrawer {
	sd.WinOffsetV = v
	return sd
}

func (sd SpriteDrawer) GetSprite(v logical.Vec) *pixel.Sprite {
	return pixel.NewSprite(sd.SpriteSheet, v.ToPixelRect(sd.SpriteSheetSizeV))
}

func (sd SpriteDrawer) GetSpriteMatrix(winPos logical.Vec) pixel.Matrix {
	mat := pixel.IM
	v := winPos.Multiply(sd.WinSizeV).Add(sd.WinOffsetV).ToPixelVec()
	mat = mat.Moved(v)
	mat = mat.ScaledXY(v, pixel.V(CHAR_PIXELS/SPRITE_SIZE, CHAR_PIXELS/SPRITE_SIZE))
	return mat.Moved(pixel.V(CHAR_PIXELS/2-1, CHAR_PIXELS/2-1))
}

func (sd SpriteDrawer) DrawSprite(ssPos, winPos logical.Vec, target pixel.Target) {
	sd.GetSprite(ssPos).Draw(target, sd.GetSpriteMatrix(winPos))
}

func (sd *SpriteDrawer) GetNewBatch() *pixel.Batch {
	batch := pixel.NewBatch(&pixel.TrianglesData{}, sd.SpriteSheet)
	batch.Clear()
	return batch
}
