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
    return logical.V(CHAR_PIXELS / 2, CHAR_PIXELS)
}

func NewSpriteDrawer(ss pixel.Picture, windowOffset logical.Vec) *SpriteDrawer {
	return &SpriteDrawer{
		SpriteSheet:    ss,
		SSConverterMin: logical.NewVecConverter(logical.ZeroVec(), SSSpriteSizeVec()),
		SSConverterMax: logical.NewVecConverter(SSSpriteSizeVec(), SSSpriteSizeVec()),
		WinConverter:   logical.NewVecConverter(windowOffset, WinSpriteSizeVec()),
	}
}

func NewTextDrawer(ss pixel.Picture, windowOffset logical.Vec) *SpriteDrawer {
	return &SpriteDrawer{
		SpriteSheet:    ss,
		SSConverterMin: logical.NewVecConverter(logical.ZeroVec(), SSTextSizeVec()),
		SSConverterMax: logical.NewVecConverter(SSTextSizeVec(), SSTextSizeVec()),
		WinConverter:   logical.NewVecConverter(windowOffset, WinTextSizeVec()),
	}
}

type SpriteDrawer struct {
	SpriteSheet    pixel.Picture
	SSConverterMin logical.VecConverter
	SSConverterMax logical.VecConverter
	WinConverter   logical.VecConverter
}

func (sd *SpriteDrawer) GetSprite(v logical.Vec) *pixel.Sprite {
	return pixel.NewSprite(
		sd.SpriteSheet,
		pixel.Rect{
			Min: sd.SSConverterMin.ToPixelVec(v),
			Max: sd.SSConverterMax.ToPixelVec(v),
		},
	)
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
