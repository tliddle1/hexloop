package main

import (
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	rows       = 5  // Number of hexagon rows
	cols       = 24 // Number of hexagon columns
	marginSize = 10
)

var connectionPermutations = [][]Connection{
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

// Game represents the game state
type Game struct {
	hexes                     []*Hex // List of hexagons
	possibleConnections       [][]Connection
	theme                     *Theme
	nextConnectionsIndex      int
	screenWidth, screenHeight int
}

// NewGame initializes the game state
func NewGame() *Game {
	hexGridWidth := hexSideRadius * (cols + 1)
	hexGridHeight := hexVertexRadius * (rows*3 + 0.5)
	screenWidth := hexGridWidth + marginSize*2
	screenHeight := hexGridHeight + marginSize*2
	return &Game{
		hexes:                newHexes(),
		possibleConnections:  connectionPermutations,
		theme:                NewDefaultTheme(),
		nextConnectionsIndex: rand.Intn(len(connectionPermutations)),
		screenWidth:          int(screenWidth),
		screenHeight:         int(screenHeight),
	}
}

func newHexes() (hexes []*Hex) {
	xBuffer := marginSize + hexSideRadius
	yBuffer := float64(marginSize + hexVertexRadius)
	for rowNum := 0; rowNum < rows; rowNum++ {
		for colNum := 0; colNum < cols; colNum++ {
			x := float64(colNum) * hexSideRadius
			y := float64(rowNum) * hexVertexRadius * 3
			// Offset odd rows to create staggered effect
			if colNum%2 != 0 {
				y += hexVertexRadius * 1.5
			}
			hexes = append(hexes, &Hex{col: colNum, row: rowNum, center: Coordinate{x + xBuffer, y + yBuffer}})
		}
	}
	return hexes
}

// Draw renders the game state
func (this *Game) Draw(screen *ebiten.Image) {
	screen.Fill(this.theme.BackgroundColor)
	// Draw all hexagons
	for _, hex := range this.hexes {
		drawHexagon(screen, hex, this.theme)
	}
	for _, hex := range this.hexes {
		drawHexagonConnections(screen, *hex, this.theme)
		if hex.hovered {
			drawPendingHexagonConnection(screen, *hex, this.nextConnections(), this.theme)
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
