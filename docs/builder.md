***

# MASTER AI BUILDER PROMPT: "PRESSURE DROP"

**Role:** You are an expert Golang Game Developer specializing in `github.com/hajimehoshi/ebiten/v2`. 
**Task:** We are building a 2D UI-driven deep-sea salvage game called "Pressure Drop". 
**Architecture Rules:** 
1. **Strict Separation of Concerns:** Pure game logic (math, economy, physics) must live in an `internal/core/` package and contain ZERO Ebitengine imports. Rendering and inputs live in `internal/engine/`, `internal/states/`, and `internal/render/`.
2. **Data-Driven:** Items and Upgrades must be unmarshaled from JSON files. No hardcoding entities.
3. **State Machine:** The game must use a strict State Machine interface (`Update`, `Draw`) to switch between screens (Menu, Sonar, Winch).

We will build this iteratively. Do not move to the next phase until I confirm the Visual Test passes.

---

## PHASE 1: The Skateboard (Window & State Machine)
**Goal:** Initialize the engine, set up the directory structure, and implement the State Manager. We need a window that switches between a "Main Menu" and a "Game" state when pressing Enter.

**AI Implementation Tasks:**
1. Generate `main.go` initializing an Ebiten window (800x600).
2. Create `internal/states/manager.go` defining the `GameState` interface (`Update`, `Draw`) and the `StateManager`.
3. Create `internal/states/state_menu.go` and `internal/states/state_sonar.go` with basic `ebitenutil.DebugPrint` text.
4. Implement input logic: Pressing 'Enter' in Menu swaps to Sonar. Pressing 'Escape' in Sonar swaps to Menu.

**Visual/Manual Test (My Eyes):** 
* I see an 800x600 window. 
* Top left text says "MENU STATE". 
* I press Enter, text changes to "SONAR STATE". 
* I press Escape, text changes to "MENU STATE".

**Edge Cases to Handle:**
* Rapidly pressing Enter/Escape should not cause state-stack overflows or panics.

---

## PHASE 2: The Cycle (The Sonar Grid & Math)
**Goal:** Build the "Tracking" mechanic. A grid where the player hovers the mouse to find hidden points via distance calculations.

**AI Implementation Tasks:**
1. In `internal/core/spatial.go`, write pure Go logic:
   * A `Point` struct (X, Y).
   * A function `CalculateDistance(p1, p2 Point) float64`.
   * A function `GetEchoIntensity(distance float64) float64` (closer = higher intensity).
2. Write **Unit Tests** for `spatial.go` testing exact distance math and zero-distance edge cases.
3. In `state_sonar.go`, generate 3 random hidden "Loot Points" on initialization.
4. Read mouse X/Y coordinates in Ebiten's `Update()`.
5. Draw a basic visual representation (e.g., a green circle centered on the screen) and print the "Echo Intensity" of the closest Loot Point as text based on mouse position.

**Visual/Manual Test (My Eyes):**
* I am in Sonar State. 
* I move my mouse around the blank screen.
* Text says `Echo Intensity: 0.1` and increases to `1.0` as I find the invisible X/Y coordinate.
* Once found, clicking it removes it from the array and prints "LOOT HOOKED".

**Edge Cases to Handle:**
* Mouse coordinates going outside the 800x600 window boundaries (handle gracefully, distance = max).
* Window resizing breaking coordinate mapping.

---

## PHASE 3: The Bike (The Winch Physics)
**Goal:** Build the "Extraction" mini-game. A tension bar that requires balancing using the Spacebar.

**AI Implementation Tasks:**
1. In `internal/core/tension.go`, write pure Go physics logic:
   * `WinchState` struct: `Tension (0-100)`, `CrateDepth (1000 to 0)`.
   * `UpdateTick(spacebarPressed bool)`: Gravity pulls Tension down. Spacebar pulls Tension up. CrateDepth decreases ONLY if Tension is between 30 and 70 (The "Green Zone"). 
   * If Tension < 0 or > 100, return `err == "CABLE_SNAPPED"`.
2. Write **Unit Tests** for `tension.go` ensuring tension doesn't float into negative numbers or over 100 without throwing the specific error.
3. Create `internal/states/state_winch.go`. 
4. Render 2 vertical rectangles: One background (the meter), one moving indicator (current tension). Draw a green box representing the 30-70 safe zone.
5. In `state_sonar.go`, when "LOOT HOOKED" is triggered, swap to `state_winch.go`.

