package main

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"rpg-game-go/animations"
	"rpg-game-go/components"
	"rpg-game-go/constants"
	"rpg-game-go/entities"
	"rpg-game-go/spritesheet"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Game struct {
	player            *entities.Player
	playerSpriteSheet *spritesheet.Spritesheet
	animationFrame    int
	enemies           []*entities.Enemy
	potions           []*entities.Potion
	tilemapJSON       *TilemapJSON
	tilesets          []Tileset
	tilemapImg        *ebiten.Image
	cam               *Camera
	colliders         []image.Rectangle
}

func NewGame() *Game {
	playerImg, _, err := ebitenutil.NewImageFromFile("assets/images/warrior-main-2.png")
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

	// Generate tilesets
	tilesets, err := tilemapJSON.GenTilesets()
	if err != nil {
		log.Fatal(err)
	}

	playerSpriteSheet := spritesheet.NewSpriteSheet(6, 8, 192)

	game := Game{
		player: &entities.Player{
			Sprite: &entities.Sprite{
				Img: playerImg,
				X:   50,
				Y:   50,
			},
			Health: 5,
			Animations: map[entities.PlayerState]*animations.Animation{
				entities.Right: animations.NewAnimation(6, 11, 1, 8.0),
				entities.Left:  animations.NewAnimation(48, 53, 1, 8.0),
				entities.Down:  animations.NewAnimation(26, 30, 3, 8.0),
				entities.Up:    animations.NewAnimation(38, 42, 3, 8.0),
			},
			CombatComp: components.NewBasicCombat(3, 1),
		},
		playerSpriteSheet: playerSpriteSheet,
		enemies: []*entities.Enemy{
			{
				Sprite: &entities.Sprite{
					Img: goblinFireImg,
					X:   150,
					Y:   150,
				},
				FollowsPlayer: true,
				CombatComp:    components.NewBasicCombat(3, 1),
			},
			{
				Sprite: &entities.Sprite{
					Img: goblinFireImg,
					X:   150,
					Y:   100,
				},
				FollowsPlayer: false,
				CombatComp:    components.NewBasicCombat(3, 1),
			},
		},
		potions: []*entities.Potion{
			{
				Sprite: &entities.Sprite{
					Img: potionImg,
					X:   120,
					Y:   120,
				},
				AmtHeal: 5,
			},
		},
		tilemapJSON: tilemapJSON,
		tilemapImg:  tilemapImg,
		tilesets:    tilesets,
		cam:         NewCamera(50, 50),
		colliders: []image.Rectangle{
			image.Rect(100, 100, 116, 116),
		},
	}

	return &game
}

