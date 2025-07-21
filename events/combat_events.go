package events

import (
	"github.com/bytearena/ecs"
	"github.com/caustin/rrogue/components"
)

// AttackEvent represents an attack between two entities
type AttackEvent struct {
	BaseEvent
	Attacker    *ecs.QueryResult
	Defender    *ecs.QueryResult
	AttackerPos *components.Position
	DefenderPos *components.Position
	ToHitRoll   int
	Hit         bool
}

func NewAttackEvent(attacker, defender *ecs.QueryResult, attackerPos, defenderPos *components.Position, toHitRoll int, hit bool) *AttackEvent {
	return &AttackEvent{
		BaseEvent:   NewBaseEvent(AttackEventType),
		Attacker:    attacker,
		Defender:    defender,
		AttackerPos: attackerPos,
		DefenderPos: defenderPos,
		ToHitRoll:   toHitRoll,
		Hit:         hit,
	}
}

// DamageEvent represents damage being dealt
type DamageEvent struct {
	BaseEvent
	Target       *ecs.QueryResult
	DamageAmount int
	DamageSource string
	IsFatal      bool
}

func NewDamageEvent(target *ecs.QueryResult, damageAmount int, damageSource string, isFatal bool) *DamageEvent {
	return &DamageEvent{
		BaseEvent:    NewBaseEvent(DamageEventType),
		Target:       target,
		DamageAmount: damageAmount,
		DamageSource: damageSource,
		IsFatal:      isFatal,
	}
}

// DeathEvent represents an entity dying
type DeathEvent struct {
	BaseEvent
	Entity   *ecs.QueryResult
	Position *components.Position
	IsPlayer bool
}

func NewDeathEvent(entity *ecs.QueryResult, position *components.Position, isPlayer bool) *DeathEvent {
	return &DeathEvent{
		BaseEvent: NewBaseEvent(DeathEventType),
		Entity:    entity,
		Position:  position,
		IsPlayer:  isPlayer,
	}
}

// MoveEvent represents entity movement
type MoveEvent struct {
	BaseEvent
	Entity   *ecs.QueryResult
	FromPos  *components.Position
	ToPos    *components.Position
	IsPlayer bool
}

func NewMoveEvent(entity *ecs.QueryResult, fromPos, toPos *components.Position, isPlayer bool) *MoveEvent {
	return &MoveEvent{
		BaseEvent: NewBaseEvent(MoveEventType),
		Entity:    entity,
		FromPos:   fromPos,
		ToPos:     toPos,
		IsPlayer:  isPlayer,
	}
}
