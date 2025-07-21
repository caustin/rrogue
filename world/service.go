package world

import (
	"github.com/bytearena/ecs"
	"github.com/caustin/rrogue/components"
)

// WorldService provides an interface for systems to interact with the ECS world
type WorldService interface {
	// Entity queries
	QueryPlayers() []*ecs.QueryResult
	QueryMonsters() []*ecs.QueryResult
	QueryRenderables() []*ecs.QueryResult
	QueryMessengers() []*ecs.QueryResult

	// Component access
	GetPosition(entity *ecs.QueryResult) *components.Position
	GetHealth(entity *ecs.QueryResult) *components.Health
	GetArmor(entity *ecs.QueryResult) *components.Armor
	GetMeleeWeapon(entity *ecs.QueryResult) *components.MeleeWeapon
	GetName(entity *ecs.QueryResult) *components.Name
	GetUserMessage(entity *ecs.QueryResult) *components.UserMessage
	GetRenderable(entity *ecs.QueryResult) *components.Renderable

	// Entity lifecycle
	DisposeEntity(entity *ecs.QueryResult)

	// Raw access for advanced use cases
	GetManager() *ecs.Manager
}
