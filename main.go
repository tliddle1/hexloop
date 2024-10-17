package main

import (
	"image/color"
	"log"
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	vector2 "github.com/tliddle1/game/vector"
)

const (
	screenWidth           = 800
	screenHeight          = 600
	hexVertexRadius       = 30 // Distance from center of hexagon to vertex
	rows                  = 5  // Number of hexagon rows
	cols                  = 24 // Number of hexagon columns
	numHexagonSides       = 6
	defaultStrokeWidth    = 2
	connectionStrokeWidth = 3
)

var (
	backgroundColor        = color.RGBA{R: 251, G: 217, B: 100, A: 255} // Yellow
	hexBorderColor         = color.RGBA{R: 254, G: 237, B: 161, A: 255} // Beige
	connectionColor        = color.RGBA{R: 133, G: 77, B: 13, A: 255}   // Brown
	pendingConnectionColor = color.RGBA{R: 150, G: 123, B: 182, A: 255} // Lavender
	hexSideRadius          = math.Sqrt(3) / 2 * hexVertexRadius
	possibleConnections    = [][]Connection{
		{{0, 1}, {2, 3}, {4, 5}},
		{{0, 1}, {2, 4}, {3, 5}},
		{{0, 1}, {2, 5}, {3, 4}},
		{{0, 2}, {1, 3}, {4, 5}},
		{{0, 2}, {1, 4}, {3, 5}},
		{{0, 2}, {1, 5}, {3, 4}},
		{{0, 3}, {2, 1}, {4, 5}},
		{{0, 3}, {2, 4}, {1, 5}},
		{{0, 3}, {2, 5}, {1, 4}},
		{{0, 4}, {2, 3}, {1, 5}},
		{{0, 4}, {2, 1}, {3, 5}},
		{{0, 4}, {2, 5}, {3, 1}},
		{{0, 5}, {2, 3}, {4, 1}},
		{{0, 5}, {2, 4}, {3, 1}},
		{{0, 5}, {2, 1}, {3, 4}},
	}
	possibleConnectionIndex = rand.Intn(len(possibleConnections))
)

// Hex represents a hexagonal tile
type Hex struct {
	Q, R              int     // Column and Row
	X, Y              float64 // Center of the hex
	Connections       []Connection
	PendingConnection bool
}

type Connection [2]int
type coordinate [2]float32

// Game represents the game state
type Game struct {
	Hexes []Hex // List of hexagons
}

// NewGame initializes the game state
func NewGame() *Game {
	var hexes []Hex
	for rowNum := 0; rowNum < rows; rowNum++ {
		for colNum := 0; colNum < cols; colNum++ {
			// Calculate x, y positions of each hex using axial coordinates
			x := float64(colNum) * hexSideRadius
			y := float64(rowNum) * hexVertexRadius * 3
			// Offset odd rows to create staggered effect
			if colNum%2 != 0 {
				y += hexVertexRadius * 1.5
			}
			hexes = append(hexes, Hex{Q: colNum, R: rowNum, X: x + 100, Y: y + 100})
		}
	}
	return &Game{Hexes: hexes}
}

// Update handles game logic updates
func (this *Game) Update() error {
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		mouseX, mouseY := ebiten.CursorPosition()
		for i, hex := range this.Hexes {
			if pointInHexagon(float64(mouseX), float64(mouseY), hex.X, hex.Y, hexVertexRadius) {
				if len(hex.Connections) == 0 {
					this.Hexes[i].Connections = possibleConnections[possibleConnectionIndex]
					possibleConnectionIndex = rand.Intn(len(possibleConnections))
				}
			}
		}
	} else {
		mouseX, mouseY := ebiten.CursorPosition()
		for i, hex := range this.Hexes {
			if pointInHexagon(float64(mouseX), float64(mouseY), hex.X, hex.Y, hexVertexRadius) {
				if len(hex.Connections) == 0 {
					this.Hexes[i].PendingConnection = true
				}
			} else {
				this.Hexes[i].PendingConnection = false
			}
		}
	}
	return nil
}

// Draw renders the game state
func (this *Game) Draw(screen *ebiten.Image) {
	screen.Fill(backgroundColor)
	// Draw all hexagons
	for _, hex := range this.Hexes {
		drawHexagon(screen, hex.X, hex.Y)
	}
	for _, hex := range this.Hexes {
		drawHexagonConnections(screen, hex)
		if hex.PendingConnection {
			drawPendingHexagonConnection(screen, hex)
		}
	}
}

