package main

import (
	"fmt"
	"image/color"
	"log"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	screenWidth  = 800
	screenHeight = 600
	hexRadius    = 30 // Radius of the hexagons
	rows         = 5  // Number of hexagon rows
	cols         = 24 // Number of hexagon columns

)

// Hex represents a hexagonal tile

type Hex struct {
	Q, R int     // Axial coordinates
	X, Y float64 // Center of the hex
}

// Game represents the game state

type Game struct {
	Hexes []Hex // List of hexagons
}

// NewGame initializes the game state

func NewGame() *Game {
	hexes := []Hex{}
	for r := 0; r < rows; r++ {
		for q := 0; q < cols; q++ {
			// Calculate x, y positions of each hex using axial coordinates
			x := float64(q) * hexRadius / 2 * math.Sqrt(3.0)
			y := float64(r) * hexRadius * 3
			// Offset odd rows to create staggered effect
			if q%2 != 0 {
				y += hexRadius * 1.5
			}
			hexes = append(hexes, Hex{Q: q, R: r, X: x + 100, Y: y + 100})
		}
	}
	return &Game{Hexes: hexes}
}

// Update handles game logic updates

func (g *Game) Update() error {
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		// Get the mouse position
		mouseX, mouseY := ebiten.CursorPosition()
		for _, hex := range g.Hexes {
			if pointInHexagon(float64(mouseX), float64(mouseY), hex.X, hex.Y, hexRadius) {
				fmt.Printf("Hex clicked at Q: %d, R: %d\n", hex.Q, hex.R)
			}
		}
	}
	return nil
}

// Draw renders the game state

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{0, 0, 0, 255}) // Black background

	// Draw all hexagons
	for _, hex := range g.Hexes {
		drawHexagon(screen, hex.X, hex.Y, hexRadius)
	}
}

// Layout sets the screen size

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

// drawHexagon draws a single hexagon

func drawHexagon(screen *ebiten.Image, centerX, centerY, radius float64) {
	const sides = 6
	vertices := make([][2]float32, sides+1)
	for i := 0; i <= sides; i++ {
		angle := math.Pi/3*float64(i) - math.Pi/6
		x := centerX + radius*math.Cos(angle)
		y := centerY + radius*math.Sin(angle)
		vertices[i] = [2]float32{float32(x), float32(y)}
	}

	for i := 0; i < sides; i++ {
		vector.StrokeLine(screen, vertices[i][0], vertices[i][1], vertices[i+1][0], vertices[i+1][1], 2, color.White, false)
		//ebitenutil.DrawLine(screen, vertices[i][0], vertices[i][1], vertices[i+1][0], vertices[i+1][1], color.White)
	}
}

// pointInHexagon checks if a point is inside a hexagon
func pointInHexagon(px, py, cx, cy, radius float64) bool {
	dx := math.Abs(px-cx) / radius
	dy := math.Abs(py-cy) / radius
	return dx <= 1.0 && dy <= math.Sqrt(3.0)/2.0 && dx+dy/math.Sqrt(3.0) <= 1.0
}

func main() {
	game := NewGame()
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Clickable Hexagonal Grid")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
