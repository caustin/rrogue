package game

import (
	"fmt"
	"github.com/caustin/rrogue/components"
	"github.com/caustin/rrogue/utils"

	"github.com/bytearena/ecs"
)

func AttackSystem(g *Game, attackerPosition *components.Position, defenderPosition *components.Position) {
	var attacker *ecs.QueryResult = nil
	var defender *ecs.QueryResult = nil

	//Get the attacker and defender if either is a player
	for _, playerCombatant := range g.World.QueryPlayers() {
		pos := g.World.GetPosition(playerCombatant)

		if pos.IsEqual(attackerPosition) {
			//This is the attacker
			attacker = playerCombatant
		} else if pos.IsEqual(defenderPosition) {
			//This is the defender
			defender = playerCombatant
		}
	}

	//Get the attacker and defender if either is a monster
	for _, cbt := range g.World.QueryMonsters() {
		pos := g.World.GetPosition(cbt)

		if pos.IsEqual(attackerPosition) {
			//This is the attacker
			attacker = cbt
		} else if pos.IsEqual(defenderPosition) {
			//This is the defender
			defender = cbt
		}

	}
	//If we somehow don't have an attacker or defender, just leave
	if attacker == nil || defender == nil {
		return
	}
	//Grab the required information
	defenderArmor := g.World.GetArmor(defender)
	defenderHealth := g.World.GetHealth(defender)
	defenderName := g.World.GetName(defender).Label
	defenderMessage := g.World.GetUserMessage(defender)

	attackerWeapon := g.World.GetMeleeWeapon(attacker)
	attackerName := g.World.GetName(attacker).Label
	attackerMessage := g.World.GetUserMessage(attacker)

	//if the attacker is dead, don't let them attackerWeapon
	if g.World.GetHealth(attacker).CurrentHealth <= 0 {
		return
	}
	//Roll a d10 to hit
	toHitRoll := utils.GetDiceRoll(10)

	if toHitRoll+attackerWeapon.ToHitBonus > defenderArmor.ArmorClass {
		//It's a hit!
		damageRoll := utils.GetRandomBetween(attackerWeapon.MinimumDamage, attackerWeapon.MaximumDamage)

		damageDone := damageRoll - defenderArmor.Defense
		//Let's not have the weapon heal the defender
		if damageDone < 0 {
			damageDone = 0
		}
		defenderHealth.CurrentHealth -= damageDone
		attackerMessage.AttackMessage = fmt.Sprintf("%s swings %s at %s and hits for %d health.\n", attackerName, attackerWeapon.Name, defenderName, damageDone)

		if defenderHealth.CurrentHealth <= 0 {
			defenderMessage.DeadMessage = fmt.Sprintf("%s has died!\n", defenderName)
			if defenderName == "Player" {
				defenderMessage.GameStateMessage = "Game Over!\n"
				g.Turn = GameOver
			} else {
				// Monster died - clean up tile and dispose entity
				level := g.Map.CurrentLevel
				pos := g.World.GetPosition(defender)
				tile := level.Tiles[level.GetIndexFromXY(pos.X, pos.Y)]
				tile.Blocked = false
				g.World.DisposeEntity(defender)
			}
		}

	} else {
		attackerMessage.AttackMessage = fmt.Sprintf("%s swings %s at %s and misses.\n", attackerName, attackerWeapon.Name, defenderName)
	}
}
