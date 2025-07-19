package main

type TurnState int

const (
	WaitingForPlayerInput TurnState = iota
	ProcessingPlayerAction
	ProcessingMonsterTurn
	GameOver
)

func GetNextState(state TurnState) TurnState {
	switch state {
	case WaitingForPlayerInput:
		return ProcessingPlayerAction
	case ProcessingPlayerAction:
		return ProcessingMonsterTurn
	case ProcessingMonsterTurn:
		return WaitingForPlayerInput
	case GameOver:
		return GameOver
	default:
		return ProcessingPlayerAction
	}
}
