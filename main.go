package main

import (
	"image"
	"log"
	"rpg-game-go/entities"

	"github.com/hajimehoshi/ebiten/v2"
)

func checkCollisonHorizontal(sprite *entities.Sprite, colliders []image.Rectangle) {
	for _, collider := range colliders {
		if collider.Overlaps(
			image.Rect(int(sprite.X), int(sprite.Y), int(sprite.X)+40, int(sprite.Y)+40),
		) {

			if sprite.Dx > 0.0 {
				sprite.X = float64(collider.Min.X) - 40
			} else if sprite.Dx < 0.0 {
				sprite.X = float64(collider.Max.X)
			}
		}
	}
}

func checkCollisonVertical(sprite *entities.Sprite, colliders []image.Rectangle) {
	for _, collider := range colliders {
		if collider.Overlaps(
			image.Rect(int(sprite.X), int(sprite.Y), int(sprite.X)+40, int(sprite.Y)+40),
		) {

			if sprite.Dy > 0.0 {
				sprite.Y = float64(collider.Min.Y) - 40
			} else if sprite.Dy < 0.0 {
				sprite.Y = float64(collider.Max.Y)
			}
		}
	}
}

func main() {
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	game :=NewGame()

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
