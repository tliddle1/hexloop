package game

import (
	"bytes"
	"fmt"
	"image/color"
	"log"
	"math/rand"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	color2 "github.com/tliddle1/hexloop/color"
	"github.com/tliddle1/hexloop/draw"
	"github.com/tliddle1/hexloop/hexagon"
	"github.com/tliddle1/hexloop/vector"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
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

const (
	titleScreen = iota
	gameScreen
)

var (
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

type scene uint8

// Game represents the game state
type Game struct {
	hexes                     []*hexagon.Hex // List of hexagons
	possibleConnections       [][]hexagon.Connection
	loops                     []hexagon.Loop
	theme                     *color2.Theme
	nextConnectionsIndex      int
	ScreenWidth, ScreenHeight int
	disabledTicksLeft         int
	score                     int
	highScore                 int // TODO get highScore to work on web
	gameInProgress            bool
	currentScene              scene
	titleHexes                []*hexagon.TextHexagon
	printer                   *message.Printer
}

// NewGame initializes the game state
func NewGame() *Game {
	hexGridWidth := hexagon.HexSideRadius * (cols + 1) // +3 in parentheses if you want to accommodate for the current hexagon on the sidebar
	hexGridHeight := hexagon.HexVertexRadius * (rows*3 + 0.5)
	screenWidth := int(hexGridWidth) + marginSize*2
	screenHeight := int(hexGridHeight) + marginSize*2 + scoreTextSize + scoreTextSize
	return &Game{
		hexes:                newHexes(rows, cols, hexagon.HexVertexRadius, draw.HexagonStrokeWidth, draw.ConnectionWidth, getGameBoardFirstHexCoordinate()),
		possibleConnections:  connectionPermutations,
		theme:                color2.NewDefaultTheme(),
		nextConnectionsIndex: rand.Intn(len(connectionPermutations)),
		ScreenWidth:          screenWidth,
		ScreenHeight:         screenHeight,
		gameInProgress:       true,
		currentScene:         titleScreen,
		titleHexes:           newTitleHexes(screenWidth, screenHeight),
		printer:              message.NewPrinter(language.English),
	}
}

func newTitleHexes(screenWidth, screenHeight int) (titleHexes []*hexagon.TextHexagon) {
	originX := float64(screenWidth) / 2
	originY := float64(screenHeight)/2 - (hexagon.HexVertexRadiusTest * 2.5)
	for row := range 2 {
		for col := -3; col < 4; col++ {
			str := ""
			textSize := float64(scoreTextSize)
			if row == 0 && col == -2 {
				str = "H"
				textSize = float64(scoreTextSize * 3)
			} else if row == 0 && col == 0 {
				str = "E"
				textSize = float64(scoreTextSize * 3)
			} else if row == 0 && col == 2 {
				str = "X"
				textSize = float64(scoreTextSize * 3)
			} else if row == 0 && col == -3 {
				str = "L"
				textSize = float64(scoreTextSize * 3)
			} else if row == 0 && (col == -1 || col == 1) {
				str = "O"
				textSize = float64(scoreTextSize * 3)
			} else if row == 0 && col == 3 {
				str = "P"
				textSize = float64(scoreTextSize * 3)
			} else if row == 1 && col == -3 {
				str = "Start"
			} else if row == 1 && col == 3 {
				str = "Tutorial"
			}
			hex := hexagon.NewTextHexagon(col, row, originX, originY, hexagon.HexVertexRadiusTest, draw.TitleHexagonStrokeWidth, draw.TitleConnectionWidth, str, textSize)
			titleHexes = append(titleHexes, hex)
		}
	}
	return titleHexes
}

func newHexes(numRows, numCols int, vertexRadius float64, edgeWidth, connectionWidth float32, origin hexagon.Coordinate) (hexes []*hexagon.Hex) {
	for rowNum := 0; rowNum < numRows; rowNum++ {
		for colNum := 0; colNum < numCols; colNum++ {
			hexes = append(hexes, hexagon.NewHex(colNum, rowNum, origin[0], origin[1], vertexRadius, edgeWidth, connectionWidth))
		}
	}
	return hexes
}

// Draw renders the game state
func (this *Game) Draw(screen *ebiten.Image) {
	screen.Fill(this.theme.BackgroundColor)
	//this.drawScreenBorder(screen)
	if this.currentScene == titleScreen {
		this.drawTitleScreen(screen)
	} else if this.currentScene == gameScreen {
		this.drawGameScreen(screen)
	} else {
		panic("unknown scene")
	}
}

// Layout sets the screen size
func (this *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	//return screenWidth, screenHeight
	return outsideWidth, outsideHeight
}

// Update handles game logic updates
func (this *Game) Update() error {
	if this.currentScene == titleScreen {
		this.updateTitleScreen()
	} else if this.currentScene == gameScreen {
		if this.gameInProgress {
			this.updateGameInProgress()
		} else {
			if this.disabledTicksLeft > 0 {
				this.disabledTicksLeft--
			} else if this.disabledTicksLeft == 0 {
				this.startOver()
				this.gameInProgress = true
			}
		}
	}
	return nil
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

// Score

func (this *Game) updateScore(points int) {
	this.score += points
}

func calculatePoints(loops []hexagon.Loop) int {
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
		if hex.PointInHexagon(float64(mouseX), float64(mouseY)) && hex.Empty() {
			hex.Connections = this.nextConnections()
			this.updateNextConnections()
			this.loops = this.getCompleteLoops(hex)
			if len(this.loops) > 0 {
				this.disabledTicksLeft = 50
			}
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

func (this *Game) getCompleteLoops(hex *hexagon.Hex) (loops []hexagon.Loop) {
	if !hex.Empty() {
		for _, connection := range hex.Connections {
			side := connection[0]
			loop, completedLoop, _ := this.findLoop(
				side,
				hex,
				hexagon.HexConnection{
					Hex:        hex,
					Connection: connection,
				},
				hexagon.Loop{{
					Hex:        hex,
					Connection: connection,
				}})
			if completedLoop {
				loops = append(loops, loop)
			}
		}
	}
	return loops
}

func (this *Game) getIncompleteLoops(hex *hexagon.Hex) (loops []hexagon.Loop, touchesWall bool) {
	if !hex.Empty() {
		for _, connection := range hex.Connections {
			side := connection[0]
			connectedHexes, loopFound, _ := this.findLoop(
				side,
				hex,
				hexagon.HexConnection{
					Hex:        hex,
					Connection: connection,
				},
				hexagon.Loop{{
					Hex:        hex,
					Connection: connection,
				}})
			if loopFound {
				loops = append(loops, connectedHexes)
			}
		}
	}
	return loops, false
}

const (
	connectedToEdge = iota
	connectedToEmpty
)

func (this *Game) findLoop(previousConnectedSide int, curHex *hexagon.Hex, startHexConnection hexagon.HexConnection, connectedHexes hexagon.Loop) (loop hexagon.Loop, completed bool, reason int) {
	nextHex := this.getBorderHex(curHex.Row, curHex.Col, previousConnectedSide)
	if nextHex == nil {
		return connectedHexes, false, connectedToEdge
	} else if nextHex.Empty() {
		return connectedHexes, false, connectedToEmpty
	}
	//connection with opposite side
	originSide := previousConnectedSide - 3
	if originSide < 0 {
		originSide += 6
	}
	connectedSide := nextHex.ConnectedSide(originSide)
	if nextHex.Equals(startHexConnection.Hex) && (connectedSide == startHexConnection.Connection[0] || connectedSide == startHexConnection.Connection[1]) {
		return connectedHexes, true, -1
	}
	return this.findLoop(connectedSide, nextHex, startHexConnection, append(connectedHexes, hexagon.HexConnection{
		Hex:        nextHex,
		Connection: hexagon.Connection{originSide, connectedSide},
	}))
}

func (this *Game) updateHoveredHex(mouseX, mouseY int) {
	for _, hex := range this.hexes {
		if hex.PointInHexagon(float64(mouseX), float64(mouseY)) {
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
	text.Draw(screen, this.scoreString(this.score), getTextFace(scoreTextSize), getDrawScoreOptions(this.theme.ConnectionColor))
}

func (this *Game) drawHexagonGameBoard(screen *ebiten.Image) {
	for _, hex := range this.hexes {
		draw.Hexagon(screen, hex, this.theme.HexBorderColor)
	}
}

func (this *Game) drawHexagonTitleBoard(screen *ebiten.Image) {
	for _, hex := range this.titleHexes {
		draw.TextHexagon(screen, hex, this.theme.HexBorderColor, this.theme.ConnectionColor)
	}
}

func (this *Game) drawPlacedHexagons(screen *ebiten.Image) {
	for _, hex := range this.hexes {
		draw.HexagonConnections(screen, hex, this.theme.ConnectionColor, this.theme)
	}
}

func (this *Game) drawPendingConnections(screen *ebiten.Image, hex *hexagon.Hex) {
	// TODO make connection normal color (or black) if it touches a wall
	// TODO make connection completed color if looped
	nextConns := this.nextConnections()
	for i := range nextConns {
		for j := range nextConns[i] {
			side := nextConns[i][j]
			nextHex := hex
			draw.HexagonConnection(screen, nextHex, nextConns[i], this.theme.PendingConnectionColors[i%3], this.theme.BackgroundColor)
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
	origin := getGameBoardFirstHexCoordinate()
	currentHex := hexagon.NewHex(cols+2, 0, origin[0], origin[1], hexagon.HexVertexRadius, draw.HexagonStrokeWidth, draw.ConnectionWidth)
	draw.Hexagon(screen, currentHex, this.theme.PendingHexBorderColor)
	this.drawPendingConnections(screen, currentHex)
}

// TODO make iterator
func (this *Game) drawPendingLoops(screen *ebiten.Image, hex *hexagon.Hex, side int, hoveredHex *hexagon.Hex, color color.RGBA) (nextHex *hexagon.Hex, nextSide int, drawn bool) {
	nextHexConnection := this.getNextHexConnection(hex, side, hoveredHex)
	if nextHexConnection.Hex != nil {
		draw.HexagonConnection(screen, nextHexConnection.Hex, nextHexConnection.Connection, color, this.theme.BackgroundColor)
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

func getGameBoardFirstHexCoordinate() hexagon.Coordinate {
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
	text.Draw(screen, this.highScoreString(this.highScore), getTextFace(scoreTextSize), getDrawHighScoreOptions(this.theme.ConnectionColor))
}

func (this *Game) startOver() {
	for _, hex := range this.hexes {
		hex.Reset()
	}
	this.score = 0
	this.loops = nil
}

func (this *Game) updateGameInProgress() {
	if this.gameOver() {
		this.gameInProgress = false
		this.highScore = max(this.highScore, this.score)
	}
	if this.disabledTicksLeft > 0 {
		this.disabledTicksLeft--
		if this.disabledTicksLeft == 0 {
			newGame := this.boardEmpty()
			this.removeCompletedLoops()
			if this.boardEmpty() && !newGame {
				this.updateScore(clearBoardBonus)
			}
		}
		return
	}

	mouseX, mouseY := ebiten.CursorPosition()
	// TODO don't allow dragging?
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		this.updateClickedHex(mouseX, mouseY)
		this.updateScore(calculatePoints(this.loops))
	} else {
		this.updateHoveredHex(mouseX, mouseY)
	}
}

func (this *Game) drawGameScreen(screen *ebiten.Image) {
	this.drawScore(screen)
	this.drawHighScore(screen)
	this.drawHexagonGameBoard(screen)
	this.drawPlacedHexagons(screen)
	//this.drawCurrentHexPattern(screen)
	this.drawPendingHex(screen)
	this.drawCompletedLoops(screen)
}

func (this *Game) drawScreenBorder(screen *ebiten.Image) {
	strokeWidth := float32(10)
	vector.StrokeLine(screen, 0, 0, float32(this.ScreenWidth)+(strokeWidth/2), 0, strokeWidth, this.theme.ConnectionColor, true)
	vector.StrokeLine(screen, 0, 0, 0, float32(this.ScreenHeight)+(strokeWidth/2), strokeWidth, this.theme.ConnectionColor, true)
	vector.StrokeLine(screen, 0, float32(this.ScreenHeight), float32(this.ScreenWidth)+(strokeWidth/2), float32(this.ScreenHeight), strokeWidth, this.theme.ConnectionColor, true)
	vector.StrokeLine(screen, float32(this.ScreenWidth), 0, float32(this.ScreenWidth), float32(this.ScreenHeight)+(strokeWidth/2), strokeWidth, this.theme.ConnectionColor, true)
}

func (this *Game) drawTitleScreen(screen *ebiten.Image) {
	this.drawHexagonTitleBoard(screen)
}

func (this *Game) updateTitleScreen() {
	mouseX, mouseY := ebiten.CursorPosition()
	// TODO don't allow dragging?
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		startButton := this.getStartButton()
		if startButton != nil && startButton.PointInHexagon(float64(mouseX), float64(mouseY)) {
			this.currentScene = gameScreen
			// TODO handle transition gracefully
			this.disabledTicksLeft = 50
		}
	}
}

func (this *Game) getStartButton() *hexagon.TextHexagon {
	for _, hex := range this.titleHexes {
		if hex.Row == 1 && hex.Col == -3 {
			return hex
		}
	}
	return nil
}

func (this *Game) highScoreString(score int) string {
	scoreStr := strconv.Itoa(score)
	return "High Score: " + this.withCommas(scoreStr)
}

func (this *Game) scoreString(score int) string {
	scoreStr := strconv.Itoa(score)
	return "Score: " + this.withCommas(scoreStr)
}

func getTextFace(textSize float64) text.Face {
	fontFaceSource, err := text.NewGoTextFaceSource(bytes.NewReader(fonts.MPlus1pRegular_ttf))
	if err != nil {
		log.Fatal(err)
	}
	textFace := &text.GoTextFace{
		Source: fontFaceSource,
		Size:   textSize,
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

func (this *Game) withCommas(s string) string {
	if len(s) <= 3 {
		return s
	}
	return this.withCommas(s[:len(s)-3]) + "," + s[len(s)-3:]
}