**Visual/Manual Test (My Eyes):**
* I click the hidden loot in Sonar State. Screen swaps to Winch State.
* I see a vertical tension bar. It falls when I do nothing. It rises when I hold Space.
* Text shows `Depth: 1000m`. It only goes down when my tension indicator is inside the green zone.
* If it hits 0 or 100, text says "CABLE SNAPPED" and returns to Sonar State.
* If Depth hits 0, text says "EXTRACTED" and returns to Sonar State.

**Edge Cases to Handle:**
* Frame-rate dependence (Ensure tension math is multiplied by DeltaTime / TPS so it runs the same speed on 60hz vs 144hz monitors).

---

## PHASE 4: The Car (JSON Economy & Durability)
**Goal:** Introduce the Business logic, Loot Tables, Equipment Degradation, and UI Menus.

**AI Implementation Tasks:**
1. Create `assets/data/items.json` (3 dummy items with names, prices) and `upgrades.json` (2 dummy winches with stats).
2. In `internal/core/economy.go`, create `Player` struct (`Money`, `WinchDurability`).
3. Write pure Go functions `Repair(p *Player, amount float64) error` and `SellItem(p *Player, item Item)`.
4. Write **Unit Tests** for `economy.go` (Insufficient funds error, Durability capped at 100).
5. In `state_menu.go`, render the Player's Money, Winch Durability, and a "Repair (Cost: $X)" button using Ebiten mouse-click logic.

**Visual/Manual Test (My Eyes):**
* I start on the Menu State. Money is $0. Winch Durability is 100%.
* *Cheat key:* Press 'M' to add $500. Money updates to $500.
* I lower durability (Press 'D'). It goes to 80%.
* I click "Repair". Money drops, Durability goes to 100%. If I don't have enough money, click does nothing and logs an error to console.

**Edge Cases to Handle:**
* JSON parsing fails (missing file, bad syntax) -> game must panic with a clear error string on startup, not fail silently.
* Floating point currency issues (Ensure `Money` is an `int` acting as cents, or tightly controlled).

---

## PHASE 5: The Plane (The Full Integration Loop)
**Goal:** Tie Phases 1-4 together into the actual game loop. Apply penalties to the state machine.

**AI Implementation Tasks:**
1. Update `state_winch.go`: When returning to Sonar State, subtract 5% from `Player.WinchDurability`. 
2. Tie the Winch "Green Zone" to Durability. (e.g., If Durability is 50%, the Green Zone shrinks from 30-70 down to 45-55, making it harder).
3. Update `state_sonar.go`: When clicking a loot dot, randomly pull an item from the JSON loot table. Pass it to the Winch State. If successfully extracted, add it to `Player.Inventory`.
4. Update `state_menu.go`: Add a "Sell All" button that empties inventory and adds the JSON value to `Player.Money`.

**Visual/Manual Test (My Eyes):**
* Start in Menu -> Go to Sonar -> Track the dot -> Click it -> Go to Winch -> Successfully balance tension until Depth 0 -> Go back to Menu.
* In Menu, I see Durability dropped to 95%. I see "Item: Cartel Lockbox" in inventory. 
* I click "Sell All". Inventory empties, Money goes up by $15,000. 

**Edge Cases to Handle:**
* Attempting to start a dive (leaving Menu) with 0% Winch Durability should be blocked by the UI.

---

## PHASE 6: The Rocket (Audio & The Coast Guard Threat)
**Goal:** Add the atmospheric audio and the high-stakes timer threat.

**AI Implementation Tasks:**
1. Implement the `audio` package from Ebitengine in `internal/engine/audio.go`. 
2. Add functionality to play looping background `.wav`, and one-shot sound effects (Ping, Snap).
3. In `state_winch.go`, implement the **Coast Guard Timer**. Add a visual counter (Time left: 45s). 
4. If Time reaches 0 before Depth reaches 0, trigger "BUSTED". Remove all current items in `Player.Inventory`, deduct a $1,000 fine from `Player.Money`, and return to Menu.

**Visual/Manual Test (My Eyes):**
* Audio plays in the background.
* In Winch State, a timer counts down from 45 seconds. 
* If I purposely hold the crate too long and the timer hits 0, I am sent back to the Menu. My money is reduced by $1,000. 

**Edge Cases to Handle:**
* Player money going into negative numbers. (If Money < 0, render "GAME OVER - DEBTOR" text in the Menu state and block further gameplay).
* Audio context suspending when switching browser tabs (if compiling to WASM). Ebiten handles this mostly, but ensure audio resumes cleanly.