package spritesheet

import (
	"image"
)

type Spritesheet struct {
	WidthInTiles  int
	HeightInTiles int
	Tilesize      int
}

func (s *Spritesheet) Rect(index int) image.Rectangle {
	x := (index % s.WidthInTiles) * s.Tilesize
	y := (index / s.WidthInTiles) * s.Tilesize

	return image.Rect(
		x, y, x+s.Tilesize, y+s.Tilesize,
	)
}

func NewSpriteSheet(w, h, t int) *Spritesheet {
	return &Spritesheet{
		w, h, t,
	}
}
