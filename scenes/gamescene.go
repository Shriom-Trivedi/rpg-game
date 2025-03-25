package scenes

import "github.com/hajimehoshi/ebiten/v2"

type GameScene struct{}

func (g *GameScene) Draw(screen *ebiten.Image) {
}

func (g *GameScene) FirstLoad() {
}

func (g *GameScene) OnEnter() {
}

func (g *GameScene) OnExit() {
}

func (g *GameScene) Update() SceneId {
}

var _ Scene = (*GameScene)(nil)
