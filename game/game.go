package game

import (
	"bytes"
	"image/color"
	"log"
	"math/rand"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	color2 "github.com/tliddle1/hexloop/color"
	"github.com/tliddle1/hexloop/draw"
	"github.com/tliddle1/hexloop/hexagon"
	"github.com/tliddle1/hexloop/vector"
)

const (
	// board
	rows          = 5          // Number of hexagon rows
	cols          = rows*4 - 2 // Number of hexagon columns
	marginSize    = 30
	smallTextSize = 24
	// points
	clearBoardBonus  = 5_000
	lowestPointValue = 1.0
	increment        = 1.0
)

const (
	titleScreen = iota
	gameScreen
	tutorialScreenExplanation
	tutorialScreen1
	tutorialScreen2
	//hexGridWidth = hexagon.HexSideRadius * (cols + 1) // +3 in parentheses if you want to accommodate for the current hexagon on the sidebar
	hexGridHeight = hexagon.HexVertexRadius * (rows*3 + 0.5)
	screenWidth   = 494 + marginSize*2 // int(hexGridWidth) + marginSize*2
	screenHeight  = int(hexGridHeight) + marginSize*2 + smallTextSize + smallTextSize
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

// TODO add drawings to image to improve performance
// TODO make unit tests
// TODO make clickableShape interface (arrow, hexagon, etc.)
// TODO Play Game button then start button (don't start until cursor is up again)
// TODO Game Over (Message that it's over and option to play again)
// TODO Start Over (Menu)
// TODO Landing Page (Play, Themes, How To Play)
// TODO Smaller board with three options ?
// TODO Puzzle Mode (set of tiles to make one loop?)
// TODO Challenge Mode (obstacles?, hexes to clear?
// TODO Add internal timer (https://arc.net/l/quote/vlbnnjos)

type sceneType uint8

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
	highScore                 int // TODO get highScore to save on web
	gameInProgress            bool
	currentSceneType          sceneType
	titleHexes                []*hexagon.TextHexagon
	titleBoardImage           *ebiten.Image
	nextArrowHovered          bool
	startButton               *hexagon.TextHexagon
	tutorialButton            *hexagon.TextHexagon
}

// NewGame initializes the game state
func NewGame() *Game {
	titleHexes, startButton, tutorialButton := newTitleHexes(screenWidth, screenHeight)
	g := Game{
		hexes:                newHexes(rows, cols, hexagon.HexVertexRadius, draw.HexagonStrokeWidth, draw.ConnectionWidth, getGameBoardFirstHexCoordinate()),
		possibleConnections:  connectionPermutations,
		theme:                color2.NewDefaultTheme(),
		nextConnectionsIndex: rand.Intn(len(connectionPermutations)),
		ScreenWidth:          screenWidth,
		ScreenHeight:         screenHeight,
		gameInProgress:       true,
		currentSceneType:     titleScreen,
		titleHexes:           titleHexes,
		startButton:          startButton,
		tutorialButton:       tutorialButton,
	}
	g.generateTitleBoardImage(screenWidth, screenHeight)
	return &g
}

