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
	connectionPermutations = [][]Connection{
		{{0, 1}, {2, 3}, {4, 5}},
		{{0, 1}, {2, 4}, {3, 5}},
		{{0, 1}, {2, 5}, {3, 4}},
		{{0, 2}, {1, 3}, {4, 5}},
		{{0, 2}, {1, 4}, {3, 5}},
		{{0, 2}, {1, 5}, {3, 4}},
		{{0, 3}, {1, 2}, {4, 5}},
		{{0, 3}, {1, 4}, {2, 5}},
		{{0, 3}, {1, 5}, {2, 4}},
		{{0, 4}, {1, 2}, {3, 5}},
		{{0, 4}, {1, 3}, {2, 5}},
		{{0, 4}, {1, 5}, {2, 3}},
		{{0, 5}, {1, 2}, {3, 4}},
		{{0, 5}, {1, 3}, {2, 4}},
		{{0, 5}, {1, 4}, {2, 3}},
	}
)

type Connection [2]int
type Coordinate [2]float64

// Game represents the game state
type Game struct {
	hexes                []*Hex // List of hexagons
	possibleConnections  [][]Connection
	nextConnectionsIndex int
}

// NewGame initializes the game state
func NewGame() *Game {
	return &Game{
		hexes:                newHexes(),
		possibleConnections:  connectionPermutations,
		nextConnectionsIndex: rand.Intn(len(connectionPermutations)),
	}
}

func newHexes() (hexes []*Hex) {
	for rowNum := 0; rowNum < rows; rowNum++ {
		for colNum := 0; colNum < cols; colNum++ {
			x := float64(colNum) * hexSideRadius
			y := float64(rowNum) * hexVertexRadius * 3
			// Offset odd rows to create staggered effect
			if colNum%2 != 0 {
				y += hexVertexRadius * 1.5
			}
			hexes = append(hexes, &Hex{col: colNum, row: rowNum, center: Coordinate{x + 100, y + 100}})
		}
	}
	return hexes
}

// Draw renders the game state
func (this *Game) Draw(screen *ebiten.Image) {
	screen.Fill(backgroundColor)
	// Draw all hexagons
	for _, hex := range this.hexes {
		drawHexagon(screen, hex)
	}
	for _, hex := range this.hexes {
		drawHexagonConnections(screen, *hex)
		if hex.hovered {
			drawPendingHexagonConnection(screen, *hex, this.nextConnections())
		}
	}
}

// Layout sets the screen size
func (this *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	//return screenWidth, screenHeight
	return outsideWidth, outsideHeight
}

// Update handles game logic updates
func (this *Game) Update() error {
	mouseX, mouseY := ebiten.CursorPosition()
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		this.updateClickedHex(mouseX, mouseY)
	} else {
		this.updateHoveredHex(mouseX, mouseY)
	}
	return nil
}

func (this *Game) updateClickedHex(mouseX, mouseY int) {
	for _, hex := range this.hexes {
		if hex.pointInHexagon(float64(mouseX), float64(mouseY), hexVertexRadius) && hex.available() {
			hex.connections = this.nextConnections()
			this.updateNextConnections()
		}
	}
}

func (this *Game) updateHoveredHex(mouseX, mouseY int) {
	for _, hex := range this.hexes {
		if hex.pointInHexagon(float64(mouseX), float64(mouseY), hexVertexRadius) {
			if hex.available() {
				hex.hovered = true
			}
		} else {
			hex.hovered = false
		}
	}
}

func (this *Game) updateNextConnections() {
	this.nextConnectionsIndex = rand.Intn(len(this.possibleConnections))
}

func (this *Game) nextConnections() []Connection {
	return this.possibleConnections[this.nextConnectionsIndex]
}

// Hex represents a hexagonal tile
type Hex struct {
	col, row    int        // Column and Row
	center      Coordinate // center of the hex
	connections []Connection
	hovered     bool
}

func (this *Hex) available() bool {
	return len(this.connections) == 0
}

// vertexCoordinates returns a slice of coordinates of each vertex of the hexagon
func (this *Hex) vertexCoordinates() []Coordinate {
	vertices := make([]Coordinate, numHexagonSides)
	for i := 0; i < numHexagonSides; i++ {
		angle := math.Pi/3*float64(i) - math.Pi/6
		x := getXCoordinateFromPolar(this.center[0], hexVertexRadius, angle)
		y := getYCoordinateFromPolar(this.center[1], hexVertexRadius, angle)
		vertices[i] = Coordinate{x, y}
	}
	return vertices
}

// hexagonSideCoordinates returns a slice of coordinates of the midpoint of each side of the hexagon
func (this *Hex) hexagonSideCoordinates() []Coordinate {
	sides := make([]Coordinate, numHexagonSides)
	for i := 0; i < numHexagonSides; i++ {
		angle := math.Pi / 3 * float64(i-1)
		x := getXCoordinateFromPolar(this.center[0], hexSideRadius, angle)
		y := getYCoordinateFromPolar(this.center[1], hexSideRadius, angle)
		sides[i] = Coordinate{x, y}
	}
	return sides
}

