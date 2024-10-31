package main

import (
	"image/color"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	color2 "github.com/tliddle1/game/color"
	"github.com/tliddle1/game/draw"
	"github.com/tliddle1/game/hexagon"
)

const (
	rows       = 5          // Number of hexagon rows
	cols       = rows*4 - 2 //24 // Number of hexagon columns
	marginSize = 30
)

// TODO add points (exponentially more for more hexes, double for multiple loops at once, bonus for clearing board)

var (
	threeLinesConnections  = []hexagon.Connection{{0, 3}, {1, 4}, {2, 5}}
	connectionPermutations = [][]hexagon.Connection{
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

// Game represents the game state
type Game struct {
	hexes                     []*hexagon.Hex // List of hexagons
	possibleConnections       [][]hexagon.Connection
	theme                     *color2.Theme
	nextConnectionsIndex      int
	screenWidth, screenHeight int
	disabledTicksLeft         int
	loops                     [][]hexagon.HexConnection
}

// NewGame initializes the game state
func NewGame() *Game {
	hexGridWidth := hexagon.HexSideRadius * (cols + 1)
	hexGridHeight := hexagon.HexVertexRadius * (rows*3 + 0.5)
	screenWidth := hexGridWidth + marginSize*2
	screenHeight := hexGridHeight + marginSize*2
	return &Game{
		hexes:                newHexes(),
		possibleConnections:  connectionPermutations,
		theme:                color2.NewDefaultTheme(),
		nextConnectionsIndex: rand.Intn(len(connectionPermutations)),
		screenWidth:          int(screenWidth),
		screenHeight:         int(screenHeight),
	}
}

func newHexes() (hexes []*hexagon.Hex) {
	xBuffer := marginSize + hexagon.HexSideRadius
	yBuffer := float64(marginSize + hexagon.HexVertexRadius)
	for rowNum := 0; rowNum < rows; rowNum++ {
		for colNum := 0; colNum < cols; colNum++ {
			x := float64(colNum) * hexagon.HexSideRadius
			y := float64(rowNum) * hexagon.HexVertexRadius * 3
			// Offset odd rows to create staggered effect
			if colNum%2 != 0 {
				y += hexagon.HexVertexRadius * 1.5
			}
			newHex := hexagon.Hex{Col: colNum, Row: rowNum, Center: hexagon.Coordinate{x + xBuffer, y + yBuffer}}
			hexes = append(hexes, &newHex)
		}
	}
	return hexes
}

// Draw renders the game state
func (this *Game) Draw(screen *ebiten.Image) {
	var hoveredHex *hexagon.Hex = nil
	screen.Fill(this.theme.BackgroundColor)
	for _, hex := range this.hexes {
		draw.Hexagon(screen, hex, this.theme.HexBorderColor)
	}
	for _, hex := range this.hexes {
		draw.HexagonConnections(screen, *hex, this.theme)
		if hex.Hovered && this.disabledTicksLeft == 0 {
			//drawPendingHexagonConnection(screen, *hex, this.nextConnections(), this.theme)
			if len(hex.Connections) == 0 {
				hoveredHex = hex
				draw.Hexagon(screen, hoveredHex, this.theme.PendingHexBorderColor)
			}
		}
	}
	draw.Loops(screen, this.loops, this.theme.BackgroundColor)
	// TODO make unit test
	if hoveredHex != nil {
		nextConns := this.nextConnections()
		for i := range nextConns {
			for j := range nextConns[i] {
				side := nextConns[i][j]
				nextHex := hoveredHex
				draw.HexagonConnection(screen, *nextHex, nextConns[i], this.theme.PendingConnectionColors[i%3], this.theme.BackgroundColor)
				nextSide := side
				drawn := true
				k := 0
				for {
					nextHex, nextSide, drawn = this.drawPendingLoops(screen, nextHex, nextSide, hoveredHex, this.theme.PendingConnectionColors[i%3])
					if !drawn || nextHex == nil || k > rows*cols*5 {
						break
					}
					k++
				}
			}
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
	if this.disabledTicksLeft > 0 {
		this.disabledTicksLeft--
		if this.disabledTicksLeft == 1 {
			for _, loop := range this.loops {
				for _, hexConnection := range loop {
					hexConnection.Hex.Connections = nil
				}
			}
			this.loops = nil
		}
		return nil
	}
	mouseX, mouseY := ebiten.CursorPosition()
	// TODO don't allow dragging?
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		this.updateClickedHex(mouseX, mouseY)
	} else {
		this.updateHoveredHex(mouseX, mouseY)
	}

	return nil
}

func (this *Game) updateClickedHex(mouseX, mouseY int) {
	for _, hex := range this.hexes {
		if hex.PointInHexagon(float64(mouseX), float64(mouseY), hexagon.HexVertexRadius) && hex.Empty() {
			hex.Connections = this.nextConnections()
			this.updateNextConnections()
			this.loops = this.checkForCompleteLoops(hex)
		}
	}
}

func (this *Game) updateHoveredHex(mouseX, mouseY int) {
	for _, hex := range this.hexes {
		if hex.PointInHexagon(float64(mouseX), float64(mouseY), hexagon.HexVertexRadius) {
			if hex.Empty() {
				hex.Hovered = true
			}
		} else {
			hex.Hovered = false
		}
	}
}

func (this *Game) updateNextConnections() {
	this.nextConnectionsIndex = rand.Intn(len(this.possibleConnections))
}

func (this *Game) nextConnections() []hexagon.Connection {
	return this.possibleConnections[this.nextConnectionsIndex]
}

func (this *Game) getBorderHex(row, col, side int) *hexagon.Hex {
	r, c := getBorderHexGridPosition(row, col, side)
	return this.getHexFromGridPosition(r, c)
}

func getBorderHexGridPosition(row, col, side int) (r, c int) {
	switch side {
	case 0:
		if isEven(col) {
			return row - 1, col + 1
		} else {
			return row, col + 1
		}
	case 1:
		return row, col + 2
	case 2:
		if isOdd(col) {
			return row + 1, col + 1
		} else {
			return row, col + 1
		}
	case 3:
		if isOdd(col) {
			return row + 1, col - 1
		} else {
			return row, col - 1
		}
	case 4:
		return row, col - 2
	case 5:
		if isEven(col) {
			return row - 1, col - 1
		} else {
			return row, col - 1
		}
	default:
		panic("bad side")
		return row, col
	}
}

func (this *Game) getHexFromGridPosition(row, col int) *hexagon.Hex {
	for _, hex := range this.hexes {
		if hex.Row == row && hex.Col == col {
			return hex
		}
	}
	return nil
}

func (this *Game) checkForCompleteLoops(hex *hexagon.Hex) (loops [][]hexagon.HexConnection) {
	if !hex.Empty() {
		for _, connection := range hex.Connections {
			side := connection[0]
			loopFound, connectedHexes := this.findLoop(
				side,
				hex,
				hexagon.HexConnection{
					Hex:        hex,
					Connection: connection,
				},
				[]hexagon.HexConnection{{
					Hex:        hex,
					Connection: connection,
				}})
			if loopFound {
				this.disabledTicksLeft = 50
				loops = append(loops, connectedHexes)
			}
		}
	}
	// check loops for duplicates
	return loops
}

func (this *Game) findLoop(previousConnectedSide int, curHex *hexagon.Hex, startHexConnection hexagon.HexConnection, connectedHexes []hexagon.HexConnection) (bool, []hexagon.HexConnection) {
	nextHex := this.getBorderHex(curHex.Row, curHex.Col, previousConnectedSide)
	if nextHex == nil || nextHex.Empty() {
		return false, nil
	}
	//connection with opposite side
	originSide := previousConnectedSide - 3
	if originSide < 0 {
		originSide += 6
	}
	connectedSide := nextHex.ConnectedSide(originSide)
	if nextHex.Equals(startHexConnection.Hex) && (connectedSide == startHexConnection.Connection[0] || connectedSide == startHexConnection.Connection[1]) {
		return true, connectedHexes
	}
	return this.findLoop(connectedSide, nextHex, startHexConnection, append(connectedHexes, hexagon.HexConnection{
		Hex:        nextHex,
		Connection: hexagon.Connection{originSide, connectedSide},
	}))
}

// TODO make iterator
func (this *Game) drawPendingLoops(screen *ebiten.Image, hex *hexagon.Hex, side int, hoveredHex *hexagon.Hex, color color.RGBA) (nextHex *hexagon.Hex, nextSide int, drawn bool) {
	nextHexConnection := this.getNextHexConnection(hex, side, hoveredHex)
	if nextHexConnection.Hex != nil {
		draw.HexagonConnection(screen, *nextHexConnection.Hex, nextHexConnection.Connection, color, this.theme.BackgroundColor)
		return nextHexConnection.Hex, nextHexConnection.Connection[1], true
	}
	return nil, -1, false
}

func (this *Game) getNextHexConnection(hex *hexagon.Hex, connectedSide int, hoveredHex *hexagon.Hex) hexagon.HexConnection {
	borderHex := this.getBorderHex(hex.Row, hex.Col, connectedSide)
	if borderHex == hoveredHex {
		originSide := connectedSide - 3
		if originSide < 0 {
			originSide += 6
		}
		var nextConnectedSide int
		for _, connection := range this.nextConnections() {
			if connection[0] == originSide {
				nextConnectedSide = connection[1]
			}
			if connection[1] == originSide {
				nextConnectedSide = connection[0]
			}
		}
		return hexagon.HexConnection{
			Hex:        borderHex,
			Connection: hexagon.Connection{originSide, nextConnectedSide},
		}
	}
	if borderHex == nil || borderHex == hex || len(borderHex.Connections) == 0 {
		return hexagon.HexConnection{
			Hex: nil,
		}
	}
	originSide := connectedSide - 3
	if originSide < 0 {
		originSide += 6
	}
	nextConnectedSide := borderHex.ConnectedSide(originSide)
	return hexagon.HexConnection{
		Hex:        borderHex,
		Connection: hexagon.Connection{originSide, nextConnectedSide},
	}
}

func isEven(x int) bool {
	return x%2 == 0
}

func isOdd(x int) bool {
	return x%2 != 0
}
