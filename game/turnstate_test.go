package game

import "testing"

func TestTurnStateConstants(t *testing.T) {
	if WaitingForPlayerInput != 0 {
		t.Errorf("WaitingForPlayerInput = %d, expected 0", WaitingForPlayerInput)
	}
	if ProcessingPlayerAction != 1 {
		t.Errorf("ProcessingPlayerAction = %d, expected 1", ProcessingPlayerAction)
	}
	if ProcessingMonsterTurn != 2 {
		t.Errorf("ProcessingMonsterTurn = %d, expected 2", ProcessingMonsterTurn)
	}
	if GameOver != 3 {
		t.Errorf("GameOver = %d, expected 3", GameOver)
	}
}

func TestGetNextState(t *testing.T) {
	tests := []struct {
		name     string
		current  TurnState
		expected TurnState
	}{
		{
			name:     "WaitingForPlayerInput to ProcessingPlayerAction",
			current:  WaitingForPlayerInput,
			expected: ProcessingPlayerAction,
		},
		{
			name:     "ProcessingPlayerAction to ProcessingMonsterTurn",
			current:  ProcessingPlayerAction,
			expected: ProcessingMonsterTurn,
		},
		{
			name:     "ProcessingMonsterTurn to WaitingForPlayerInput",
			current:  ProcessingMonsterTurn,
			expected: WaitingForPlayerInput,
		},
		{
			name:     "GameOver stays GameOver",
			current:  GameOver,
			expected: GameOver,
		},
		{
			name:     "Invalid state defaults to ProcessingPlayerAction",
			current:  TurnState(999),
			expected: ProcessingPlayerAction,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetNextState(tt.current)
			if result != tt.expected {
				t.Errorf("GetNextState(%d) = %d, expected %d", tt.current, result, tt.expected)
			}
		})
	}
}

func TestTurnStateCycle(t *testing.T) {
	var state TurnState = WaitingForPlayerInput

	// Test a complete cycle
	state = GetNextState(state) // Should be ProcessingPlayerAction
	if state != ProcessingPlayerAction {
		t.Errorf("First transition failed: got %d, expected %d", state, ProcessingPlayerAction)
	}

	state = GetNextState(state) // Should be ProcessingMonsterTurn
	if state != ProcessingMonsterTurn {
		t.Errorf("Second transition failed: got %d, expected %d", state, ProcessingMonsterTurn)
	}

	state = GetNextState(state) // Should be WaitingForPlayerInput
	if state != WaitingForPlayerInput {
		t.Errorf("Third transition failed: got %d, expected %d", state, WaitingForPlayerInput)
	}
}
