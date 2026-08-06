package states

import (
	"fmt"
	"image/color"

	"pressure-drop/internal/core"
	"pressure-drop/internal/engine"
	"pressure-drop/internal/render"

	"github.com/hajimehoshi/ebiten/v2"
)

type WinchState struct {
	manager    *StateManager
	item       core.Item
	winchState *core.WinchState
	timer      float64
	doneMsg    string
}

func NewWinchState(manager *StateManager, item core.Item) *WinchState {
	return &WinchState{
		manager: manager,
		item:    item,
		winchState: &core.WinchState{
			Tension:    50,
			CrateDepth: 1000,
			ItemWeight: item.Weight,
		},
		timer: 45.0,
	}
}

func (s *WinchState) Update(input *engine.InputManager, player *core.Player, reg *engine.Registry) error {
	dt := 1.0 / 60.0

	s.timer -= dt

	if s.timer <= 0 {
		player.Money -= 1000
		player.Inventory = []core.Item{}
		s.doneMsg = "BUSTED! Coast Guard caught you. -$1000"
		s.returnToMenu(player)
		return nil
	}

	spacePressed := input.IsKeyPressed(ebiten.KeySpace)
	extracted, err := core.UpdateTick(s.winchState, player, spacePressed, dt)

	if err != nil {
		s.doneMsg = "CABLE SNAPPED!"
		core.ApplyDiveDegradation(player)
		s.returnToMenu(player)
		return nil
	}

	if extracted {
		player.Inventory = append(player.Inventory, s.item)
		core.ApplyDiveDegradation(player)
		s.doneMsg = "EXTRACTED: " + s.item.Name
		s.returnToMenu(player)
		return nil
	}

	return nil
}

func (s *WinchState) returnToMenu(player *core.Player) {
	menu := NewMenuState(s.manager)
	menu.msg = s.doneMsg
	s.manager.SwitchState(menu)
}

func (s *WinchState) Draw(screen *ebiten.Image, player *core.Player) {
	render.DrawText(screen, 30, 30, "=== EXTRACTION PHASE ===")
	render.DrawText(screen, 30, 50, "OBJECTIVE: Keep the RED line inside the GREEN Box!")
	render.DrawText(screen, 30, 70, "- [HOLD SPACE] = Pull harder (Red line goes UP)")
	render.DrawText(screen, 30, 90, "- [LET GO] = Let gravity pull (Red line goes DOWN)")
	render.DrawText(screen, 30, 110, "WARNING: Touching completely top or bottom SNAPS THE CABLE.")

	render.DrawText(screen, 30, 150, "%s", fmt.Sprintf("Depth Remaining: %.1f meters", s.winchState.CrateDepth))
	render.DrawText(screen, 30, 170, "%s", fmt.Sprintf("Coast Guard ETA: %.1f seconds", s.timer))

	// Draw Green Zone and Tension indicator
	minZone, maxZone := core.GetGreenZone(player)

	// Draw Background Meter (100x400)
	barX := float64(400)
	barY := float64(100)
	barH := float64(400)
	barW := float64(30)
	render.DrawRect(screen, barX, barY, barW, barH, color.RGBA{100, 100, 100, 255})

	// Draw Green Zone mapped (invert Y because 0 is bottom in UI, 100 is top)
	gzTop := barY + barH - (maxZone/100.0)*barH
	gzBottom := barY + barH - (minZone/100.0)*barH
	render.DrawRect(screen, barX, gzTop, barW, gzBottom-gzTop, color.RGBA{0, 255, 0, 150})

	// Draw Tension marker
	tensionY := barY + barH - (s.winchState.Tension/100.0)*barH
	render.DrawRect(screen, barX-10, tensionY-5, barW+20, 10, color.RGBA{255, 0, 0, 255})
}