func newTitleHexes(screenWidth, screenHeight int) (titleHexes []*hexagon.TextHexagon, startButton, tutorialButton *hexagon.TextHexagon) {
	originX := float64(screenWidth) / 2
	originY := float64(screenHeight)/2 - (hexagon.HexVertexRadiusTest * 2.5)
	startButtonText := "Start"
	tutorialButtonText := "How to Play"
	for row := range 2 {
		for col := -3; col < 4; col++ {
			addConnections := false
			str := ""
			textSize := float64(smallTextSize)
			if row == 0 && col == -2 {
				str = "H"
				textSize = float64(smallTextSize * 3)
			} else if row == 0 && col == 0 {
				str = "E"
				textSize = float64(smallTextSize * 3)
			} else if row == 0 && col == 2 {
				str = "X"
				textSize = float64(smallTextSize * 3)
			} else if row == 0 && col == -3 {
				str = "L"
				textSize = float64(smallTextSize * 3)
			} else if row == 0 && (col == -1 || col == 1) {
				str = "O"
				textSize = float64(smallTextSize * 3)
			} else if row == 0 && col == 3 {
				str = "P"
				textSize = float64(smallTextSize * 3)
			} else if row == 1 && col == -3 {
				str = startButtonText
			} else if row == 1 && col == 3 {
				str = tutorialButtonText
			} else {
				addConnections = true
			}
			hex := hexagon.NewTextHexagon(col, row, originX, originY, hexagon.HexVertexRadiusTest, draw.TitleHexagonStrokeWidth, draw.TitleConnectionWidth, str, textSize)
			if str == startButtonText {
				startButton = hex
			}
			if str == tutorialButtonText {
				tutorialButton = hex
			}
			if addConnections {
				hex.Connections = connectionPermutations[rand.Intn(len(connectionPermutations))]
			}
			titleHexes = append(titleHexes, hex)
		}
	}
	return titleHexes, startButton, tutorialButton
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
	if this.currentSceneType == titleScreen {
		this.drawTitleScreen(screen)
	} else if this.currentSceneType == gameScreen {
		this.drawGameScreen(screen)
	} else if this.currentSceneType == tutorialScreenExplanation {
		this.drawTutorialScreenExplanation(screen)
	} else if this.currentSceneType == tutorialScreen1 {
		this.drawTutorialScreen1(screen)
	} else if this.currentSceneType == tutorialScreen2 {
		this.drawTutorialScreen2(screen)
	} else {
		panic("unknown sceneType")
	}
}

// Layout sets the screen size
func (this *Game) Layout(_, _ int) (int, int) {
	return screenWidth, screenHeight
	//return outsideWidth, outsideHeight
}

