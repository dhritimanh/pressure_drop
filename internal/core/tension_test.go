package core

import "testing"

func TestTensionGreenZone(t *testing.T) {
	ws := &WinchState{
		Tension:    50,
		CrateDepth: 100,
		ItemWeight: 0,
	}

	upgrade := &WinchUpgrade{PullSpeed: 100, GreenZoneSize: 50}
	player := &Player{WinchHealth: 100, EquippedWinch: upgrade}

	// Stay perfectly at ~50. No spacebar, but gravity reduces it.
	// We'll just test that being at 50 reduces depth.
	_, err := UpdateTick(ws, player, false, 0.1) // Tension will drop, but starts at 50.
	
	if err != nil {
		t.Errorf("Unexpected cable snap")
	}
	if ws.CrateDepth >= 100 {
		t.Errorf("Depth should have decreased inside green zone")
	}

	// Force cable snap (less than 0)
	ws.Tension = 1.0
	_, err = UpdateTick(ws, player, false, 1.0)
	if err == nil || err.Error() != "CABLE_SNAPPED" {
		t.Errorf("Expected CABLE_SNAPPED from low tension")
	}
}
