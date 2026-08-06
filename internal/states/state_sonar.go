package states

import (
	"fmt"
	"image/color"
	"math/rand"
	"time"

	"pressure-drop/internal/core"
	"pressure-drop/internal/engine"
	"pressure-drop/internal/render"

	"github.com/hajimehoshi/ebiten/v2"
)

type lootDot struct {
	P    core.Point
	Item core.Item
}

type SonarState struct {
	manager *StateManager
	dots    []lootDot
}

func NewSonarState(manager *StateManager) *SonarState {
	s := &SonarState{manager: manager}
	rand.Seed(time.Now().UnixNano())
	return s
}

func (s *SonarState) Update(input *engine.InputManager, player *core.Player, reg *engine.Registry) error {
	if len(s.dots) == 0 {
		// Spawn 3 random ones attached to random registry items
		for i := 0; i < 3; i++ {
			itm := reg.Items[rand.Intn(len(reg.Items))]
			s.dots = append(s.dots, lootDot{
				P: core.Point{
					X: 100 + rand.Float64()*600,
					Y: 100 + rand.Float64()*400,
				},
				Item: itm,
			})
		}
	}

	if input.IsKeyJustPressed(ebiten.KeyEscape) {
		s.manager.SwitchState(NewMenuState(s.manager))
		return nil
	}

	mx, my := input.MousePosition()
	mp := core.Point{X: float64(mx), Y: float64(my)}

	// Find if clicked on any
	if input.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		for _, d := range s.dots {
			dist := core.CalculateDistance(mp, d.P)
			intensity := core.GetEchoIntensity(dist, d.Item.Weight)
			if intensity > 0.95 {
				// HOOKED! Switch to winch state
				s.manager.SwitchState(NewWinchState(s.manager, d.Item))
				return nil
			}
		}
	}

	return nil
}

func (s *SonarState) Draw(screen *ebiten.Image, player *core.Player) {
	render.DrawText(screen, 30, 30, "=== SONAR TRACKING ===")
	render.DrawText(screen, 30, 50, "Move your mouse blindly around the dark window.")
	render.DrawText(screen, 30, 70, "The Echo Intensity rises as you get closer to hidden loot.")
	render.DrawText(screen, 30, 90, "When Intensity is > 0.95, LEFT CLICK to hook the crate!")
	render.DrawText(screen, 30, 120, "[ESC] Abort Dive")

	mx, my := ebiten.CursorPosition()
	mp := core.Point{X: float64(mx), Y: float64(my)}

	// Logic for UI drawing echo intensity
	// NOTE: Because physics is state-agnostic, we borrow it here loosely just for UI reading
	// Wait, we need the heaviest intensity since multiple can overlap
	maxIntensity := 0.0
	for _, d := range s.dots {
		// Just to debug visually, let's draw them faintly or invisibly.
		// For now we don't draw dots, so it's a real sonar test.
		dist := core.CalculateDistance(mp, d.P)
		intensity := core.GetEchoIntensity(dist, d.Item.Weight)
		if intensity > maxIntensity {
			maxIntensity = intensity
		}
	}

	render.DrawText(screen, 350, 300, "%s", fmt.Sprintf("ECHO INTENSITY: %.2f", maxIntensity))

	// Draw a circle roughly around mouse based on max intensity
	cSize := float64(maxIntensity * 50)
	if cSize > 0 {
		render.DrawRect(screen, float64(mx)-cSize/2, float64(my)-cSize/2, cSize, cSize, color.RGBA{0, 255, 0, 100})
	}
}
