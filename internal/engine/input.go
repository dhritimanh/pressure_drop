package engine

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// InputManager wraps ebiten input checks so states don't rely globally
type InputManager struct{}

func (i *InputManager) IsKeyJustPressed(key ebiten.Key) bool {
	return inpututil.IsKeyJustPressed(key)
}

func (i *InputManager) IsMouseButtonJustPressed(button ebiten.MouseButton) bool {
	return inpututil.IsMouseButtonJustPressed(button)
}

func (i *InputManager) MousePosition() (int, int) {
	return ebiten.CursorPosition()
}

func (i *InputManager) IsKeyPressed(key ebiten.Key) bool {
	return ebiten.IsKeyPressed(key)
}