// Update handles game logic updates
func (this *Game) Update() error {
	switch this.currentSceneType {
	case gameScreen:
		this.updateGameScreen()
	case tutorialScreenExplanation:
		this.updateTutorialExplanationScreen()
	case tutorialScreen1:
		this.updateTutorial1Screen()
	case tutorialScreen2:
		this.updateTutorial2Screen()
	case titleScreen:
		this.updateTitleScreen()
	default:
		this.currentSceneType = titleScreen
		this.updateTitleScreen()
	}
	return nil
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func (this *Game) drawScore(screen *ebiten.Image) {
	text.Draw(screen, this.scoreString(this.score), getTextFace(smallTextSize), getDrawScoreOptions(this.theme.ConnectionColor))
}

func (this *Game) drawHexagonGameBoard(screen *ebiten.Image) {
	for _, hex := range this.hexes {
		draw.Hexagon(screen, hex, this.theme.HexBorderColor)
	}
}

func (g *Game) generateTitleBoardImage(width, height int) {
	img := ebiten.NewImage(width, height)

	for _, hex := range g.titleHexes {
		draw.TextHexagon(img, hex, g.theme.HexBorderColor, g.theme.ConnectionColor)
	}
	for _, hex := range g.titleHexes {
		draw.HexagonConnections(img, hex.Hex, g.theme.ConnectionColor, g.theme)
	}

	g.titleBoardImage = img
}

func (this *Game) drawHexagonTitleBoard(screen *ebiten.Image) {
	if this.titleBoardImage != nil {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(0, 0) // Or wherever you want to place it
		screen.DrawImage(this.titleBoardImage, op)
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

func (this *Game) drawPendingLoops(screen *ebiten.Image, hex *hexagon.Hex, side int, hoveredHex *hexagon.Hex, color color.RGBA) (nextHex *hexagon.Hex, nextSide int, drawn bool) {
	// TODO make iterator
	nextHexConnection := this.getNextHexConnection(hex, side, hoveredHex)
	if nextHexConnection.Hex != nil {
		draw.HexagonConnection(screen, nextHexConnection.Hex, nextHexConnection.Connection, color, this.theme.BackgroundColor)
		return nextHexConnection.Hex, nextHexConnection.Connection[1], true
	}
	return nil, -1, false
}

func (this *Game) drawPendingHex(screen *ebiten.Image, hoveredHex *hexagon.Hex) {
	if hoveredHex != nil {
		draw.Hexagon(screen, hoveredHex, this.theme.PendingHexBorderColor)
		this.drawPendingConnections(screen, hoveredHex)
	}
}

func (this *Game) drawCompletedLoops(screen *ebiten.Image) {
	draw.Loops(screen, this.loops, this.theme.CompletedLoopColor, this.theme.BackgroundColor)
}

func (this *Game) drawHighScore(screen *ebiten.Image) {
	text.Draw(screen, this.highScoreString(this.highScore), getTextFace(smallTextSize), getDrawHighScoreOptions(this.theme.ConnectionColor))
}

func (this *Game) drawGameScreen(screen *ebiten.Image) {
	this.drawScore(screen)
	this.drawHighScore(screen)
	this.drawHexagonGameBoard(screen)
	this.drawPlacedHexagons(screen)
	//this.drawCurrentHexPattern(screen)
	this.drawPendingHex(screen, this.getHoveredHex())
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
	if this.startButton.Hovered {
		draw.Hexagon(screen, this.startButton.Hex, this.theme.PendingHexBorderColor)
	}
	if this.tutorialButton.Hovered {
		draw.Hexagon(screen, this.tutorialButton.Hex, this.theme.PendingHexBorderColor)
	}
}

func (this *Game) drawNextArrow(screen *ebiten.Image, clr color.RGBA) {
	startX := float32(this.ScreenWidth - 70)
	startY := float32(marginSize + smallTextSize/2)
	unit := float32(smallTextSize / 3)
	this.drawArrow(screen, startX, startY, unit, clr)
}

func (this *Game) drawArrow(screen *ebiten.Image, startX, startY, unit float32, clr color.RGBA) {
	// vertical line down
	vector.StrokeLine(screen, startX, startY, startX, startY+(unit), draw.HexagonStrokeWidth, clr, true)
	// bottom line
	vector.StrokeLine(screen, startX, startY+unit, startX+(unit*4), startY+unit, draw.HexagonStrokeWidth, clr, true)
	// then down
	vector.StrokeLine(screen, startX+(unit*4), startY+unit, startX+(unit*4), startY+unit*2, draw.HexagonStrokeWidth, clr, true)
	// then diagonal
	vector.StrokeLine(screen, startX+(unit*4), startY+unit*2, startX+unit*7, startY, draw.HexagonStrokeWidth, clr, true)

	// vertical line up
	vector.StrokeLine(screen, startX, startY, startX, startY-unit, draw.HexagonStrokeWidth, clr, true)
	// top line
	vector.StrokeLine(screen, startX, startY-unit, startX+(unit*4), startY-unit, draw.HexagonStrokeWidth, clr, true)
	// then up
	vector.StrokeLine(screen, startX+(unit*4), startY-unit, startX+(unit*4), startY-unit*2, draw.HexagonStrokeWidth, clr, true)
	// then diagonal
	vector.StrokeLine(screen, startX+(unit*4), startY-unit*2, startX+unit*7, startY, draw.HexagonStrokeWidth, clr, true)

}

func (this *Game) drawTutorialText(screen *ebiten.Image, content string) {
	drawOptions := &text.DrawOptions{}
	drawOptions.GeoM.Translate(20, (marginSize+smallTextSize)/2)
	drawOptions.ColorScale.ScaleWithColor(this.theme.ConnectionColor)
	text.Draw(screen, content, getTextFace(smallTextSize), drawOptions)
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func (this *Game) updateGameScreen() {
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

func (this *Game) updateScore(points int) {
	this.score += points
}

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
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		this.updateClickedHex(mouseX, mouseY)
		this.updateScore(calculatePoints(this.loops))
	} else {
		this.updateHoveredHex(mouseX, mouseY)
	}
}

func (this *Game) updateTitleScreen() {
	mouseX, mouseY := ebiten.CursorPosition()
	if this.startButton.PointInHexagon(float64(mouseX), float64(mouseY)) {
		this.startButton.Hovered = true
	} else {
		this.startButton.Hovered = false
	}
	if this.tutorialButton.PointInHexagon(float64(mouseX), float64(mouseY)) {
		this.tutorialButton.Hovered = true
	} else {
		this.tutorialButton.Hovered = false
	}
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		if this.startButton.Hovered {
			this.currentSceneType = gameScreen
		}
		if this.tutorialButton.Hovered {
			this.currentSceneType = tutorialScreenExplanation
			//todo initialize scene functions
			originX := float64(this.ScreenWidth) / 2
			originY := float64(this.ScreenHeight) - (hexagon.HexVertexRadiusTest * 1.5)
			this.startButton = hexagon.NewTextHexagon(0, 0, originX, originY, hexagon.HexVertexRadiusTest, draw.TitleHexagonStrokeWidth, draw.TitleConnectionWidth, "Start", smallTextSize)
		}
	}
}

func (this *Game) updateNextArrow() {
	mouseX, mouseY := ebiten.CursorPosition()
	startX := float32(this.ScreenWidth - 70)
	startY := float32(marginSize + smallTextSize/2)
	if float32(mouseX) > startX &&
		float32(mouseX) < startX+(smallTextSize*7/3) &&
		float32(mouseY) < startY+smallTextSize &&
		float32(mouseY) > startY-smallTextSize {
		this.nextArrowHovered = true
	} else {
		this.nextArrowHovered = false
	}
}

func (this *Game) updateTutorialExplanationScreen() {
	mouseX, mouseY := ebiten.CursorPosition()
	if this.startButton.PointInHexagon(float64(mouseX), float64(mouseY)) {
		this.startButton.Hovered = true
	} else {
		this.startButton.Hovered = false
	}
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) && this.startButton.Hovered {
		this.currentSceneType = gameScreen
	}
}

func (this *Game) updateTutorial1Screen() {
	var checkHex *hexagon.Hex
	hexes := newHexes(rows, cols, hexagon.HexVertexRadius, draw.HexagonStrokeWidth, draw.ConnectionWidth, getGameBoardFirstHexCoordinate())
	for _, hex := range hexes {
		if hex.Row == 2 && hex.Col == 5 {
			hex.Connections = []hexagon.Connection{{0, 5}, {1, 3}, {2, 4}}
		}
		if hex.Row == 2 && hex.Col == 12 {
			hex.Connections = []hexagon.Connection{{0, 5}, {1, 2}, {3, 4}}
			checkHex = hex
		}
		if hex.Row == 2 && hex.Col == 6 {
			hex.Connections = []hexagon.Connection{{0, 1}, {2, 5}, {3, 4}}
		}
		// Top Left
		if hex.Row == 2 && hex.Col == 8 {
			hex.Connections = []hexagon.Connection{{0, 5}, {1, 4}, {2, 3}}
		}
		// Top Right
		if hex.Row == 2 && hex.Col == 10 {
			hex.Connections = []hexagon.Connection{{0, 3}, {1, 2}, {4, 5}}
		}
		// Left
		if hex.Row == 2 && hex.Col == 7 {
			hex.Connections = []hexagon.Connection{{0, 2}, {1, 3}, {4, 5}}
		}
		// Center
		if hex.Row == 2 && hex.Col == 9 {
			hex.Connections = []hexagon.Connection{{0, 4}, {1, 3}, {2, 5}}
		}
		// Right
		if hex.Row == 2 && hex.Col == 11 {
			hex.Connections = []hexagon.Connection{{0, 4}, {1, 2}, {3, 5}}
		}
		// Bottom Left
		if hex.Row == 3 && hex.Col == 8 {
			hex.Connections = []hexagon.Connection{{0, 5}, {1, 3}, {2, 4}}
		}
		// Bottom Right
		if hex.Row == 3 && hex.Col == 10 {
			hex.Connections = []hexagon.Connection{{0, 5}, {1, 2}, {3, 4}}
		}
	}
	this.loops = this.getCompleteLoops(checkHex)
	this.hexes = hexes
	this.updateNextArrow()
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) && this.nextArrowHovered {
		this.currentSceneType = tutorialScreen2
	}
}

