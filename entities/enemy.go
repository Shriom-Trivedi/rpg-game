package entities

import (
	"rpg-game-go/animations"
	"rpg-game-go/components"
)

type EnemyState uint8

type Enemy struct {
	*Sprite
	FollowsPlayer bool
	Animations    map[Direction]*animations.Animation
	CombatComp    *components.EnemyCombat
}

func (e *Enemy) ActiveAnimation(dx, dy int) *animations.Animation {
	if dx > 0 {
		return e.Animations[Right]
	}
	if dx < 0 {
		return e.Animations[Left]
	}
	if dy > 0 {
		return e.Animations[Down]
	}
	if dy < 0 {
		return e.Animations[Up]
	}

	return nil
}
