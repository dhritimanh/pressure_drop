package render

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

func DrawText(screen *ebiten.Image, x, y int, format string, args ...interface{}) {
	text := fmt.Sprintf(format, args...)
	ebitenutil.DebugPrintAt(screen, text, x, y)
}

func DrawRect(screen *ebiten.Image, x, y, width, height float64, clr color.Color) {
	vector.DrawFilledRect(screen, float32(x), float32(y), float32(width), float32(height), clr, false)
}

func DrawTooltip(screen *ebiten.Image, msg string) {
	DrawRect(screen, 10, 560, 780, 30, color.RGBA{0, 0, 0, 200})
	DrawText(screen, 20, 565, msg)
}
