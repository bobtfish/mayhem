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

func SSSpriteSizeVec() logical.Vec {
	return logical.V(SPRITE_SIZE, SPRITE_SIZE)
}

func SSTextSizeVec() logical.Vec {
	return logical.V(SPRITE_SIZE/2, SPRITE_SIZE)
}

func WinSpriteSizeVec() logical.Vec {
	return logical.V(CHAR_PIXELS, CHAR_PIXELS)
}

func WinTextSizeVec() logical.Vec {
	return logical.V(CHAR_PIXELS/2, CHAR_PIXELS)
}

func NewSpriteDrawer(ss pixel.Picture, windowOffset logical.Vec) *SpriteDrawer {
	return &SpriteDrawer{
		SpriteSheet:      ss,
		SpriteSheetSizeV: SSSpriteSizeVec(),
		WinConverter:     logical.NewVecConverter(windowOffset, WinSpriteSizeVec()),
	}
}

func NewTextDrawer(ss pixel.Picture, windowOffset logical.Vec) *SpriteDrawer {
	return &SpriteDrawer{
		SpriteSheet:      ss,
		SpriteSheetSizeV: SSTextSizeVec(),
		WinConverter:     logical.NewVecConverter(windowOffset, WinTextSizeVec()),
	}
}

type SpriteDrawer struct {
	SpriteSheet      pixel.Picture
	SpriteSheetSizeV logical.Vec
	WinConverter     logical.VecConverter
}

func (sd *SpriteDrawer) GetPixelRect(v logical.Vec) pixel.Rect {
	return pixel.Rect{
		Min: v.Multiply(sd.SpriteSheetSizeV).ToPixelVec(),
		Max: v.Multiply(sd.SpriteSheetSizeV).Add(sd.SpriteSheetSizeV).ToPixelVec(),
	}
}

func (sd *SpriteDrawer) GetSprite(v logical.Vec) *pixel.Sprite {
	return pixel.NewSprite(sd.SpriteSheet, sd.GetPixelRect(v))
}

func (sd *SpriteDrawer) GetSpriteMatrix(win logical.Vec) pixel.Matrix {
	mat := pixel.IM
	v := sd.WinConverter.ToPixelVec(win)
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
