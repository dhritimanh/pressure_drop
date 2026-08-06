package core

import (
	"math"
)

type Point struct {
	X float64
	Y float64
}

func CalculateDistance(p1, p2 Point) float64 {
	dx := p1.X - p2.X
	dy := p1.Y - p2.Y
	return math.Sqrt(dx*dx + dy*dy)
}

func GetEchoIntensity(distance float64, itemWeight float64) float64 {
	// Base range. Heavy items have a larger detection radius.
	maxRange := 150.0 + itemWeight

	if distance > maxRange {
		return 0.0
	}

	// 1.0 (very close) to 0.0 (far)
	intensity := 1.0 - (distance / maxRange)
	
	// Ensure it strictly bounds between 0 and 1
	if intensity < 0 {
		return 0
	}
	if intensity > 1 {
		return 1
	}

	return intensity
}
