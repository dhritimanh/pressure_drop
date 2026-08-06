# Implementation Plan: Pressure Drop

This plan outlines the architecture and execution strategy to build "Pressure Drop", a 2D UI-driven deep-sea salvage game using Ebitengine. The implementation adheres strictly to a **Modular State Machine** and **Data-Driven Design** to ensure scalability, DLC-readiness, and testability.

## Time Estimate
As an AI agent, I can write the core infrastructure, state machine, and game logic sequentially within **3-4 prompt turns (approx. 20-30 minutes)**. For a skilled human Ebitengine developer, this project scope as defined would require about **3 to 5 hours** of focused effort.

## Visualizing Gaps & Incorporating the Business Economy

1. **The Business Economy (Wear & Tear):**
   - We will implement a deep strategic layer based on equipment tiers. Upgrades will not just be linear buffs.
   - **Performance:** dictating speed and accuracy.
   - **Durability (0-100%):** Degrades per dive. Lower durability actively degrades performance (e.g. winch green zone shrinks).
   - **Maintenance Cost:** Cost to repair 1% of durability.
   - **Tiers:** Early (Scrapper) gear loses durability fast but is cheap to fix. High-end (Corporate/Mil-Spec) gear loses durability slowly but costs a fortune to repair if damaged severely, introducing psychological gambling ("just one more dive...").

2. **Game Loop Consistency:**
   - The structured loop prevents state confusion: `Menu (Prep/Repair/Buy) -> Sonar (Scan) -> Winch (Extract) -> Menu (Resolve/Sell)`. 
   - Concluding an extraction (success or failure) always returns to the Menu to evaluate the financial toll.

3. **Loot Agency vs. Randomness:**
   - Rather than pure randomness, crate value will dictate sonar response and winch physics. This allows players to track heavier, more valuable artifacts at a much higher risk of cable snaps and durability loss.

## Proposed Architecture & File Structure

The project will strictly separate pure game logic from Ebitengine rendering/input.

### [NEW] `cmd/game/main.go`
- Entry point. Initializes Ebitengine window and injects dependencies to the `StateManager`.

### [NEW] `assets/data/`
- `items.json`: ALL crates/loot defined here (`id`, `name`, `price`, `weight`).
- `upgrades.json`: ALL upgrades defined here for DLC readiness (`id`, `name`, `price`, `pull_speed`, `degradation_per_dive`, `repair_cost_per_point`).

### [NEW] `internal/core/` (PURE GO - NO EBITEN)
- `economy.go`: Handles `Player` struct (Money, Equipment Health), repair math (`RepairWinch`), purchasing.
- `tension.go`: Winch mini-game physics math. Calculates tension drop/rise based on `Equipment.pull_speed` and `Equipment.degradation`.
- `spatial.go`: Pure distance/radar calculations.
- *Includes comprehensive unit tests (`economy_test.go`, `tension_test.go`) to guarantee bug-free logic.*

### [NEW] `internal/engine/` (EBITEN WRAPPERS)
- `audio.go`: Audio manager for `.wav` background loops and SFX.
- `input.go`: Mouse/Keyboard wrapping.

### [NEW] `internal/states/` (STATE MACHINE)
- `manager.go`: `GameState` interface and `StateManager` struct to prevent logic bleeding between active screens.
- `state_menu.go`: Handles buying, selling, and repairing.
- `state_sonar.go`: The tracking phase (audio pings, hovering over grid).
- `state_winch.go`: The extraction phase mini-game (spacebar tension control, 45s Coast Guard timer).

### [NEW] `internal/render/`
- `ui.go`: Abstractions for drawing buttons, text, and layouts.
- `radar.go`: Specific drawing logic for the sonar grid.

## Open Questions
- Do you want me to utilize the generic Ebiten debug text for rapid rendering ("Phase 1: The Skateboard"), or would you prefer I immediately implement TrueType fonts (TTF) and pixel-art shapes for a more polished aesthetic?
- We will be initializing the Go modules. What is the preferred module path? `github.com/raynwow/pressure-drop` or simply `pressure-drop`?

## Verification Plan
1. **Core Verification:** Run `go test ./internal/core/...` to instantly prove the economy and physics math is bug-free.
2. **Phase Validation:** Run `go run ./cmd/game/main.go` to simulate a full game loop, proving the state machine isolates logic correctly and JSON driven data populates the economy properly.