// todo add numbers
func (this *Game) updateTutorial2Screen() {
	var checkHex *hexagon.Hex
	hexes := newHexes(rows, cols, hexagon.HexVertexRadius, draw.HexagonStrokeWidth, draw.ConnectionWidth, getGameBoardFirstHexCoordinate())
	for _, hex := range hexes {
		if hex.Row == 2 && hex.Col == 5 {
			hex.Connections = []hexagon.Connection{{0, 5}, {1, 3}, {2, 4}}
		}
		if hex.Row == 2 && hex.Col == 12 {
			hex.Connections = []hexagon.Connection{{0, 5}, {1, 2}, {3, 4}}
			checkHex = hex
		}
		if hex.Row == 2 && hex.Col == 6 {
			hex.Connections = []hexagon.Connection{{0, 1}, {2, 5}, {3, 4}}
		}
		// Top Left
		if hex.Row == 2 && hex.Col == 8 {
			hex.Connections = []hexagon.Connection{{0, 5}, {1, 4}, {2, 3}}
		}
		// Top Right
		if hex.Row == 2 && hex.Col == 10 {
			hex.Connections = []hexagon.Connection{{0, 3}, {1, 2}, {4, 5}}
		}
		// Left
		if hex.Row == 2 && hex.Col == 7 {
			hex.Connections = []hexagon.Connection{{0, 2}, {1, 3}, {4, 5}}
		}
		// Center
		if hex.Row == 2 && hex.Col == 9 {
			hex.Connections = []hexagon.Connection{{0, 4}, {1, 3}, {2, 5}}
		}
		// Right
		if hex.Row == 2 && hex.Col == 11 {
			hex.Connections = []hexagon.Connection{{0, 4}, {1, 2}, {3, 5}}
		}
		// Bottom Left
		if hex.Row == 3 && hex.Col == 8 {
			hex.Connections = []hexagon.Connection{{0, 5}, {1, 3}, {2, 4}}
		}
		// Bottom Right
		if hex.Row == 3 && hex.Col == 10 {
			hex.Connections = []hexagon.Connection{{0, 5}, {1, 2}, {3, 4}}
		}
	}
	this.loops = this.getCompleteLoops(checkHex)
	this.hexes = hexes
	this.updateNextArrow()
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) && this.nextArrowHovered {
		// TODO update tutorial screen
		this.currentSceneType = tutorialScreen2
	}
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

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

