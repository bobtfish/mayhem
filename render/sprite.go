package render

import (
	"image"

	"image/color"
	_ "image/png" // For the side effects
	"io"

	"github.com/faiface/pixel"

	"github.com/bobtfish/mayhem/logical"
)

func GetSpriteSheet(io io.Reader) pixel.Picture {
	img, _, err := image.Decode(io)
	if err != nil {
		panic(err)
	}

	// Recreate an image of double the width of the sprite sheet.
	// Copy the original sprite sheet, and then to the right copy
	// all the same pixels but with the colours inverted.
	// This allows us to draw inverse video :)
	size := img.Bounds().Size()
	wImg := image.NewNRGBA(image.Rect(0, 0, size.X*2, size.Y))
	for x := 0; x < size.X; x++ {
		for y := 0; y < size.Y; y++ {
			r, g, b, a := img.At(x, y).RGBA()
			pixel := color.RGBA{uint8(r), uint8(g), uint8(b), uint8(a)}
			newPixel := color.RGBA{uint8(^r), uint8(^g), uint8(^b), uint8(a)}
			wImg.Set(x, y, pixel)
			wImg.Set(x+size.X, y, newPixel)
		}
	}
	return pixel.PictureDataFromImage(wImg)
}

func NewSpriteDrawer(ss pixel.Picture) SpriteDrawer {
	return SpriteDrawer{
		SpriteSheet:      ss,
		SpriteSheetSizeV: logical.V(SpriteSize, SpriteSize),
		WinSizeV:         logical.V(CharPixels, CharPixels),
	}
}

func NewTextDrawer(ss pixel.Picture) SpriteDrawer {
	return SpriteDrawer{
		SpriteSheet:      ss,
		SpriteSheetSizeV: logical.V(SpriteSize/2, SpriteSize),
		WinSizeV:         logical.V(CharPixels/2, CharPixels),
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
	mat = mat.ScaledXY(v, pixel.V(CharPixels/SpriteSize, CharPixels/SpriteSize))
	return mat.Moved(pixel.V(CharPixels/2-1, CharPixels/2-1))
}

func (sd SpriteDrawer) DrawSpriteColor(ssPos, winPos logical.Vec, mask color.Color, target pixel.Target) {
	sd.GetSprite(ssPos).DrawColorMask(target, sd.GetSpriteMatrix(winPos), mask)
}

func (sd SpriteDrawer) DrawSprite(ssPos, winPos logical.Vec, target pixel.Target) {
	sd.GetSprite(ssPos).Draw(target, sd.GetSpriteMatrix(winPos))
}

func (sd *SpriteDrawer) GetNewBatch() *pixel.Batch {
	batch := pixel.NewBatch(&pixel.TrianglesData{}, sd.SpriteSheet)
	batch.Clear()
	return batch
}

func GetColor(r, g, b int) color.Color {
	if r == 0 && g == 0 && b == 0 {
		return pixel.RGB(1, 1, 1)
	}
	return pixel.RGB(float64(r)/255, float64(g)/255, float64(b)/255)
}
