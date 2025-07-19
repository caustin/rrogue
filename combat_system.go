package main

import (
	"fmt"

	"github.com/bytearena/ecs"
)

func AttackSystem(g *Game, attackerPosition *Position, defenderPosition *Position) {
	var attacker *ecs.QueryResult = nil
	var defender *ecs.QueryResult = nil

	//Get the attacker and defender if either is a player
	for _, playerCombatant := range g.World.Query(g.WorldTags["players"]) {
		pos := playerCombatant.Components[g.Components.Position].(*Position)

		if pos.IsEqual(attackerPosition) {
			//This is the attacker
			attacker = playerCombatant
		} else if pos.IsEqual(defenderPosition) {
			//This is the defender
			defender = playerCombatant
		}
	}

	//Get the attacker and defender if either is a monster
	for _, cbt := range g.World.Query(g.WorldTags["monsters"]) {
		pos := cbt.Components[g.Components.Position].(*Position)

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
	defenderArmor := defender.Components[g.Components.Armor].(*Armor)
	defenderHealth := defender.Components[g.Components.Health].(*Health)
	defenderName := defender.Components[g.Components.Name].(*Name).Label
	defenderMessage := defender.Components[g.Components.UserMessage].(*UserMessage)

	attackerWeapon := attacker.Components[g.Components.MeleeWeapon].(*MeleeWeapon)
	attackerName := attacker.Components[g.Components.Name].(*Name).Label
	attackerMessage := attacker.Components[g.Components.UserMessage].(*UserMessage)

	//if the attacker is dead, don't let them attackerWeapon
	if attacker.Components[g.Components.Health].(*Health).CurrentHealth <= 0 {
		return
	}
	//Roll a d10 to hit
	toHitRoll := GetDiceRoll(10)

	if toHitRoll+attackerWeapon.ToHitBonus > defenderArmor.ArmorClass {
		//It's a hit!
		damageRoll := GetRandomBetween(attackerWeapon.MinimumDamage, attackerWeapon.MaximumDamage)

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
				pos := defender.Components[g.Components.Position].(*Position)
				tile := level.Tiles[level.GetIndexFromXY(pos.X, pos.Y)]
				tile.Blocked = false
				g.World.DisposeEntity(defender.Entity)
			}
		}

	} else {
		attackerMessage.AttackMessage = fmt.Sprintf("%s swings %s at %s and misses.\n", attackerName, attackerWeapon.Name, defenderName)
	}
}
