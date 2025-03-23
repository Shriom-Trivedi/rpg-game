package entities

import "rpg-game-go/components"

type Enemy struct {
	*Sprite
	FollowsPlayer bool
	CombatComp    *components.BasicCombat
}