package entities

import (
	"rpg-game-go/animations"
	"rpg-game-go/components"
)

type PlayerState uint8

type Player struct {
	*Sprite
	Health     uint
	Animations map[Direction]*animations.Animation
	CombatComp *components.BasicCombat
}

func (p *Player) ActiveAnimation(dx, dy int) *animations.Animation {
	if dx > 0 {
		return p.Animations[Right]
	}
	if dx < 0 {
		return p.Animations[Left]
	}
	if dy > 0 {
		return p.Animations[Down]
	}
	if dy < 0 {
		return p.Animations[Up]
	}

	return nil
}

func (p *Player) CombatAnimation(button Direction) *animations.Animation {
	return p.Animations[Direction(button)]
}
