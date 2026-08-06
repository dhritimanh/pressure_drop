package engine

import (
	"image/color"

	"pressure-drop/internal/core"

	"github.com/hajimehoshi/ebiten/v2"
)

// StateManagerInterface to avoid circular dependencies between states/engine
type StateManagerInterface interface {
	Update(input *InputManager, player *core.Player, reg *Registry) error
	Draw(screen *ebiten.Image, player *core.Player)
}

type Game struct {
	Player  *core.Player
	Reg     *Registry
	Input   *InputManager
	Manager StateManagerInterface
}

func NewGame(player *core.Player, reg *Registry, sm StateManagerInterface) *Game {
	return &Game{
		Player:  player,
		Reg:     reg,
		Input:   &InputManager{},
		Manager: sm,
	}
}

func (g *Game) Update() error {
	return g.Manager.Update(g.Input, g.Player, g.Reg)
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{30, 40, 50, 255})
	g.Manager.Draw(screen, g.Player)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 800, 600
}
