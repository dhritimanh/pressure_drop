package states

import (
	"fmt"

	"pressure-drop/internal/core"
	"pressure-drop/internal/engine"
	"pressure-drop/internal/render"

	"github.com/hajimehoshi/ebiten/v2"
)

type MenuState struct {
	manager *StateManager
	msg     string
}

func NewMenuState(manager *StateManager) *MenuState {
	return &MenuState{manager: manager, msg: "Welcome to Pressure Drop."}
}

func (s *MenuState) Update(input *engine.InputManager, player *core.Player, reg *engine.Registry) error {
	mx, my := input.MousePosition()
	clicked := input.IsMouseButtonJustPressed(ebiten.MouseButtonLeft)

	inBounds := func(y int) bool {
		return mx >= 30 && mx <= 300 && my >= y-15 && my <= y+5
	}

	if (clicked && inBounds(80)) || input.IsKeyJustPressed(ebiten.KeyS) {
		earned := core.SellAll(player)
		s.msg = fmt.Sprintf("Sold items for $%d", earned)
	}

	if (clicked && inBounds(100)) || input.IsKeyJustPressed(ebiten.KeyR) {
		if player.EquippedWinch != nil {
			missing := 100.0 - player.WinchHealth
			if missing > 0 {
				err := core.RepairWinch(player, missing)
				if err != nil {
					s.msg = "Cannot afford repair!"
				} else {
					s.msg = "Winch fully repaired."
				}
			} else {
				s.msg = "Winch is already at 100%."
			}
		} else {
			s.msg = "Buy a winch first!"
		}
	}

	if (clicked && inBounds(120)) || input.IsKeyJustPressed(ebiten.KeyU) {
		for i := range reg.Upgrades {
			if player.EquippedWinch == nil || player.EquippedWinch.Price < reg.Upgrades[i].Price {
				err := core.BuyUpgrade(player, &reg.Upgrades[i])
				if err != nil {
					s.msg = "Cannot afford upgrade!"
				} else {
					s.msg = fmt.Sprintf("Purchased %s!", reg.Upgrades[i].Name)
				}
				break
			}
		}
	}

	if (clicked && inBounds(150)) || input.IsKeyJustPressed(ebiten.KeyEnter) {
		if player.EquippedWinch == nil {
			s.msg = "You must buy a Winch first!"
		} else if player.WinchHealth <= 0 {
			s.msg = "Winch is completely broken. Repair it first!"
		} else {
			s.manager.SwitchState(NewSonarState(s.manager))
		}
	}

	if player.Money < 0 {
		s.msg = "GAME OVER - DEBTOR. The Coast Guard seized your vessel."
	}
	return nil
}

func (s *MenuState) Draw(screen *ebiten.Image, player *core.Player) {
	render.DrawText(screen, 30, 30, "=== PRESSURE DROP MENU ===")

	stats := fmt.Sprintf("Money: $%d | Winch Health: %.1f%%", player.Money, player.WinchHealth)
	if player.EquippedWinch != nil {
		stats += fmt.Sprintf(" | Equipped: %s", player.EquippedWinch.Name)
	} else {
		stats += " | NO WINCH EQUIPPED. BUY ONE!"
	}
	render.DrawText(screen, 30, 55, "%s", stats)

	render.DrawText(screen, 30, 80, "[S] Sell All Inventory")
	render.DrawText(screen, 30, 100, "[R] Repair Winch to 100%%")
	render.DrawText(screen, 30, 120, "[U] Buy Best Upgrade")
	render.DrawText(screen, 30, 150, "PRESS [ENTER] TO START DIVE")

	render.DrawTooltip(screen, s.msg)
}
