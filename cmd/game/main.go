package main

import (
	"log"

	"pressure-drop/internal/core"
	"pressure-drop/internal/engine"
	"pressure-drop/internal/states"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	ebiten.SetWindowSize(800, 600)
	ebiten.SetWindowTitle("Pressure Drop")

	reg, err := engine.LoadRegistry()
	if err != nil {
		log.Fatalf("Failed to load registry (Data-Driven JSON missing): %v", err)
	}

	player := &core.Player{
		Money:         100, // Starting allowance
		WinchHealth:   100,
		EquippedWinch: &reg.Upgrades[0], // Give them the Rusty Winch for free immediately
	}

	// Wait, we need the initial menu state to bootstrap the manager
	manager := states.NewStateManager(nil)
	menu := states.NewMenuState(manager)
	manager.SwitchState(menu)

	// In the next tick it will switch perfectly
	g := engine.NewGame(player, reg, manager)

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
