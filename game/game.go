package game

import (
	"bytes"
	"fmt"
	"image/color"
	"log"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	color2 "github.com/tliddle1/hexloop/color"
	"github.com/tliddle1/hexloop/draw"
	"github.com/tliddle1/hexloop/hexagon"
)

const (
	// board
	rows          = 5          // Number of hexagon rows
	cols          = rows*4 - 2 // Number of hexagon columns
	marginSize    = 30
	scoreTextSize = 24
	// points
	clearBoardBonus  = 5_000
	lowestPointValue = 1.0
	increment        = 1.0
)

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

// TODO make unit tests
// TODO Game Over (Message that it's over and option to play again)
// TODO Start Over (Menu)
// TODO Landing Page (Play, Themes, How To Play)
// TODO Puzzle Mode (set of tiles to make one loop?)
// TODO Challenge Mode (obstacles?, hexes to clear?
// TODO Add internal timer (https://arc.net/l/quote/vlbnnjos)

// Game represents the game state
type Game struct {
	hexes                     []*hexagon.Hex // List of hexagons
	possibleConnections       [][]hexagon.Connection
	loops                     [][]hexagon.HexConnection
	theme                     *color2.Theme
	nextConnectionsIndex      int
	ScreenWidth, ScreenHeight int
	disabledTicksLeft         int
	score                     int
	highScore                 int // TODO
	gameInProgress            bool
}

// NewGame initializes the game state
func NewGame() *Game {
	hexGridWidth := hexagon.HexSideRadius * (cols + 1) // +3 in parentheses if you want to accommodate for the current hexagon on the sidebar
	hexGridHeight := hexagon.HexVertexRadius * (rows*3 + 0.5)
	screenWidth := int(hexGridWidth) + marginSize*2
	screenHeight := int(hexGridHeight) + marginSize*2 + scoreTextSize + scoreTextSize
	return &Game{
		hexes:                newHexes(),
		possibleConnections:  connectionPermutations,
		theme:                color2.NewDefaultTheme(),
		nextConnectionsIndex: rand.Intn(len(connectionPermutations)),
		ScreenWidth:          screenWidth,
		ScreenHeight:         screenHeight,
		gameInProgress:       true,
	}
}

func newHexes() (hexes []*hexagon.Hex) {
	origin := getFirstHexCoordinate()
	for rowNum := 0; rowNum < rows; rowNum++ {
		for colNum := 0; colNum < cols; colNum++ {
			hexes = append(hexes, hexagon.NewHex(colNum, rowNum, origin[0], origin[1]))
		}
	}
	return hexes
}

// Draw renders the game state
func (this *Game) Draw(screen *ebiten.Image) {
	screen.Fill(this.theme.BackgroundColor)
	this.drawScore(screen)
	this.drawHighScore(screen)
	this.drawHexagonBoard(screen)
	this.drawPlacedHexagons(screen)
	//this.drawCurrentHexPattern(screen)
	this.drawPendingHex(screen)
	this.drawCompletedLoops(screen)
}

func (this *Game) getHoveredHex() *hexagon.Hex {
	for _, hex := range this.hexes {
		if hex.Hovered && this.disabledTicksLeft == 0 {
			if len(hex.Connections) == 0 {
				return hex
			}
		}
	}
	return nil
}

// Layout sets the screen size
func (this *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	//return screenWidth, screenHeight
	return outsideWidth, outsideHeight
}

// Update handles game logic updates
func (this *Game) Update() error {
	if this.gameOver() && this.gameInProgress {
		this.gameInProgress = false
		this.highScore = max(this.highScore, this.score)
		this.startOver()
	}
	if this.gameInProgress {
		// this.updateGameInProgress()
		if this.disabledTicksLeft > 0 {
			this.disabledTicksLeft--
			if this.disabledTicksLeft == 0 {
				this.removeCompletedLoops()
				if this.boardEmpty() {
					this.updateScore(clearBoardBonus)
				}
			}
			return nil
		}

		mouseX, mouseY := ebiten.CursorPosition()
		// TODO don't allow dragging?
		if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
			this.updateClickedHex(mouseX, mouseY)
			this.updateScore(calculatePoints(this.loops))
		} else {
			this.updateHoveredHex(mouseX, mouseY)
		}
	} else {
		if this.disabledTicksLeft > 0 {
			this.disabledTicksLeft--
		}
		if this.disabledTicksLeft == 0 {
			this.gameInProgress = true
		}
	}
	return nil
}

// Score

func (this *Game) updateScore(points int) {
	this.score += points
}

func calculatePoints(loops [][]hexagon.HexConnection) int {
	connectionPoints := 0
	for _, loop := range loops {
		connectionPoints += loopPointFormula(len(loop))
	}
	return connectionPoints * len(loops)
}

func loopPointFormula(n int) int {
	nFloat := float64(n)
	return int((nFloat / 2) * ((2 * lowestPointValue) + (nFloat-1)*increment))
}

// Board

func (this *Game) updateClickedHex(mouseX, mouseY int) {
	for _, hex := range this.hexes {
		if hex.PointInHexagon(float64(mouseX), float64(mouseY), hexagon.HexVertexRadius) && hex.Empty() {
			hex.Connections = this.nextConnections()
			this.updateNextConnections()
			this.loops = this.checkForCompleteLoops(hex)
		}
	}
}

