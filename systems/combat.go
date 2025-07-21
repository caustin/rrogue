package systems

import (
	"fmt"
	"github.com/bytearena/ecs"
	"github.com/caustin/rrogue/components"
	"github.com/caustin/rrogue/events"
	"github.com/caustin/rrogue/utils"
	"github.com/caustin/rrogue/world"
)

// CombatSystem handles all combat-related operations
type CombatSystem struct {
	world    world.WorldService
	eventBus *events.EventBus
}

// NewCombatSystem creates a new combat system with dependencies
func NewCombatSystem(world world.WorldService, eventBus *events.EventBus) *CombatSystem {
	return &CombatSystem{
		world:    world,
		eventBus: eventBus,
	}
}

// RegisterHandlers subscribes the combat system to relevant events
func (cs *CombatSystem) RegisterHandlers() {
	cs.eventBus.Subscribe(events.AttackEventType, cs.HandleAttack)
	cs.eventBus.Subscribe(events.DamageEventType, cs.HandleDamage)
}

// HandleAttack processes attack events and determines hit/miss
func (cs *CombatSystem) HandleAttack(event events.Event) {
	attackEvent := event.(*events.AttackEvent)

	// Get component data
	defenderArmor := cs.world.GetArmor(attackEvent.Defender)
	attackerWeapon := cs.world.GetMeleeWeapon(attackEvent.Attacker)
	attackerName := cs.world.GetName(attackEvent.Attacker).Label
	defenderName := cs.world.GetName(attackEvent.Defender).Label
	attackerMessage := cs.world.GetUserMessage(attackEvent.Attacker)

	// Check if attacker is alive
	if cs.world.GetHealth(attackEvent.Attacker).CurrentHealth <= 0 {
		return
	}

	// Determine hit/miss based on the attack event data
	if attackEvent.Hit {
		// Calculate damage
		damageRoll := utils.GetRandomBetween(attackerWeapon.MinimumDamage, attackerWeapon.MaximumDamage)
		damageDone := damageRoll - defenderArmor.Defense

		// Ensure no negative damage (healing)
		if damageDone < 0 {
			damageDone = 0
		}

		// Set attack message
		attackerMessage.AttackMessage = fmt.Sprintf("%s swings %s at %s and hits for %d health.\n",
			attackerName, attackerWeapon.Name, defenderName, damageDone)

		// Publish damage event
		damageEvent := events.NewDamageEvent(attackEvent.Defender, damageDone, attackerWeapon.Name, false)
		cs.eventBus.Publish(damageEvent)

	} else {
		// Miss
		attackerMessage.AttackMessage = fmt.Sprintf("%s swings %s at %s and misses.\n",
			attackerName, attackerWeapon.Name, defenderName)
	}
}

// HandleDamage processes damage events and applies damage
func (cs *CombatSystem) HandleDamage(event events.Event) {
	damageEvent := event.(*events.DamageEvent)

	// Apply damage
	defenderHealth := cs.world.GetHealth(damageEvent.Target)
	defenderHealth.CurrentHealth -= damageEvent.DamageAmount

	// Check for death
	if defenderHealth.CurrentHealth <= 0 {
		defenderPos := cs.world.GetPosition(damageEvent.Target)
		defenderName := cs.world.GetName(damageEvent.Target).Label
		isPlayer := defenderName == "Player"

		// Set death message
		defenderMessage := cs.world.GetUserMessage(damageEvent.Target)
		defenderMessage.DeadMessage = fmt.Sprintf("%s has died!\n", defenderName)

		// Handle death immediately (temporary until we have proper event handlers)
		if isPlayer {
			defenderMessage.GameStateMessage = "Game Over!\n"
			// TODO: Set game over state via event
		} else {
			// Clean up the monster entity
			cs.world.DisposeEntity(damageEvent.Target)
			// Note: Map tile unblocking will need to be handled externally
			// until we implement a proper MapSystem
		}

		// Publish death event for future event handlers
		deathEvent := events.NewDeathEvent(damageEvent.Target, defenderPos, isPlayer)
		cs.eventBus.Publish(deathEvent)
	}
}

// ProcessAttack is a helper function to initiate an attack between two positions
func (cs *CombatSystem) ProcessAttack(attackerPos, defenderPos *components.Position) {
	var attacker, defender *ecs.QueryResult = nil, nil

	// Find attacker and defender entities at the given positions
	for _, player := range cs.world.QueryPlayers() {
		pos := cs.world.GetPosition(player)
		if pos.IsEqual(attackerPos) {
			attacker = player
		}
		if pos.IsEqual(defenderPos) {
			defender = player
		}
	}

	for _, monster := range cs.world.QueryMonsters() {
		pos := cs.world.GetPosition(monster)
		if pos.IsEqual(attackerPos) {
			attacker = monster
		}
		if pos.IsEqual(defenderPos) {
			defender = monster
		}
	}

	// Ensure we have both attacker and defender
	if attacker == nil || defender == nil {
		return
	}

	// Roll to hit
	toHitRoll := utils.GetDiceRoll(10)
	attackerWeapon := cs.world.GetMeleeWeapon(attacker)
	defenderArmor := cs.world.GetArmor(defender)
	hit := toHitRoll+attackerWeapon.ToHitBonus > defenderArmor.ArmorClass

	// Publish attack event
	attackEvent := events.NewAttackEvent(attacker, defender, attackerPos, defenderPos, toHitRoll, hit)
	cs.eventBus.Publish(attackEvent)
}
