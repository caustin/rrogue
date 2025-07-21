package world

import (
	"github.com/bytearena/ecs"
	"github.com/caustin/rrogue/components"
	"github.com/caustin/rrogue/level"
	"github.com/caustin/rrogue/utils"
	"log"

	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

// ComponentReferences holds all ECS component references
type ComponentReferences struct {
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

// GameWorld implements WorldService and manages the ECS world
type GameWorld struct {
	manager    *ecs.Manager
	tags       map[string]ecs.Tag
	components *ComponentReferences
}

// NewGameWorld creates a new GameWorld with initialized entities
func NewGameWorld(startingLevel level.Level) *GameWorld {
	manager, tags, components := initializeWorld(startingLevel)
	return &GameWorld{
		manager:    manager,
		tags:       tags,
		components: components,
	}
}

// QueryPlayers returns all player entities
func (w *GameWorld) QueryPlayers() []*ecs.QueryResult {
	return w.manager.Query(w.tags["players"])
}

// QueryMonsters returns all monster entities
func (w *GameWorld) QueryMonsters() []*ecs.QueryResult {
	return w.manager.Query(w.tags["monsters"])
}

// QueryRenderables returns all renderable entities
func (w *GameWorld) QueryRenderables() []*ecs.QueryResult {
	return w.manager.Query(w.tags["renderables"])
}

// QueryMessengers returns all entities with user messages
func (w *GameWorld) QueryMessengers() []*ecs.QueryResult {
	return w.manager.Query(w.tags["messengers"])
}

// GetPosition returns the position component of an entity
func (w *GameWorld) GetPosition(entity *ecs.QueryResult) *components.Position {
	return entity.Components[w.components.Position].(*components.Position)
}

// GetHealth returns the health component of an entity
func (w *GameWorld) GetHealth(entity *ecs.QueryResult) *components.Health {
	return entity.Components[w.components.Health].(*components.Health)
}

// GetArmor returns the armor component of an entity
func (w *GameWorld) GetArmor(entity *ecs.QueryResult) *components.Armor {
	return entity.Components[w.components.Armor].(*components.Armor)
}

// GetMeleeWeapon returns the melee weapon component of an entity
func (w *GameWorld) GetMeleeWeapon(entity *ecs.QueryResult) *components.MeleeWeapon {
	return entity.Components[w.components.MeleeWeapon].(*components.MeleeWeapon)
}

// GetName returns the name component of an entity
func (w *GameWorld) GetName(entity *ecs.QueryResult) *components.Name {
	return entity.Components[w.components.Name].(*components.Name)
}

// GetUserMessage returns the user message component of an entity
func (w *GameWorld) GetUserMessage(entity *ecs.QueryResult) *components.UserMessage {
	return entity.Components[w.components.UserMessage].(*components.UserMessage)
}

// GetRenderable returns the renderable component of an entity
func (w *GameWorld) GetRenderable(entity *ecs.QueryResult) *components.Renderable {
	return entity.Components[w.components.Renderable].(*components.Renderable)
}

// DisposeEntity removes an entity from the world
func (w *GameWorld) DisposeEntity(entity *ecs.QueryResult) {
	w.manager.DisposeEntity(entity.Entity)
}

// GetManager returns the raw ECS manager for advanced use cases
func (w *GameWorld) GetManager() *ecs.Manager {
	return w.manager
}

// initializeWorld creates and populates the ECS world (moved from game package)
func initializeWorld(startingLevel level.Level) (*ecs.Manager, map[string]ecs.Tag, *ComponentReferences) {
	tags := make(map[string]ecs.Tag)
	manager := ecs.NewManager()

	// Create components struct
	cr := &ComponentReferences{
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
		AddComponent(cr.Player, components.Player{}).
		AddComponent(cr.Renderable, &components.Renderable{
			Image: playerImg,
		}).
		AddComponent(movable, components.Movable{}).
		AddComponent(cr.Position, &components.Position{
			X: x,
			Y: y,
		}).
		AddComponent(cr.Health, &components.Health{
			MaxHealth:     30,
			CurrentHealth: 30,
		}).
		AddComponent(cr.MeleeWeapon, &components.MeleeWeapon{
			Name:          "Battle Axe",
			MinimumDamage: 10,
			MaximumDamage: 20,
			ToHitBonus:    3,
		}).
		AddComponent(cr.Armor, &components.Armor{
			Name:       "Plate Armor",
			Defense:    15,
			ArmorClass: 18,
		}).
		AddComponent(cr.Name, &components.Name{Label: "Player"}).
		AddComponent(cr.UserMessage, &components.UserMessage{
			AttackMessage:    "",
			DeadMessage:      "",
			GameStateMessage: "",
		})

	//Add a Monster in each room except the player's room
	for _, room := range startingLevel.Rooms {
		if room.X1 != startingRoom.X1 {
			mX, mY := room.Center()

			//Flip a coin to see what to add...
			mobSpawn := utils.GetDiceRoll(2)

			if mobSpawn == 1 {
				manager.NewEntity().
					AddComponent(cr.Monster, &components.Monster{}).
					AddComponent(cr.Renderable, &components.Renderable{
						Image: orcImg,
					}).
					AddComponent(cr.Position, &components.Position{
						X: mX,
						Y: mY,
					}).
					AddComponent(cr.Health, &components.Health{
						MaxHealth:     30,
						CurrentHealth: 30,
					}).
					AddComponent(cr.MeleeWeapon, &components.MeleeWeapon{
						Name:          "Machete",
						MinimumDamage: 4,
						MaximumDamage: 8,
						ToHitBonus:    1,
					}).
					AddComponent(cr.Armor, &components.Armor{
						Name:       "Leather",
						Defense:    5,
						ArmorClass: 6,
					}).
					AddComponent(cr.Name, &components.Name{Label: "Orc"}).
					AddComponent(cr.UserMessage, &components.UserMessage{
						AttackMessage:    "",
						DeadMessage:      "",
						GameStateMessage: "",
					})
			} else {
				manager.NewEntity().
					AddComponent(cr.Monster, &components.Monster{}).
					AddComponent(cr.Renderable, &components.Renderable{
						Image: skellyImg,
					}).
					AddComponent(cr.Position, &components.Position{
						X: mX,
						Y: mY,
					}).
					AddComponent(cr.Health, &components.Health{
						MaxHealth:     10,
						CurrentHealth: 10,
					}).
					AddComponent(cr.MeleeWeapon, &components.MeleeWeapon{
						Name:          "Short Sword",
						MinimumDamage: 2,
						MaximumDamage: 6,
						ToHitBonus:    0,
					}).
					AddComponent(cr.Armor, &components.Armor{
						Name:       "Bone",
						Defense:    3,
						ArmorClass: 4,
					}).
					AddComponent(cr.Name, &components.Name{Label: "Skeleton"}).
					AddComponent(cr.UserMessage, &components.UserMessage{
						AttackMessage:    "",
						DeadMessage:      "",
						GameStateMessage: "",
					})
			}

		}
	}

	players := ecs.BuildTag(cr.Player, cr.Position, cr.Health, cr.MeleeWeapon, cr.Armor, cr.Name, cr.UserMessage)
	tags["players"] = players

	renderables := ecs.BuildTag(cr.Renderable, cr.Position)
	tags["renderables"] = renderables

	monsters := ecs.BuildTag(cr.Monster, cr.Position, cr.Health, cr.MeleeWeapon, cr.Armor, cr.Name, cr.UserMessage)
	tags["monsters"] = monsters

	messengers := ecs.BuildTag(cr.UserMessage)
	tags["messengers"] = messengers

	return manager, tags, cr
}
