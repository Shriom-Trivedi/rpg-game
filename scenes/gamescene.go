package scenes

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"rpg-game-go/animations"
	"rpg-game-go/camera"
	"rpg-game-go/components"
	"rpg-game-go/constants"
	"rpg-game-go/entities"
	"rpg-game-go/spritesheet"
	"rpg-game-go/tilemap"
	"rpg-game-go/tileset"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

func checkCollisonHorizontal(sprite *entities.Sprite, colliders []image.Rectangle) {
	for _, collider := range colliders {
		if collider.Overlaps(
			image.Rect(int(sprite.X), int(sprite.Y), int(sprite.X)+constants.Tilesize, int(sprite.Y)+constants.Tilesize),
		) {

			if sprite.Dx > 0.0 {
				sprite.X = float64(collider.Min.X) - constants.Tilesize
			} else if sprite.Dx < 0.0 {
				sprite.X = float64(collider.Max.X)
			}
		}
	}
}

func checkCollisonVertical(sprite *entities.Sprite, colliders []image.Rectangle) {
	for _, collider := range colliders {
		if collider.Overlaps(
			image.Rect(int(sprite.X), int(sprite.Y), int(sprite.X)+constants.Tilesize, int(sprite.Y)+constants.Tilesize),
		) {

			if sprite.Dy > 0.0 {
				sprite.Y = float64(collider.Min.Y) - constants.Tilesize
			} else if sprite.Dy < 0.0 {
				sprite.Y = float64(collider.Max.Y)
			}
		}
	}
}

type GameScene struct {
	player            *entities.Player
	playerSpriteSheet *spritesheet.Spritesheet
	enemySpriteSheet  *spritesheet.Spritesheet
	animationFrame    int
	tilemapJSON       *tilemap.TilemapJSON
	tilemapImg        *ebiten.Image
	cam               *camera.Camera
	enemies           []*entities.Enemy
	potions           []*entities.Potion
	colliders         []image.Rectangle
	tilesets          []tileset.Tileset
}

func NewGameScene() *GameScene {
	return &GameScene{
		player:            nil,
		playerSpriteSheet: nil,
		enemies:           make([]*entities.Enemy, 0),
		potions:           make([]*entities.Potion, 0),
		tilemapJSON:       nil,
		tilemapImg:        nil,
		cam:               nil,
		colliders:         make([]image.Rectangle, 0),
	}
}

