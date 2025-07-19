package main

import "testing"

func TestManhattanDistance(t *testing.T) {
	tests := []struct {
		name     string
		pos1     Position
		pos2     Position
		expected int
	}{
		{
			name:     "same position",
			pos1:     Position{X: 5, Y: 5},
			pos2:     Position{X: 5, Y: 5},
			expected: 0,
		},
		{
			name:     "adjacent horizontal",
			pos1:     Position{X: 5, Y: 5},
			pos2:     Position{X: 6, Y: 5},
			expected: 1,
		},
		{
			name:     "adjacent vertical",
			pos1:     Position{X: 5, Y: 5},
			pos2:     Position{X: 5, Y: 6},
			expected: 1,
		},
		{
			name:     "diagonal distance",
			pos1:     Position{X: 0, Y: 0},
			pos2:     Position{X: 3, Y: 4},
			expected: 7,
		},
		{
			name:     "negative coordinates",
			pos1:     Position{X: -2, Y: -3},
			pos2:     Position{X: 1, Y: 2},
			expected: 8,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.pos1.GetManhattanDistance(&tt.pos2)
			if result != tt.expected {
				t.Errorf("GetManhattanDistance() = %d, expected %d", result, tt.expected)
			}

			// Test symmetry
			result2 := tt.pos2.GetManhattanDistance(&tt.pos1)
			if result2 != tt.expected {
				t.Errorf("GetManhattanDistance() symmetry failed: %d != %d", result, result2)
			}
		})
	}
}

func TestPositionEquality(t *testing.T) {
	tests := []struct {
		name     string
		pos1     Position
		pos2     Position
		expected bool
	}{
		{
			name:     "equal positions",
			pos1:     Position{X: 5, Y: 5},
			pos2:     Position{X: 5, Y: 5},
			expected: true,
		},
		{
			name:     "different x",
			pos1:     Position{X: 5, Y: 5},
			pos2:     Position{X: 6, Y: 5},
			expected: false,
		},
		{
			name:     "different y",
			pos1:     Position{X: 5, Y: 5},
			pos2:     Position{X: 5, Y: 6},
			expected: false,
		},
		{
			name:     "completely different",
			pos1:     Position{X: 1, Y: 2},
			pos2:     Position{X: 8, Y: 9},
			expected: false,
		},
		{
			name:     "negative coordinates equal",
			pos1:     Position{X: -5, Y: -10},
			pos2:     Position{X: -5, Y: -10},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.pos1.IsEqual(&tt.pos2)
			if result != tt.expected {
				t.Errorf("IsEqual() = %v, expected %v", result, tt.expected)
			}

			// Test symmetry
			result2 := tt.pos2.IsEqual(&tt.pos1)
			if result2 != tt.expected {
				t.Errorf("IsEqual() symmetry failed: %v != %v", result, result2)
			}
		})
	}
}

func TestHealthComponent(t *testing.T) {
	tests := []struct {
		name           string
		maxHealth      int
		currentHealth  int
		damage         int
		expectedHealth int
		shouldBeDead   bool
	}{
		{
			name:           "healthy character takes damage",
			maxHealth:      100,
			currentHealth:  100,
			damage:         30,
			expectedHealth: 70,
			shouldBeDead:   false,
		},
		{
			name:           "character dies from damage",
			maxHealth:      50,
			currentHealth:  20,
			damage:         25,
			expectedHealth: -5,
			shouldBeDead:   true,
		},
		{
			name:           "no damage taken",
			maxHealth:      75,
			currentHealth:  50,
			damage:         0,
			expectedHealth: 50,
			shouldBeDead:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			health := Health{
				MaxHealth:     tt.maxHealth,
				CurrentHealth: tt.currentHealth,
			}

			// Simulate taking damage
			health.CurrentHealth -= tt.damage

			if health.CurrentHealth != tt.expectedHealth {
				t.Errorf("Expected health %d, got %d", tt.expectedHealth, health.CurrentHealth)
			}

			isDead := health.CurrentHealth <= 0
			if isDead != tt.shouldBeDead {
				t.Errorf("Expected dead status %v, got %v", tt.shouldBeDead, isDead)
			}
		})
	}
}

func TestMeleeWeaponComponent(t *testing.T) {
	weapon := MeleeWeapon{
		Name:          "Test Sword",
		MinimumDamage: 5,
		MaximumDamage: 15,
		ToHitBonus:    3,
	}

	// Test that weapon properties are set correctly
	if weapon.Name != "Test Sword" {
		t.Errorf("Expected weapon name 'Test Sword', got '%s'", weapon.Name)
	}

	if weapon.MinimumDamage != 5 {
		t.Errorf("Expected minimum damage 5, got %d", weapon.MinimumDamage)
	}

	if weapon.MaximumDamage != 15 {
		t.Errorf("Expected maximum damage 15, got %d", weapon.MaximumDamage)
	}

	if weapon.ToHitBonus != 3 {
		t.Errorf("Expected to hit bonus 3, got %d", weapon.ToHitBonus)
	}

	// Test damage range is valid
	if weapon.MinimumDamage > weapon.MaximumDamage {
		t.Error("Minimum damage should not exceed maximum damage")
	}
}

func TestArmorComponent(t *testing.T) {
	armor := Armor{
		Name:       "Test Plate",
		Defense:    10,
		ArmorClass: 15,
	}

	// Test that armor properties are set correctly
	if armor.Name != "Test Plate" {
		t.Errorf("Expected armor name 'Test Plate', got '%s'", armor.Name)
	}

	if armor.Defense != 10 {
		t.Errorf("Expected defense 10, got %d", armor.Defense)
	}

	if armor.ArmorClass != 15 {
		t.Errorf("Expected armor class 15, got %d", armor.ArmorClass)
	}
}

func TestUserMessageComponent(t *testing.T) {
	msg := UserMessage{
		AttackMessage:    "Player attacks!",
		DeadMessage:      "Monster dies!",
		GameStateMessage: "Game Over!",
	}

	// Test message properties
	if msg.AttackMessage != "Player attacks!" {
		t.Errorf("Expected attack message 'Player attacks!', got '%s'", msg.AttackMessage)
	}

	if msg.DeadMessage != "Monster dies!" {
		t.Errorf("Expected dead message 'Monster dies!', got '%s'", msg.DeadMessage)
	}

	if msg.GameStateMessage != "Game Over!" {
		t.Errorf("Expected game state message 'Game Over!', got '%s'", msg.GameStateMessage)
	}

	// Test clearing messages
	msg.AttackMessage = ""
	msg.DeadMessage = ""
	msg.GameStateMessage = ""

	if msg.AttackMessage != "" || msg.DeadMessage != "" || msg.GameStateMessage != "" {
		t.Error("Messages should be clearable")
	}
}
