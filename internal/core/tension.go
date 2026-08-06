package core

import "errors"

type WinchState struct {
	Tension    float64
	CrateDepth float64
	ItemWeight float64
}

// UpdateTick returns true if successfully extracted (depth <= 0).
// Returns an error if snapped.
func UpdateTick(ws *WinchState, player *Player, spacebarPressed bool, dt float64) (bool, error) {
	if player.EquippedWinch == nil {
		return false, errors.New("no winch equipped")
	}

	// Reduced physics speed to make the needle float lazily (~5 second buffer)
	gravity := 10.0 + (ws.ItemWeight * 0.1)
	pullStrength := player.EquippedWinch.PullSpeed * 0.5

	// Gravity ALWAYS applies. Spacebar fights against it.
	netForce := -gravity
	if spacebarPressed {
		netForce += pullStrength
	}

	ws.Tension += netForce * dt

	if ws.Tension < 0 || ws.Tension > 100 {
		return false, errors.New("CABLE_SNAPPED")
	}

	// Calculate Dynamic Green Zone based on baseline + durability decay
	baseSize := player.EquippedWinch.GreenZoneSize
	healthFactor := player.WinchHealth / 100.0 // 0.0 to 1.0
	actualSize := baseSize * healthFactor

	// Center around 50
	halfSize := actualSize / 2.0
	greenZoneMin := 50.0 - halfSize
	greenZoneMax := 50.0 + halfSize

	// Must have at least a tiny window (e.g. 5 units) so it's not impossible
	if greenZoneMax-greenZoneMin < 5.0 {
		greenZoneMin = 47.5
		greenZoneMax = 52.5
	}

	if ws.Tension >= greenZoneMin && ws.Tension <= greenZoneMax {
		// Base depth reduction speed (e.g. 50 depth per sec)
		ws.CrateDepth -= 50.0 * dt
	}

	if ws.CrateDepth <= 0 {
		ws.CrateDepth = 0
		return true, nil
	}

	return false, nil
}

// GetGreenZone helper for rendering UI
func GetGreenZone(player *Player) (float64, float64) {
	if player.EquippedWinch == nil {
		return 40, 60
	}
	baseSize := player.EquippedWinch.GreenZoneSize
	healthFactor := player.WinchHealth / 100.0
	actualSize := baseSize * healthFactor
	halfSize := actualSize / 2.0
	min := 50.0 - halfSize
	max := 50.0 + halfSize

	if max-min < 5.0 {
		min = 47.5
		max = 52.5
	}
	return min, max
}
