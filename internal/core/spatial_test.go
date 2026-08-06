package core

import (
	"testing"
)

func TestCalculateDistance(t *testing.T) {
	p1 := Point{0, 0}
	p2 := Point{3, 4}
	dist := CalculateDistance(p1, p2)
	if dist != 5.0 {
		t.Errorf("Expected 5.0, got %v", dist)
	}

	distZero := CalculateDistance(p1, p1)
	if distZero != 0.0 {
		t.Errorf("Expected 0.0 distance for matching points")
	}
}

func TestGetEchoIntensity(t *testing.T) {
	// For weight 10, maxRange = 160.0
	intensityZero := GetEchoIntensity(0, 10)
	if intensityZero != 1.0 {
		t.Errorf("Expected intensity 1.0 at exact point, got %v", intensityZero)
	}

	intensityFar := GetEchoIntensity(2000, 10)
	if intensityFar != 0 {
		t.Errorf("Expected 0 intensity completely out of bounds, got %v", intensityFar)
	}
}
