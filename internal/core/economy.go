package core

import (
	"errors"
	"math"
)

type Item struct {
	ID     string
	Name   string
	Price  int
	Weight float64
}

type WinchUpgrade struct {
	ID                 string
	Name               string
	Price              int
	PullSpeed          float64
	DegradationPerDive float64
	RepairCostPerPoint int
	GreenZoneSize      float64
}

type Player struct {
	Money         int
	WinchHealth   float64
	EquippedWinch *WinchUpgrade
	Inventory     []Item
}

func RepairWinch(p *Player, pointsToRepair float64) error {
	if p.EquippedWinch == nil {
		return errors.New("no winch equipped")
	}

	cost := int(math.Ceil(pointsToRepair)) * p.EquippedWinch.RepairCostPerPoint
	if p.Money < cost {
		return errors.New("insufficient funds for repair")
	}

	p.Money -= cost
	p.WinchHealth += pointsToRepair
	if p.WinchHealth > 100 {
		p.WinchHealth = 100
	}
	return nil
}

func ApplyDiveDegradation(p *Player) {
	if p.EquippedWinch == nil {
		return
	}
	p.WinchHealth -= p.EquippedWinch.DegradationPerDive
	if p.WinchHealth < 0 {
		p.WinchHealth = 0
	}
}

func SellAll(p *Player) int {
	totalEarned := 0
	for _, item := range p.Inventory {
		totalEarned += item.Price
	}
	p.Money += totalEarned
	p.Inventory = []Item{}
	return totalEarned
}

func BuyUpgrade(p *Player, upgrade *WinchUpgrade) error {
    if p.Money < upgrade.Price {
        return errors.New("insufficient funds for upgrade")
    }
    p.Money -= upgrade.Price
    p.EquippedWinch = upgrade
    return nil
}
