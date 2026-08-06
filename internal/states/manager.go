package states

import (
	"pressure-drop/internal/core"
	"pressure-drop/internal/engine"

	"github.com/hajimehoshi/ebiten/v2"
)

type GameState interface {
	Update(input *engine.InputManager, player *core.Player, reg *engine.Registry) error
	Draw(screen *ebiten.Image, player *core.Player)
}

type StateManager struct {
	currentState GameState
	nextState    GameState
}

func NewStateManager(initial GameState) *StateManager {
	return &StateManager{
		currentState: initial,
	}
}

func (sm *StateManager) SwitchState(newState GameState) {
	sm.nextState = newState
}

func (sm *StateManager) Update(input *engine.InputManager, player *core.Player, reg *engine.Registry) error {
	if sm.nextState != nil {
		sm.currentState = sm.nextState
		sm.nextState = nil
	}

	if sm.currentState != nil {
		return sm.currentState.Update(input, player, reg)
	}
	return nil
}

func (sm *StateManager) Draw(screen *ebiten.Image, player *core.Player) {
	if sm.currentState != nil {
		sm.currentState.Draw(screen, player)
	}
}