func (this *Game) getBorderHex(row, col, side int) *hexagon.Hex {
	r, c := getBorderHexGridPosition(row, col, side)
	return this.getHexFromGridPosition(r, c)
}

func (this *Game) getHexFromGridPosition(row, col int) *hexagon.Hex {
	for _, hex := range this.hexes {
		if hex.Row == row && hex.Col == col {
			return hex
		}
	}
	return nil
}

func (this *Game) getCompleteLoops(hex *hexagon.Hex) []hexagon.Loop {
	var loops []hexagon.Loop
	if !hex.Empty() {
		for _, connection := range hex.Connections {
			side := connection[0]
			startingHexConnection := hexagon.HexConnection{
				Hex:        hex,
				Connection: connection,
			}
			if duplicateConnection(loops, startingHexConnection) {
				continue
			}
			loop, completedLoop, _ := this.findLoop(
				side,
				hex,
				startingHexConnection,
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

func duplicateConnection(loops []hexagon.Loop, connection hexagon.HexConnection) bool {
	for _, loop := range loops {
		if loop.Contains(connection) {
			return true
		}
	}
	return false
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

func (this *Game) gameOver() bool {
	for _, hex := range this.hexes {
		if hex.Empty() {
			return false
		}
	}
	return this.disabledTicksLeft == 0
}

func (this *Game) startOver() {
	for _, hex := range this.hexes {
		hex.Reset()
	}
	this.score = 0
	this.loops = nil
}

func (this *Game) drawTutorialScreenExplanation(screen *ebiten.Image) {
	var clr color.RGBA
	if this.startButton.Hovered {
		clr = this.theme.PendingHexBorderColor
	} else {
		clr = this.theme.HexBorderColor
	}
	draw.TextHexagon(screen, this.startButton, clr, this.theme.ConnectionColor)

	drawOptions := &text.DrawOptions{}
	drawOptions.ColorScale.ScaleWithColor(this.theme.ConnectionColor)
	yTranslate := float64(marginSize+smallTextSize) / 2
	drawOptions.GeoM.Translate(20, yTranslate)
	text.Draw(screen, "The object of Hexloop is to place tiles to", getTextFace(smallTextSize), drawOptions)
	drawOptions.GeoM.Translate(0, yTranslate)
	text.Draw(screen, "connect loops. The tiles used to make a loop", getTextFace(smallTextSize), drawOptions)
	drawOptions.GeoM.Translate(0, yTranslate)
	text.Draw(screen, "will disappear to give you room to place more", getTextFace(smallTextSize), drawOptions)
	drawOptions.GeoM.Translate(0, yTranslate)
	text.Draw(screen, "tiles. You will also gets points for each", getTextFace(smallTextSize), drawOptions)
	drawOptions.GeoM.Translate(0, yTranslate)
	text.Draw(screen, "connection in the loop. The longer the loop,", getTextFace(smallTextSize), drawOptions)
	drawOptions.GeoM.Translate(0, yTranslate)
	text.Draw(screen, "the more points you'll get for each segment.", getTextFace(smallTextSize), drawOptions)
	drawOptions.GeoM.Translate(0, yTranslate)
	text.Draw(screen, "For example, a loop with 4 connections is", getTextFace(smallTextSize), drawOptions)
	drawOptions.GeoM.Translate(0, yTranslate)
	text.Draw(screen, "worth 10 points, but a loop with 8", getTextFace(smallTextSize), drawOptions)
	drawOptions.GeoM.Translate(0, yTranslate)
	text.Draw(screen, "connections is worth 36 points! If you make", getTextFace(smallTextSize), drawOptions)
	drawOptions.GeoM.Translate(0, yTranslate)
	text.Draw(screen, "two loops at once, you'll get double the", getTextFace(smallTextSize), drawOptions)
	drawOptions.GeoM.Translate(0, yTranslate)
	text.Draw(screen, "points for those loops. If you're able to get", getTextFace(smallTextSize), drawOptions)
	drawOptions.GeoM.Translate(0, yTranslate)
	text.Draw(screen, "every placed tile cleared, you'll get a 5,000", getTextFace(smallTextSize), drawOptions)
	drawOptions.GeoM.Translate(0, yTranslate)
	text.Draw(screen, "point bonus! How many points can you get?", getTextFace(smallTextSize), drawOptions)

}

func (this *Game) drawTutorialScreen1(screen *ebiten.Image) {
	this.drawHexagonGameBoard(screen)
	this.drawPlacedHexagons(screen)
	this.drawCompletedLoops(screen)
	this.drawTutorialText(screen, "The object of Hexloop is to make loops.")
	var clr color.RGBA
	if this.nextArrowHovered {
		clr = this.theme.PendingHexBorderColor
	} else {
		clr = this.theme.ConnectionColor
	}
	this.drawNextArrow(screen, clr)
}

func (this *Game) drawTutorialScreen2(screen *ebiten.Image) {
	this.drawHexagonGameBoard(screen)
	this.drawPlacedHexagons(screen)
	this.drawCompletedLoops(screen)
	this.drawTutorialText(screen, "This loop is 10 connections long...")
	var clr color.RGBA
	if this.nextArrowHovered {
		clr = this.theme.PendingHexBorderColor
	} else {
		clr = this.theme.ConnectionColor
	}
	this.drawNextArrow(screen, clr)
}

func (this *Game) highScoreString(score int) string {
	scoreStr := strconv.Itoa(score)
	return "High Score: " + withCommas(scoreStr)
}

func (this *Game) scoreString(score int) string {
	scoreStr := strconv.Itoa(score)
	return "Score: " + withCommas(scoreStr)
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

func getGameBoardFirstHexCoordinate() hexagon.Coordinate {
	xBuffer := marginSize + hexagon.HexSideRadius
	yBuffer := float64(marginSize + hexagon.HexVertexRadius + (smallTextSize * 2))
	return [2]float64{xBuffer, yBuffer}
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

func withCommas(s string) string {
	if len(s) <= 3 {
		return s
	}
	return withCommas(s[:len(s)-3]) + "," + s[len(s)-3:]
}