func (this *Game) updateNextConnections() {
	this.nextConnectionsIndex = rand.Intn(len(this.possibleConnections))
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

func (this *Game) boardEmpty() bool {
	for _, hex := range this.hexes {
		if !hex.Empty() {
			return false
		}
	}
	return true
}

func (this *Game) nextConnections() []hexagon.Connection {
	return this.possibleConnections[this.nextConnectionsIndex]
}

func (this *Game) removeCompletedLoops() {
	for _, loop := range this.loops {
		for _, hexConnection := range loop {
			hexConnection.Hex.Connections = nil
		}
	}
	this.loops = nil
}

func (this *Game) drawScore(screen *ebiten.Image) {
	text.Draw(screen, scoreString(this.score), getTextFace(), getDrawScoreOptions(this.theme.ConnectionColor))
}

func (this *Game) drawHexagonBoard(screen *ebiten.Image) {
	for _, hex := range this.hexes {
		draw.Hexagon(screen, hex, this.theme.HexBorderColor)
	}
}

func (this *Game) drawPlacedHexagons(screen *ebiten.Image) {
	for _, hex := range this.hexes {
		draw.HexagonConnections(screen, *hex, this.theme)
	}
}

func (this *Game) drawPendingConnections(screen *ebiten.Image, hex *hexagon.Hex) {
	// TODO make connection normal color (or black) if it touches a wall
	nextConns := this.nextConnections()
	for i := range nextConns {
		for j := range nextConns[i] {
			side := nextConns[i][j]
			nextHex := hex
			draw.HexagonConnection(screen, *nextHex, nextConns[i], this.theme.PendingConnectionColors[i%3], this.theme.BackgroundColor)
			nextSide := side
			drawn := true
			k := 0
			for {
				nextHex, nextSide, drawn = this.drawPendingLoops(screen, nextHex, nextSide, hex, this.theme.PendingConnectionColors[i%3])
				if !drawn || nextHex == nil || k > rows*cols*5 {
					break
				}
				k++
			}
		}
	}
}

func (this *Game) drawCurrentHexPattern(screen *ebiten.Image) {
	origin := getFirstHexCoordinate()
	currentHex := hexagon.NewHex(cols+2, 0, origin[0], origin[1])
	draw.Hexagon(screen, currentHex, this.theme.PendingHexBorderColor)
	this.drawPendingConnections(screen, currentHex)
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

func getFirstHexCoordinate() hexagon.Coordinate {
	xBuffer := marginSize + hexagon.HexSideRadius
	yBuffer := float64(marginSize + hexagon.HexVertexRadius + (scoreTextSize * 2))
	return [2]float64{xBuffer, yBuffer}
}

func (this *Game) drawPendingHex(screen *ebiten.Image) {
	hoveredHex := this.getHoveredHex()
	if hoveredHex != nil {
		draw.Hexagon(screen, hoveredHex, this.theme.PendingHexBorderColor)
		this.drawPendingConnections(screen, hoveredHex)
	}
}

func (this *Game) drawCompletedLoops(screen *ebiten.Image) {
	draw.Loops(screen, this.loops, this.theme.CompletedLoopColor, this.theme.BackgroundColor)
}

func (this *Game) gameOver() bool {
	for _, hex := range this.hexes {
		if hex.Empty() {
			return false
		}
	}
	return this.disabledTicksLeft == 0
}

func (this *Game) drawHighScore(screen *ebiten.Image) {
	text.Draw(screen, highScoreString(this.highScore), getTextFace(), getDrawHighScoreOptions(this.theme.ConnectionColor))
}

func (this *Game) startOver() {
	for _, hex := range this.hexes {
		hex.Reset()
	}
	this.score = 0
	this.loops = nil
	this.disabledTicksLeft = 200
}

func highScoreString(score int) string {
	// todo add commas
	return fmt.Sprintf("High Score: %d", score)
}

func scoreString(score int) string {
	// todo add commas
	return fmt.Sprintf("Score: %d", score)
}

func getTextFace() text.Face {
	fontFaceSource, err := text.NewGoTextFaceSource(bytes.NewReader(fonts.MPlus1pRegular_ttf))
	if err != nil {
		log.Fatal(err)
	}
	textFace := &text.GoTextFace{
		Source: fontFaceSource,
		Size:   scoreTextSize,
	}
	return textFace
}

func getDrawScoreOptions(clr color.RGBA) *text.DrawOptions {
	drawOptions := &text.DrawOptions{}
	drawOptions.GeoM.Translate(79, marginSize/2+marginSize)
	drawOptions.ColorScale.ScaleWithColor(clr)
	return drawOptions
}
func getDrawHighScoreOptions(clr color.RGBA) *text.DrawOptions {
	drawOptions := &text.DrawOptions{}
	drawOptions.GeoM.Translate(20, marginSize/2)
	drawOptions.ColorScale.ScaleWithColor(clr)
	return drawOptions
}

func isEven(x int) bool {
	return x%2 == 0
}

func isOdd(x int) bool {
	return x%2 != 0
}