func (g *GameScene) Draw(screen *ebiten.Image) {
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
	if g.player.CombatComp.Attacking() {
		playerFrame = g.player.CombatAnimation(entities.MouseLeftClick).Frame()
	} else if activeAnimation != nil {
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

		activeAnimation := enemy.ActiveAnimation(
			int(enemy.Dx),
			int(enemy.Dy),
		)

		enemyFrame := 0
		if activeAnimation != nil {
			enemyFrame = activeAnimation.Frame()
		}

		screen.DrawImage(
			enemy.Img.SubImage(
				// image.Rect(0, 0, 150, 150),
				g.enemySpriteSheet.Rect(enemyFrame),
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

func (g *GameScene) FirstLoad() {
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
	tilemapJSON, err := tilemap.NewTilemapJSON("assets/maps/spawn.tmj")
	if err != nil {
		log.Fatal(err)
	}

	// Generate tilesets
	tilesets, err := tilemapJSON.GenTilesets()
	if err != nil {
		log.Fatal(err)
	}

	playerSpriteSheet := spritesheet.NewSpriteSheet(6, 8, 192)
	enemySpriteSheet := spritesheet.NewSpriteSheet(6, 5, 192)

	g.player = &entities.Player{
		Sprite: &entities.Sprite{
			Img: playerImg,
			X:   50,
			Y:   50,
		},
		Health: 5,
		Animations: map[entities.Direction]*animations.Animation{
			entities.Right:          animations.NewAnimation(6, 11, 1, 8.0),
			entities.Left:           animations.NewAnimation(48, 53, 1, 8.0),
			entities.Down:           animations.NewAnimation(26, 30, 3, 8.0),
			entities.Up:             animations.NewAnimation(38, 42, 3, 8.0),
			entities.MouseLeftClick: animations.NewAnimation(12, 17, 1, 1.0),
		},
		CombatComp: components.NewBasicCombat(3, 1),
	}

	g.playerSpriteSheet = playerSpriteSheet
	g.enemySpriteSheet = enemySpriteSheet

	g.enemies = []*entities.Enemy{
		{
			Sprite: &entities.Sprite{
				Img: goblinFireImg,
				X:   150,
				Y:   150,
			},
			FollowsPlayer: true,
			Animations: map[entities.Direction]*animations.Animation{
				entities.Right: animations.NewAnimation(7, 12, 1, 8.0),
				entities.Left:  animations.NewAnimation(7, 12, 1, 8.0),
				entities.Up:    animations.NewAnimation(7, 12, 1, 8.0),
				entities.Down:  animations.NewAnimation(7, 12, 1, 8.0),
			},
			CombatComp: components.NewEnemyCombat(3, 1, 30),
		},
		{
			Sprite: &entities.Sprite{
				Img: goblinFireImg,
				X:   150,
				Y:   100,
			},
			FollowsPlayer: false,
			Animations: map[entities.Direction]*animations.Animation{
				entities.Right: animations.NewAnimation(7, 12, 1, 8.0),
				entities.Left:  animations.NewAnimation(7, 12, 1, 8.0),
				entities.Down:  animations.NewAnimation(7, 12, 1, 8.0),
				entities.Up:    animations.NewAnimation(7, 12, 1, 8.0),
			},
			CombatComp: components.NewEnemyCombat(3, 1, 30),
		},
	}

	g.potions = []*entities.Potion{
		{
			Sprite: &entities.Sprite{
				Img: potionImg,
				X:   120,
				Y:   120,
			},
			AmtHeal: 5,
		},
	}

	g.tilemapJSON = tilemapJSON
	g.tilemapImg = tilemapImg
	g.tilesets = tilesets
	g.cam = camera.NewCamera(50, 50)

	g.colliders = []image.Rectangle{
		image.Rect(100, 100, 116, 116),
	}

}

func (g *GameScene) OnEnter() {
}

func (g *GameScene) OnExit() {
}

func (g *GameScene) Update() SceneId {
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

	clicked := inpututil.IsMouseButtonJustPressed(ebiten.MouseButton0)

	if clicked {
		g.player.CombatComp.Attack()
		g.player.CombatAnimation(entities.MouseLeftClick).Reset()
	}

	if g.player.CombatComp.Attacking() {
		playerCombatAnimation := g.player.CombatAnimation(entities.MouseLeftClick)
		playerCombatAnimation.Update()

		if playerCombatAnimation.IsLastFrame() {
			g.player.CombatComp.AttackingStop()
		}
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

		enemyAnimation := enemy.ActiveAnimation(
			int(enemy.Dx),
			int(enemy.Dy),
		)
		if enemyAnimation != nil {
			enemyAnimation.Update()
		}
	}

	cX, cY := ebiten.CursorPosition()
	cX -= int(g.cam.X)
	cY -= int(g.cam.Y)

	g.player.CombatComp.Update()

	// player rectangle
	pRect := image.Rect(
		int(g.player.X),
		int(g.player.Y),
		int(g.player.X)+constants.Tilesize,
		int(g.player.Y)+constants.Tilesize,
	)

	deadEnemies := make(map[int]struct{})
	for index, enemy := range g.enemies {
		enemy.CombatComp.Update()
		rect := image.Rect(
			int(enemy.X),
			int(enemy.Y),
			int(enemy.X)+constants.Tilesize,
			int(enemy.Y)+constants.Tilesize,
		)

		// if enemy overlaps player
		if rect.Overlaps(pRect) {
			if enemy.CombatComp.Attack() {
				g.player.CombatComp.Damage(enemy.CombatComp.AttackPower())

				if g.player.CombatComp.Health() <= 0 {
					fmt.Println("The Player has died...   ")
				}
			}
		}

		if cX > rect.Min.X && cX < rect.Max.X && cY > rect.Min.Y && cY < rect.Max.Y {
			if clicked {
				enemy.CombatComp.Damage(g.player.CombatComp.AttackPower())

				if enemy.CombatComp.Health() <= 0 {
					deadEnemies[index] = struct{}{}
				}
			}
		}
	}

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

	return GameSceneId
}

var _ Scene = (*GameScene)(nil)
