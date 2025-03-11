package main

import (
	"fmt"
	"image"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Sprite struct {
	Img  *ebiten.Image
	X, Y float64
}

type Enemy struct {
	*Sprite
	FollowsPlayer bool
}

type Potion struct {
	*Sprite
	AmtHeal uint
}

type Player struct {
	*Sprite
	Health uint
}

type Game struct {
	player      *Player
	enemies     []*Enemy
	potions     []*Potion
	tilemapJSON *TilemapJSON
	tilemapImg  *ebiten.Image
}

func (g *Game) Update() error {

	if ebiten.IsKeyPressed(ebiten.KeyD) || ebiten.IsKeyPressed(ebiten.KeyRight) {
		g.player.X += 2
	}

	if ebiten.IsKeyPressed(ebiten.KeyA) || ebiten.IsKeyPressed(ebiten.KeyLeft) {
		g.player.X -= 2
	}

	if ebiten.IsKeyPressed(ebiten.KeyW) || ebiten.IsKeyPressed(ebiten.KeyUp) {
		g.player.Y -= 2
	}

	if ebiten.IsKeyPressed(ebiten.KeyS) || ebiten.IsKeyPressed(ebiten.KeyDown) {
		g.player.Y += 2
	}

	for _, enemy := range g.enemies {
		if enemy.FollowsPlayer {
			if enemy.X < g.player.X {
				enemy.X += 1
			} else if enemy.X > g.player.X {
				enemy.X -= 1
			}

			if enemy.Y < g.player.Y {
				enemy.Y += 1
			} else if enemy.Y > g.player.Y {
				enemy.Y -= 1
			}
		}
	}

	for _, potion := range g.potions {
		if g.player.X > potion.X {
			g.player.Health += potion.AmtHeal
			fmt.Printf("Picked up potion! Health: %d\n", g.player.Health)
		}
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{120, 180, 255, 255})

	op := &ebiten.DrawImageOptions{}

	// loop over the layers
	tilesetColumns := g.tilemapImg.Bounds().Dx() / 16 // Number of tiles per row in tileset
	for _, layer := range g.tilemapJSON.Layers {
		for index, id := range layer.Data {
			if id == 0 {
				continue
			}
			x := (index % layer.Width) * 16 // tile position x
			y := (index / layer.Width) * 16 // tile position y

			srcX := (id - 1) % tilesetColumns * 16
			srcY := (id - 1) / tilesetColumns * 16

			op.GeoM.Translate(float64(x), float64(y))

			screen.DrawImage(
				g.tilemapImg.SubImage(image.Rect(srcX, srcY, srcX+600, srcY+600)).(*ebiten.Image),
				op,
			)

			op.GeoM.Reset()
		}
	}

	// Scale factors (reduce size)
	scaleX := 0.3 // Shrinks width to 50%
	scaleY := 0.3 // Shrinks height to 50%

	// Apply scaling
	op.GeoM.Scale(scaleX, scaleY)

	// Track positions
	op.GeoM.Translate(g.player.X, g.player.Y)

	// Draw our player
	screen.DrawImage(
		g.player.Img.SubImage(
			image.Rect(0, 0, 150, 150),
		).(*ebiten.Image),
		op,
	)

	op.GeoM.Reset()

	// Draw enemies (goblins)
	for _, enemy := range g.enemies {

		// Scale factors (reduce size)
		scaleX := 0.3 // Shrinks width to 50%
		scaleY := 0.3 // Shrinks height to 50%

		// Apply scaling
		op.GeoM.Scale(scaleX, scaleY)

		op.GeoM.Translate(enemy.X, enemy.Y)

		screen.DrawImage(
			enemy.Img.SubImage(
				image.Rect(0, 0, 150, 150),
			).(*ebiten.Image),
			op,
		)

		op.GeoM.Reset()

	}

	// Draw potions
	for _, enemy := range g.potions {

		// Scale factors (reduce size)
		scaleX := 0.3 // Shrinks width to 50%
		scaleY := 0.3 // Shrinks height to 50%

		// Apply scaling
		op.GeoM.Scale(scaleX, scaleY)

		op.GeoM.Translate(enemy.X, enemy.Y)

		screen.DrawImage(
			enemy.Img.SubImage(
				image.Rect(0, 0, 150, 150),
			).(*ebiten.Image),
			op,
		)

		op.GeoM.Reset()

	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 320, 240
}

func main() {
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	playerImg, _, err := ebitenutil.NewImageFromFile("assets/images/warrior-main.png")
	if err != nil {
		// handle error
		log.Fatal(err)
	}

	goblinFireImg, _, err := ebitenutil.NewImageFromFile("assets/images/goblin_fire.png")
	if err != nil {
		// handle error
		log.Fatal(err)
	}

	potionImg, _, err := ebitenutil.NewImageFromFile("assets/images/meat.png")
	if err != nil {
		// handle error
		log.Fatal(err)
	}

	tilemapImg, _, err := ebitenutil.NewImageFromFile("assets/images/Tilemap_Flat.png")
	if err != nil {
		// handle error
		log.Fatal(err)
	}

	// Load tile map
	tilemapJSON, err := NewTilemapJSON("assets/maps/spawn.tmj")
	if err != nil {
		log.Fatal(err)
	}

	game := Game{
		player: &Player{
			Sprite: &Sprite{
				Img: playerImg,
				X:   50,
				Y:   50,
			},
			Health: 5,
		},

		enemies: []*Enemy{
			{
				&Sprite{
					Img: goblinFireImg,
					X:   150,
					Y:   150,
				},
				true,
			},
			{
				&Sprite{
					Img: goblinFireImg,
					X:   150,
					Y:   100,
				},
				false,
			},
		},
		potions: []*Potion{
			{
				&Sprite{
					Img: potionImg,
					X:   120,
					Y:   120,
				},
				5,
			},
		},
		tilemapJSON: tilemapJSON,
		tilemapImg:  tilemapImg,
	}

	if err := ebiten.RunGame(&game); err != nil {
		log.Fatal(err)
	}
}
