package entities

import (
	// "fmt"
	// "image"
	// "image/color"
	// "log"

	"github.com/hajimehoshi/ebiten/v2"
	// "github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Sprite struct {
	Img  *ebiten.Image
	X, Y, Dx, Dy float64 // Dx is change in x and Dy is change in y.
}