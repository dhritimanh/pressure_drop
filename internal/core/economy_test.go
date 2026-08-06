package core
import "testing"

func TestRepairWinch(t *testing.T) {
	upgrade := &WinchUpgrade{RepairCostPerPoint: 10}
	player := &Player{Money: 500, WinchHealth: 50, EquippedWinch: upgrade}
	
	// Test successful repair
	err := RepairWinch(player, 20) 
	if err != nil || player.Money != 300 || player.WinchHealth != 70 {
		t.Errorf("Repair failed logic test")
	}
	
	// Test insufficient funds
	err = RepairWinch(player, 50)
	if err == nil {
		t.Errorf("Expected failure due to lack of funds, but succeeded")
	}

	// Test cap at 100
	player.Money = 9999
	_ = RepairWinch(player, 500)
	if player.WinchHealth > 100 {
		t.Errorf("Health overcapped to %v", player.WinchHealth)
	}
}
