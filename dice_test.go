package main

import "testing"

func TestGetRandomBetween(t *testing.T) {
	// Test the function multiple times to check range
	tests := []struct {
		name string
		low  int
		high int
	}{
		{
			name: "simple range 1-5",
			low:  1,
			high: 5,
		},
		{
			name: "range 3-7",
			low:  3,
			high: 7,
		},
		{
			name: "single number range",
			low:  5,
			high: 5,
		},
		{
			name: "large range 10-20",
			low:  10,
			high: 20,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test multiple iterations to check range bounds
			for i := 0; i < 100; i++ {
				result := GetRandomBetween(tt.low, tt.high)

				if result < tt.low {
					t.Errorf("GetRandomBetween(%d, %d) returned %d, which is below minimum %d",
						tt.low, tt.high, result, tt.low)
				}

				if result > tt.high {
					t.Errorf("GetRandomBetween(%d, %d) returned %d, which is above maximum %d",
						tt.low, tt.high, result, tt.high)
				}
			}
		})
	}
}

func TestGetDiceRoll(t *testing.T) {
	tests := []struct {
		name string
		num  int
	}{
		{
			name: "d6 roll",
			num:  6,
		},
		{
			name: "d20 roll",
			num:  20,
		},
		{
			name: "single sided die",
			num:  1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test multiple rolls to verify range
			for i := 0; i < 100; i++ {
				result := GetDiceRoll(tt.num)

				if result < 1 {
					t.Errorf("GetDiceRoll(%d) returned %d, which is below minimum 1", tt.num, result)
				}

				if result > tt.num {
					t.Errorf("GetDiceRoll(%d) returned %d, which is above maximum %d", tt.num, result, tt.num)
				}
			}
		})
	}
}

func TestGetRandomInt(t *testing.T) {
	tests := []struct {
		name string
		num  int
	}{
		{
			name: "0 to 5",
			num:  6,
		},
		{
			name: "0 to 9",
			num:  10,
		},
		{
			name: "single value",
			num:  1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test multiple iterations
			for i := 0; i < 100; i++ {
				result := GetRandomInt(tt.num)

				if result < 0 {
					t.Errorf("GetRandomInt(%d) returned %d, which is below minimum 0", tt.num, result)
				}

				if result >= tt.num {
					t.Errorf("GetRandomInt(%d) returned %d, which should be less than %d", tt.num, result, tt.num)
				}
			}
		})
	}
}

// Test the specific bug case
func TestGetRandomBetweenSpecificCase(t *testing.T) {
	// Test the case mentioned in the bug report: GetRandomBetween(3, 7)
	// Should return values in [3, 4, 5, 6, 7]

	validValues := map[int]bool{3: true, 4: true, 5: true, 6: true, 7: true}

	for i := 0; i < 1000; i++ {
		result := GetRandomBetween(3, 7)

		if !validValues[result] {
			t.Errorf("GetRandomBetween(3, 7) returned invalid value %d, expected one of [3,4,5,6,7]", result)
		}
	}
}