func (g *Game) Update() error {
	// set velocity to 0 initially to make it stop going in one direction on key press.
	g.player.Dx = 0
	g.player.Dy = 0

	if ebiten.IsKeyPressed(ebiten.KeyD) || ebiten.IsKeyPressed(ebiten.KeyRight) {
		g.player.Dx += 2
	}

	if ebiten.IsKeyPressed(ebiten.KeyA) || ebiten.IsKeyPressed(ebiten.KeyLeft) {
		g.player.Dx -= 2
	}

	if ebiten.IsKeyPressed(ebiten.KeyW) || ebiten.IsKeyPressed(ebiten.KeyUp) {
		g.player.Dy -= 2
	}

	if ebiten.IsKeyPressed(ebiten.KeyS) || ebiten.IsKeyPressed(ebiten.KeyDown) {
		g.player.Dy += 2
	}

	g.player.X += g.player.Dx

	checkCollisonHorizontal(g.player.Sprite, g.colliders)

	g.player.Y += g.player.Dy

	checkCollisonVertical(g.player.Sprite, g.colliders)

	activeAnimation := g.player.ActiveAnimation(
		int(g.player.Dx),
		int(g.player.Dy),
	)

	if activeAnimation != nil {
		activeAnimation.Update()
	}

	for _, enemy := range g.enemies {

		enemy.Dx = 0.0
		enemy.Dy = 0.0

		if enemy.FollowsPlayer {
			if enemy.X < g.player.X {
				enemy.Dx += 1
			} else if enemy.X > g.player.X {
				enemy.Dx -= 1
			}

			if enemy.Y < g.player.Y {
				enemy.Dy += 1
			} else if enemy.Y > g.player.Y {
				enemy.Dy -= 1
			}
		}

		enemy.X += enemy.Dx

		checkCollisonHorizontal(enemy.Sprite, g.colliders)

		enemy.Y += enemy.Dy

		checkCollisonVertical(enemy.Sprite, g.colliders)
	}

	clicked := inpututil.IsMouseButtonJustPressed(ebiten.MouseButton0)
	cX, cY := ebiten.CursorPosition()

	deadEnemies := make(map[int]struct{})
	for index, enemy := range g.enemies {
		rect := image.Rect(
			int(enemy.X),
			int(enemy.Y),
			int(enemy.X)+constants.Tilesize,
			int(enemy.Y)+constants.Tilesize,
		)

		if cX > rect.Min.X && cX < rect.Max.X && cY > rect.Min.Y && cY < rect.Max.Y {
			if clicked {
				enemy.CombatComp.Damage(g.player.CombatComp.AttackPower())

				if enemy.CombatComp.Health() <= 0 {
					deadEnemies[index] = struct{}{}
				}
			}
		}
	}
	fmt.Println("deadEnemies", deadEnemies)
	// If there are dead enemies then remove them.
	if len(deadEnemies) > 0 {
		newEnemies := make([]*entities.Enemy, 0)
		for index, enemy := range g.enemies {
			if _, exists := deadEnemies[index]; !exists {
				newEnemies = append(newEnemies, enemy)
			}
		}
		g.enemies = newEnemies
	}

	// for _, potion := range g.potions {
	// 	if g.player.X > potion.X {
	// 		g.player.Health += potion.AmtHeal
	// 		fmt.Printf("Picked up potion! Health: %d\n", g.player.Health)
	// 	}
	// }

	// Add camera to follow player
	g.cam.FollowTarget(g.player.X+8, g.player.Y+8, 320, 240)
	g.cam.Constrain(
		float64(g.tilemapJSON.Layers[0].Width)*16,
		float64(g.tilemapJSON.Layers[0].Height)*16,
		320,
		240,
	)

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{120, 180, 255, 255})

	op := &ebiten.DrawImageOptions{}

	// loop over the layers
	// tilesetColumns := g.tilemapImg.Bounds().Dx() / 16 // Number of tiles per row in tileset
	for layerIndex, layer := range g.tilemapJSON.Layers {
		for index, id := range layer.Data {
			if id == 0 {
				continue
			}

			x := (index % layer.Width) * 16 // tile position x
			y := (index / layer.Width) * 16 // tile position y

			img := g.tilesets[layerIndex].Img(id)

			// *********** TODO: fix the yoffset issue. This is just a temporary fix.
			// Calculate Y offset (adjust based on image height)
			tileHeight := img.Bounds().Dy() // Get actual height
			yOffset := tileHeight - 16      // Adjust if needed

			op.GeoM.Translate(float64(x), float64(y-yOffset)) // Shift up by yOffset
			// ***********

			// op.GeoM.Translate(float64(x), float64(y))
			op.GeoM.Translate(g.cam.X, g.cam.Y)

			// fmt.Printf("Drawing tile ID %d at X=%d, Y=%d\n", id, x, y)
			screen.DrawImage(img, op)

			// srcX := (id - 1) % tilesetColumns * 16
			// srcY := (id - 1) / tilesetColumns * 16

			// screen.DrawImage(
			// 	g.tilemapImg.SubImage(image.Rect(srcX, srcY, srcX+600, srcY+600)).(*ebiten.Image),
			// 	op,
			// )

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
	op.GeoM.Translate(g.cam.X, g.cam.Y)

	activeAnimation := g.player.ActiveAnimation(
		int(g.player.Dx),
		int(g.player.Dy),
	)

	playerFrame := 0
	if activeAnimation != nil {
		playerFrame = activeAnimation.Frame()
	}

	// Draw our player
	screen.DrawImage(
		g.player.Img.SubImage(
			// image.Rect(0, 0, 150, 150),
			g.playerSpriteSheet.Rect(playerFrame),
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
		op.GeoM.Translate(g.cam.X, g.cam.Y)

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
		op.GeoM.Translate(g.cam.X, g.cam.Y)

		screen.DrawImage(
			enemy.Img.SubImage(
				image.Rect(0, 0, 150, 150),
			).(*ebiten.Image),
			op,
		)

		op.GeoM.Reset()

	}

	for _, collider := range g.colliders {
		vector.StrokeRect(
			screen,
			float32(collider.Min.X)+float32(g.cam.X),
			float32(collider.Min.Y)+float32(g.cam.Y),
			float32(collider.Dx()),
			float32(collider.Dy()),
			1.0, color.RGBA{255, 0, 0, 255}, true,
		)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 320, 240
}