// Layout sets the screen size
func (this *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	//return screenWidth, screenHeight
	return outsideWidth, outsideHeight
}

// drawHexagon draws a single hexagon
func drawHexagon(screen *ebiten.Image, centerX, centerY float64) {
	drawHexagonSides(screen, getHexagonVertexCoordinates(centerX, centerY))
}

func drawHexagonSides(screen *ebiten.Image, vertices []coordinate) {
	vertices = append(vertices, vertices[0])
	for i := 0; i < numHexagonSides; i++ {
		vector.StrokeLine(screen,
			vertices[i][0],
			vertices[i][1],
			vertices[i+1][0],
			vertices[i+1][1],
			defaultStrokeWidth,
			hexBorderColor,
			false)
	}
}

func getHexagonVertexCoordinates(centerX, centerY float64) []coordinate {
	vertices := make([]coordinate, numHexagonSides)
	for i := 0; i < numHexagonSides; i++ {
		angle := math.Pi/3*float64(i) - math.Pi/6
		x := getXCoordinateFromPolar(centerX, hexVertexRadius, angle)
		y := getYCoordinateFromPolar(centerY, hexVertexRadius, angle)
		vertices[i] = coordinate{float32(x), float32(y)}
	}
	return vertices
}

func getHexagonSideCoordinates(centerX, centerY float64) []coordinate {
	sides := make([]coordinate, numHexagonSides)
	for i := 0; i < numHexagonSides; i++ {
		angle := math.Pi / 3 * float64(i-1)
		x := getXCoordinateFromPolar(centerX, hexSideRadius, angle)
		y := getYCoordinateFromPolar(centerY, hexSideRadius, angle)
		sides[i] = coordinate{float32(x), float32(y)}
	}
	return sides
}

func getXCoordinateFromPolar(centerX, radius, angle float64) float64 {
	return centerX + radius*math.Cos(angle)
}

func getYCoordinateFromPolar(centerY, radius, angle float64) float64 {
	return centerY + radius*math.Sin(angle)
}

func drawHexagonConnections(screen *ebiten.Image, hex Hex) {
	if len(hex.Connections) == 0 {
		return
	}
	for _, connection := range hex.Connections {
		drawHexagonConnection(screen, connection, hex, connectionColor)
	}
}

