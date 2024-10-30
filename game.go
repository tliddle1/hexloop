package main

import (
	"image/color"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	color2 "github.com/tliddle1/game/color"
)

const (
	rows       = 5          // Number of hexagon rows
	cols       = rows*4 - 2 //24 // Number of hexagon columns
	marginSize = 30
)

// TODO add points (exponentially more for more hexes, double for multiple loops at once, bonus for clearing board)

var (
	threeLinesConnections  = []Connection{{0, 3}, {1, 4}, {2, 5}}
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

// Game represents the game state
type Game struct {
	hexes                     []*Hex // List of hexagons
	possibleConnections       [][]Connection
	theme                     *color2.Theme
	nextConnectionsIndex      int
	screenWidth, screenHeight int
	squareMode                bool
	disabledTicksLeft         int
	loops                     [][]HexConnection
}

// NewGame initializes the game state
func NewGame(squareMode bool) *Game {
	hexGridWidth := hexSideRadius * (cols + 1)
	hexGridHeight := hexVertexRadius * (rows*3 + 0.5)
	if squareMode {
		hexGridHeight -= hexVertexRadius * 1.5
	}
	screenWidth := hexGridWidth + marginSize*2
	screenHeight := hexGridHeight + marginSize*2
	return &Game{
		hexes:                newHexes(squareMode),
		possibleConnections:  connectionPermutations,
		theme:                color2.NewDefaultTheme(),
		nextConnectionsIndex: rand.Intn(len(connectionPermutations)),
		screenWidth:          int(screenWidth),
		screenHeight:         int(screenHeight),
		squareMode:           squareMode,
	}
}

func newHexes(squareMode bool) (hexes []*Hex) {
	xBuffer := marginSize + hexSideRadius
	yBuffer := float64(marginSize + hexVertexRadius)
	for rowNum := 0; rowNum < rows; rowNum++ {
		for colNum := 0; colNum < cols; colNum++ {
			if squareMode && rowNum == rows-1 && isOdd(colNum) {
				continue
			}
			x := float64(colNum) * hexSideRadius
			y := float64(rowNum) * hexVertexRadius * 3
			// Offset odd rows to create staggered effect
			if colNum%2 != 0 {
				y += hexVertexRadius * 1.5
			}
			newHex := Hex{col: colNum, row: rowNum, center: Coordinate{x + xBuffer, y + yBuffer}}
			if squareMode && rowNum == (rows-1)/2 && colNum == (cols-1)/2 {
				newHex.connections = threeLinesConnections
			}
			hexes = append(hexes, &newHex)
		}
	}
	return hexes
}

// Draw renders the game state
func (this *Game) Draw(screen *ebiten.Image) {
	var hoveredHex *Hex = nil
	screen.Fill(this.theme.BackgroundColor)
	// Draw all hexagons
	for _, hex := range this.hexes {
		drawHexagon(screen, hex, this.theme.HexBorderColor)
	}
	for _, hex := range this.hexes {
		drawHexagonConnections(screen, *hex, this.theme)
		if hex.hovered && this.disabledTicksLeft == 0 {
			//drawPendingHexagonConnection(screen, *hex, this.nextConnections(), this.theme)
			if len(hex.connections) == 0 {
				hoveredHex = hex
				drawHexagon(screen, hoveredHex, this.theme.PendingHexBorderColor)
			}
		}
	}
	this.drawLoops(screen)
	// TODO make unit test
	if hoveredHex != nil {
		nextConns := this.nextConnections()
		for i := range nextConns {
			for j := range nextConns[i] {
				side := nextConns[i][j]
				nextHex := hoveredHex
				drawHexagonConnection(screen, *nextHex, nextConns[i], this.theme.PendingConnectionColors[i%3], this.theme.BackgroundColor)
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
					hexConnection.hex.connections = nil
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
		if hex.pointInHexagon(float64(mouseX), float64(mouseY), hexVertexRadius) && hex.empty() {
			hex.connections = this.nextConnections()
			this.updateNextConnections()
			this.loops = this.checkForCompleteLoops(hex)
		}
	}
}

func (this *Game) updateHoveredHex(mouseX, mouseY int) {
	for _, hex := range this.hexes {
		if hex.pointInHexagon(float64(mouseX), float64(mouseY), hexVertexRadius) {
			if hex.empty() {
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

func (this *Game) getBorderHex(row, col, side int) *Hex {
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

func (this *Game) getHexFromGridPosition(row, col int) *Hex {
	for _, hex := range this.hexes {
		if hex.row == row && hex.col == col {
			return hex
		}
	}
	return nil
}

func (this *Game) checkForCompleteLoops(hex *Hex) [][]HexConnection {
	var loops [][]HexConnection
	if !hex.empty() {
		for _, connection := range hex.connections {
			side := connection[0]
			loopFound, connectedHexes := this.findLoop(
				side,
				hex,
				HexConnection{
					hex:        hex,
					connection: connection,
				},
				[]HexConnection{{
					hex:        hex,
					connection: connection,
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

func (this *Game) findLoop(previousConnectedSide int, curHex *Hex, startHexConnection HexConnection, connectedHexes []HexConnection) (bool, []HexConnection) {
	nextHex := this.getBorderHex(curHex.row, curHex.col, previousConnectedSide)
	if nextHex == nil || nextHex.empty() {
		return false, nil
	}
	//connection with opposite side
	originSide := previousConnectedSide - 3
	if originSide < 0 {
		originSide += 6
	}
	connectedSide := nextHex.connectedSide(originSide)
	if nextHex.Equals(startHexConnection.hex) && (connectedSide == startHexConnection.connection[0] || connectedSide == startHexConnection.connection[1]) {
		return true, connectedHexes
	}
	return this.findLoop(connectedSide, nextHex, startHexConnection, append(connectedHexes, HexConnection{
		hex:        nextHex,
		connection: Connection{originSide, connectedSide},
	}))
}

func (this *Game) drawLoops(screen *ebiten.Image) {
	for _, loop := range this.loops {
		for _, hexConnection := range loop {
			drawHexagonConnection(screen, *hexConnection.hex, hexConnection.connection, color.RGBA{R: 255, G: 0, B: 0, A: 255}, this.theme.BackgroundColor)
		}
	}
}

// TODO make iterator
func (this *Game) drawPendingLoops(screen *ebiten.Image, hex *Hex, side int, hoveredHex *Hex, color color.RGBA) (nextHex *Hex, nextSide int, drawn bool) {
	nextHexConnection := this.getNextHexConnection(hex, side, hoveredHex)
	if nextHexConnection.hex != nil {
		drawHexagonConnection(screen, *nextHexConnection.hex, nextHexConnection.connection, color, this.theme.BackgroundColor)
		return nextHexConnection.hex, nextHexConnection.connection[1], true
	}
	return nil, -1, false
}

func (this *Game) getNextHexConnection(hex *Hex, connectedSide int, hoveredHex *Hex) HexConnection {
	borderHex := this.getBorderHex(hex.row, hex.col, connectedSide)
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
		return HexConnection{
			hex:        borderHex,
			connection: Connection{originSide, nextConnectedSide},
		}
	}
	if borderHex == nil || borderHex == hex || len(borderHex.connections) == 0 {
		return HexConnection{
			hex: nil,
		}
	}
	originSide := connectedSide - 3
	if originSide < 0 {
		originSide += 6
	}
	nextConnectedSide := borderHex.connectedSide(originSide)
	return HexConnection{
		hex:        borderHex,
		connection: Connection{originSide, nextConnectedSide},
	}
}

type HexConnection struct {
	hex        *Hex
	connection Connection
}
