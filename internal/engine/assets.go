package engine

import (
	"encoding/json"
	"fmt"
	"os"

	"pressure-drop/internal/core"
)

type Registry struct {
	Items    []core.Item
	Upgrades []core.WinchUpgrade
}

func LoadRegistry() (*Registry, error) {
	reg := &Registry{}

	itemData, err := os.ReadFile("assets/data/items.json")
	if err != nil {
		return nil, fmt.Errorf("failed resolving items.json: %v", err)
	}
	if err := json.Unmarshal(itemData, &reg.Items); err != nil {
		return nil, fmt.Errorf("failed parsing items.json: %v", err)
	}

	upgradeData, err := os.ReadFile("assets/data/upgrades.json")
	if err != nil {
		return nil, fmt.Errorf("failed resolving upgrades.json: %v", err)
	}
	if err := json.Unmarshal(upgradeData, &reg.Upgrades); err != nil {
		return nil, fmt.Errorf("failed parsing upgrades.json: %v", err)
	}

	return reg, nil
}
