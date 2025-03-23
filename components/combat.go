package components

type Combat interface {
	Health() int
	AttackPower() int
	Attacking() bool
	Attack()
	Damage(amount int)
}

type BasicCombat struct {
	health      int
	attackPower int
	attacking   bool
}

func NewBasicCombat(health, attackPower int) *BasicCombat {
	return &BasicCombat{
		health,
		attackPower,
		false,
	}
}

func (b *BasicCombat) AttackPower() int {
	return b.attackPower
}

func (b *BasicCombat) Damage(amount int) {
	b.health -= amount
}

func (b *BasicCombat) Health() int {
	return b.health
}

func (b *BasicCombat) Attacking() bool {
	return b.attacking
}

func (b *BasicCombat) Attack() {

}

var _ Combat = (*BasicCombat)(nil)