func drawHexagonConnection(screen *ebiten.Image, connection Connection, hex Hex, color color.RGBA) {
	sideA := connection[0]
	sideB := connection[1]
	diff := math.Abs(float64(sideA - sideB))
	// Straight Across
	if diff == 3 {
		sides := getHexagonSideCoordinates(hex.X, hex.Y)
		x1 := sides[sideA][0]
		y1 := sides[sideA][1]
		x2 := sides[sideB][0]
		y2 := sides[sideB][1]
		vector2.StrokeLine(screen, x1, y1, x2, y2, connectionStrokeWidth, color, false)
		return
	}
	// Large Curve
	if diff == 2 {
		centerSide := (sideA + sideB) / 2
		if centerSide == 1 {
			angleToSide := math.Pi
			x := hex.X - math.Cos(angleToSide)*hexSideRadius*2
			y := hex.Y - math.Sin(angleToSide)*hexSideRadius*2
			vector2.StrokePartialCircle(screen, float32(x), float32(y), float32(hexVertexRadius+hexSideRadius/2+defaultStrokeWidth), float32(angleToSide-math.Pi/6), float32(angleToSide+math.Pi/6), connectionStrokeWidth, color, false)
		}
		if centerSide == 2 {
			angleToSide := math.Pi * 4 / 3
			x := hex.X - math.Cos(angleToSide)*hexSideRadius*2
			y := hex.Y - math.Sin(angleToSide)*hexSideRadius*2
			vector2.StrokePartialCircle(screen, float32(x), float32(y), float32(hexVertexRadius+hexSideRadius/2+defaultStrokeWidth), float32(angleToSide-math.Pi/6), float32(angleToSide+math.Pi/6), connectionStrokeWidth, color, false)
		}
		if centerSide == 3 {
			angleToSide := math.Pi * 5 / 3
			x := hex.X - math.Cos(angleToSide)*hexSideRadius*2
			y := hex.Y - math.Sin(angleToSide)*hexSideRadius*2
			vector2.StrokePartialCircle(screen, float32(x), float32(y), float32(hexVertexRadius+hexSideRadius/2+defaultStrokeWidth), float32(angleToSide-math.Pi/6), float32(angleToSide+math.Pi/6), connectionStrokeWidth, color, false)
		}
		if centerSide == 4 {
			angleToSide := 0.0
			x := hex.X - math.Cos(angleToSide)*hexSideRadius*2
			y := hex.Y - math.Sin(angleToSide)*hexSideRadius*2
			vector2.StrokePartialCircle(screen, float32(x), float32(y), float32(hexVertexRadius+hexSideRadius/2+defaultStrokeWidth), float32(angleToSide-math.Pi/6), float32(angleToSide+math.Pi/6), connectionStrokeWidth, color, false)
		}
		return
	}
	if diff == 4 {
		centerSide := (sideA + sideB) / 2
		if centerSide == 2 {
			angleToSide := math.Pi * 1 / 3
			x := hex.X - math.Cos(angleToSide)*hexSideRadius*2
			y := hex.Y - math.Sin(angleToSide)*hexSideRadius*2
			vector2.StrokePartialCircle(screen, float32(x), float32(y), float32(hexVertexRadius+hexSideRadius/2+defaultStrokeWidth), float32(angleToSide-math.Pi/6), float32(angleToSide+math.Pi/6), connectionStrokeWidth, color, false)
			return
		} else if centerSide == 3 {
			angleToSide := math.Pi * 2 / 3
			x := hex.X - math.Cos(angleToSide)*hexSideRadius*2
			y := hex.Y - math.Sin(angleToSide)*hexSideRadius*2
			vector2.StrokePartialCircle(screen, float32(x), float32(y), float32(hexVertexRadius+hexSideRadius/2+defaultStrokeWidth), float32(angleToSide-math.Pi/6), float32(angleToSide+math.Pi/6), connectionStrokeWidth, color, false)
			return
		}
	}
	// Small Curve
	vertices := getHexagonVertexCoordinates(hex.X, hex.Y)
	if diff == 1 {
		drawSmallCurve(screen, vertices, min(sideA, sideB), color)
		return
	}
	if diff == 5 {
		drawSmallCurve(screen, vertices, 5, color)
		return
	}
	panic("did I forget something?")
}

func drawSmallCurve(screen *ebiten.Image, vertices []coordinate, vertex int, color color.RGBA) {
	x := vertices[vertex][0]
	y := vertices[vertex][1]
	adjustor := math.Pi / 3 * float32(vertex)
	vector2.StrokePartialCircle(screen, x, y, hexVertexRadius/2, math.Pi/2+adjustor, -math.Pi*5/6+adjustor, connectionStrokeWidth, color, false)
}

func drawPendingHexagonConnection(screen *ebiten.Image, hex Hex) {
	if len(hex.Connections) != 0 {
		return
	}
	////drawHexagonConnection(screen, Connection{0, 2}, hex, pendingConnectionColor)
	//vertices := getHexagonVertexCoordinates(hex.X, hex.Y)
	////angle := math.Pi / 3
	//side := 0
	////vector.StrokeLine(screen, float32(hex.X), float32(hex.Y), float32(hex.X+math.Cos(angle)*hexSideRadius), float32(hex.Y+math.Sin(angle)*hexSideRadius), defaultStrokeWidth, pendingConnectionColor, false)
	//vector.StrokeLine(screen, float32(hex.X)+10, float32(hex.Y)+10, vertices[side][0], vertices[side][1], defaultStrokeWidth, pendingConnectionColor, false)
	//return
	for _, connection := range possibleConnections[possibleConnectionIndex] {
		drawHexagonConnection(screen, connection, hex, pendingConnectionColor)
	}
}

// pointInHexagon checks if a point is inside a hexagon
func pointInHexagon(px, py, cx, cy, radius float64) bool {
	buffer := .1
	dx := math.Abs(px-cx) / radius
	dy := math.Abs(py-cy) / radius
	return dx <= 1.0-buffer && dy <= math.Sqrt(3.0)/2.0-buffer && dx+dy/math.Sqrt(3.0) <= 1.0-buffer
}

func main() {
	game := NewGame()
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Clickable Hexagonal Grid")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
