package main

import (
	"log"

	"github.com/bytearena/ecs"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

// Components holds all ECS component references
type Components struct {
	Position    *ecs.Component
	Renderable  *ecs.Component
	Monster     *ecs.Component
	Health      *ecs.Component
	MeleeWeapon *ecs.Component
	Armor       *ecs.Component
	Name        *ecs.Component
	UserMessage *ecs.Component
	Player      *ecs.Component
}

func InitializeWorld(startingLevel Level) (*ecs.Manager, map[string]ecs.Tag, *Components) {
	tags := make(map[string]ecs.Tag)
	manager := ecs.NewManager()

	// Create components struct
	components := &Components{
		Player:      manager.NewComponent(),
		Position:    manager.NewComponent(),
		Renderable:  manager.NewComponent(),
		Monster:     manager.NewComponent(),
		Health:      manager.NewComponent(),
		MeleeWeapon: manager.NewComponent(),
		Armor:       manager.NewComponent(),
		Name:        manager.NewComponent(),
		UserMessage: manager.NewComponent(),
	}

	movable := manager.NewComponent()

	playerImg, _, err := ebitenutil.NewImageFromFile("assets/player.png")
	if err != nil {
		log.Fatal(err)
	}
	skellyImg, _, err := ebitenutil.NewImageFromFile("assets/skelly.png")
	if err != nil {
		log.Fatal(err)
	}
	orcImg, _, err := ebitenutil.NewImageFromFile("assets/orc.png")
	if err != nil {
		log.Fatal(err)
	}

	//Get First Room
	startingRoom := startingLevel.Rooms[0]
	x, y := startingRoom.Center()

	manager.NewEntity().
		AddComponent(components.Player, Player{}).
		AddComponent(components.Renderable, &Renderable{
			Image: playerImg,
		}).
		AddComponent(movable, Movable{}).
		AddComponent(components.Position, &Position{
			X: x,
			Y: y,
		}).
		AddComponent(components.Health, &Health{
			MaxHealth:     30,
			CurrentHealth: 30,
		}).
		AddComponent(components.MeleeWeapon, &MeleeWeapon{
			Name:          "Battle Axe",
			MinimumDamage: 10,
			MaximumDamage: 20,
			ToHitBonus:    3,
		}).
		AddComponent(components.Armor, &Armor{
			Name:       "Plate Armor",
			Defense:    15,
			ArmorClass: 18,
		}).
		AddComponent(components.Name, &Name{Label: "Player"}).
		AddComponent(components.UserMessage, &UserMessage{
			AttackMessage:    "",
			DeadMessage:      "",
			GameStateMessage: "",
		})

	//Add a Monster in each room except the player's room
	for _, room := range startingLevel.Rooms {
		if room.X1 != startingRoom.X1 {
			mX, mY := room.Center()

			//Flip a coin to see what to add...
			mobSpawn := GetDiceRoll(2)

			if mobSpawn == 1 {
				manager.NewEntity().
					AddComponent(components.Monster, &Monster{}).
					AddComponent(components.Renderable, &Renderable{
						Image: orcImg,
					}).
					AddComponent(components.Position, &Position{
						X: mX,
						Y: mY,
					}).
					AddComponent(components.Health, &Health{
						MaxHealth:     30,
						CurrentHealth: 30,
					}).
					AddComponent(components.MeleeWeapon, &MeleeWeapon{
						Name:          "Machete",
						MinimumDamage: 4,
						MaximumDamage: 8,
						ToHitBonus:    1,
					}).
					AddComponent(components.Armor, &Armor{
						Name:       "Leather",
						Defense:    5,
						ArmorClass: 6,
					}).
					AddComponent(components.Name, &Name{Label: "Orc"}).
					AddComponent(components.UserMessage, &UserMessage{
						AttackMessage:    "",
						DeadMessage:      "",
						GameStateMessage: "",
					})
			} else {
				manager.NewEntity().
					AddComponent(components.Monster, &Monster{}).
					AddComponent(components.Renderable, &Renderable{
						Image: skellyImg,
					}).
					AddComponent(components.Position, &Position{
						X: mX,
						Y: mY,
					}).
					AddComponent(components.Health, &Health{
						MaxHealth:     10,
						CurrentHealth: 10,
					}).
					AddComponent(components.MeleeWeapon, &MeleeWeapon{
						Name:          "Short Sword",
						MinimumDamage: 2,
						MaximumDamage: 6,
						ToHitBonus:    0,
					}).
					AddComponent(components.Armor, &Armor{
						Name:       "Bone",
						Defense:    3,
						ArmorClass: 4,
					}).
					AddComponent(components.Name, &Name{Label: "Skeleton"}).
					AddComponent(components.UserMessage, &UserMessage{
						AttackMessage:    "",
						DeadMessage:      "",
						GameStateMessage: "",
					})
			}

		}
	}

	players := ecs.BuildTag(components.Player, components.Position, components.Health, components.MeleeWeapon, components.Armor, components.Name, components.UserMessage)
	tags["players"] = players

	renderables := ecs.BuildTag(components.Renderable, components.Position)
	tags["renderables"] = renderables

	monsters := ecs.BuildTag(components.Monster, components.Position, components.Health, components.MeleeWeapon, components.Armor, components.Name, components.UserMessage)
	tags["monsters"] = monsters

	messengers := ecs.BuildTag(components.UserMessage)
	tags["messengers"] = messengers

	return manager, tags, components
}
