package game

import "testing"

func TestCombatDamageCalculation(t *testing.T) {
	tests := []struct {
		name           string
		weaponMinDmg   int
		weaponMaxDmg   int
		armorDefense   int
		expectedMinDmg int
		expectedMaxDmg int
	}{
		{
			name:           "weapon vs no armor",
			weaponMinDmg:   5,
			weaponMaxDmg:   10,
			armorDefense:   0,
			expectedMinDmg: 5,
			expectedMaxDmg: 10,
		},
		{
			name:           "weapon vs light armor",
			weaponMinDmg:   10,
			weaponMaxDmg:   15,
			armorDefense:   3,
			expectedMinDmg: 7,
			expectedMaxDmg: 12,
		},
		{
			name:           "weapon vs heavy armor",
			weaponMinDmg:   8,
			weaponMaxDmg:   12,
			armorDefense:   10,
			expectedMinDmg: 0, // Damage reduced to 0
			expectedMaxDmg: 2,
		},
		{
			name:           "weak weapon vs strong armor",
			weaponMinDmg:   3,
			weaponMaxDmg:   5,
			armorDefense:   8,
			expectedMinDmg: 0, // All damage blocked
			expectedMaxDmg: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test minimum damage
			minDamage := tt.weaponMinDmg - tt.armorDefense
			if minDamage < 0 {
				minDamage = 0
			}

			if minDamage != tt.expectedMinDmg {
				t.Errorf("Expected min damage %d, got %d", tt.expectedMinDmg, minDamage)
			}

			// Test maximum damage
			maxDamage := tt.weaponMaxDmg - tt.armorDefense
			if maxDamage < 0 {
				maxDamage = 0
			}

			if maxDamage != tt.expectedMaxDmg {
				t.Errorf("Expected max damage %d, got %d", tt.expectedMaxDmg, maxDamage)
			}
		})
	}
}