// pointInHexagon checks if a point is inside the hexagon
func (this *Hex) pointInHexagon(px, py, radius float64) bool {
	buffer := .1 // prevents two hexagons being selected at once
	dx := math.Abs(px-this.center[0]) / radius
	dy := math.Abs(py-this.center[1]) / radius
	return dx <= 1.0-buffer &&
		dy <= math.Sqrt(3.0)/2.0-buffer &&
		dx+dy/math.Sqrt(3.0) <= 1.0-buffer
}

func main() {
	game := NewGame()
	ebiten.SetWindowSize(screenWidth, screenHeight)
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}

///////////////////////////////////////////////// Drawing Functions /////////////////////////////////////////////////

func drawHexagon(screen *ebiten.Image, hex *Hex) {
	vertices := hex.vertexCoordinates()
	vertices = append(vertices, vertices[0])
	for i := 0; i < numHexagonSides; i++ {
		vector.StrokeLine(screen,
			float32(vertices[i][0]),
			float32(vertices[i][1]),
			float32(vertices[i+1][0]),
			float32(vertices[i+1][1]),
			defaultStrokeWidth,
			hexBorderColor,
			false)
	}
}

func drawHexagonConnections(screen *ebiten.Image, hex Hex) {
	if len(hex.connections) == 0 {
		return
	}
	for _, connection := range hex.connections {
		drawHexagonConnection(screen, hex, connection, connectionColor)
	}
}

func drawPendingHexagonConnection(screen *ebiten.Image, hex Hex, nextConnections []Connection) {
	if len(hex.connections) != 0 {
		return
	}
	for _, connection := range nextConnections {
		drawHexagonConnection(screen, hex, connection, pendingConnectionColor)
	}
}

func drawHexagonConnection(screen *ebiten.Image, hex Hex, connection Connection, color color.RGBA) {
	sideA := connection[0]
	sideB := connection[1]
	diff := math.Abs(float64(sideA - sideB))
	// Straight Across
	if diff == 3 {
		drawLineConnection(screen, hex, connection, color)
		return
	}
	// Large Curve
	if diff == 2 {
		centerSide := (sideA + sideB) / 2
		angleToSide := math.Pi*2/3 + math.Pi*1/3*float64(centerSide)
		drawLargeCurveConnection(screen, hex, angleToSide, color)
		return
	}
	if diff == 4 {
		oppositeCenterSide := (sideA + sideB) / 2
		angleToSide := -math.Pi*1/3 + math.Pi*1/3*float64(oppositeCenterSide)
		drawLargeCurveConnection(screen, hex, angleToSide, color)
		return
	}
	// Small Curve
	if diff == 1 {
		drawSmallCurveConnection(screen, hex.vertexCoordinates(), min(sideA, sideB), color)
		return
	}
	if diff == 5 {
		drawSmallCurveConnection(screen, hex.vertexCoordinates(), 5, color)
		return
	}
}

func drawLineConnection(screen *ebiten.Image, hex Hex, connection Connection, color color.RGBA) {
	sideA := connection[0]
	sideB := connection[1]
	sides := hex.hexagonSideCoordinates()
	x1 := float32(sides[sideA][0])
	y1 := float32(sides[sideA][1])
	x2 := float32(sides[sideB][0])
	y2 := float32(sides[sideB][1])
	vector2.StrokeLine(screen, x1, y1, x2, y2, connectionStrokeWidth, color, false)
}

func drawLargeCurveConnection(screen *ebiten.Image, hex Hex, angleToSide float64, color color.Color) {
	x := float32(hex.center[0] - math.Cos(angleToSide)*hexSideRadius*2)
	y := float32(hex.center[1] - math.Sin(angleToSide)*hexSideRadius*2)
	radius := float32(hexVertexRadius + hexSideRadius/2 + defaultStrokeWidth)
	startAngle := float32(angleToSide - math.Pi/6)
	endAngle := float32(angleToSide + math.Pi/6)
	vector2.StrokePartialCircle(screen, x, y, radius, startAngle, endAngle, connectionStrokeWidth, color, false)
}

func drawSmallCurveConnection(screen *ebiten.Image, vertices []Coordinate, vertex int, color color.RGBA) {
	x := float32(vertices[vertex][0])
	y := float32(vertices[vertex][1])
	adjustor := math.Pi / 3 * float32(vertex)
	vector2.StrokePartialCircle(screen, x, y, hexVertexRadius/2, math.Pi/2+adjustor, -math.Pi*5/6+adjustor, connectionStrokeWidth, color, false)
}

//////////////////////////////////////////////// Coordinate Functions ////////////////////////////////////////////////

func getXCoordinateFromPolar(centerX, radius, angle float64) float64 {
	return centerX + radius*math.Cos(angle)
}

func getYCoordinateFromPolar(centerY, radius, angle float64) float64 {
	return centerY + radius*math.Sin(angle)
}
