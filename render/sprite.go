package render

import (
	"image"
	_ "image/png"
	"os"

	"github.com/faiface/pixel"
)

func loadPicture(path string) (pixel.Picture, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}
	return pixel.PictureDataFromImage(img), nil
}

func NewSpriteDrawer() SpriteDrawer {
	ss, err := loadPicture("sprite_sheet.png")
	if err != nil {
		panic(err)
	}
	return SpriteDrawer{SpriteSheet: ss}
}

type SpriteDrawer struct {
	SpriteSheet pixel.Picture
}

func (sd *SpriteDrawer) GetSprite(ssLX, ssLY int) *pixel.Sprite {
	return pixel.NewSprite(sd.SpriteSheet, pixel.R(float64(ssLX*SPRITE_SIZE), float64(ssLY*SPRITE_SIZE), float64(ssLX*SPRITE_SIZE+SPRITE_SIZE), float64(ssLY*SPRITE_SIZE+SPRITE_SIZE)))
}

func (sd *SpriteDrawer) GetSpriteMatrix(winLX, winLY int) pixel.Matrix {
	mat := pixel.IM
	v := pixel.V(float64(winLX*CHAR_PIXELS), float64(winLY*CHAR_PIXELS))
	mat = mat.Moved(v)
	mat = mat.ScaledXY(v, pixel.V(CHAR_PIXELS/SPRITE_SIZE, CHAR_PIXELS/SPRITE_SIZE))
	return mat.Moved(pixel.V(CHAR_PIXELS/2-1, CHAR_PIXELS/2-1))
}

func (sd *SpriteDrawer) DrawSprite(ssLX, ssLY, winLX, winLY int, target pixel.Target) {
	sd.GetSprite(ssLX, ssLY).Draw(target, sd.GetSpriteMatrix(winLX, winLY))
}

func (sd *SpriteDrawer) GetNewBatch() *pixel.Batch {
	batch := pixel.NewBatch(&pixel.TrianglesData{}, sd.SpriteSheet)
	batch.Clear()
	return batch
}
