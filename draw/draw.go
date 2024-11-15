package draw

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	color2 "github.com/tliddle1/hexloop/color"
	"github.com/tliddle1/hexloop/hexagon"
	"github.com/tliddle1/hexloop/vector"
)

const (
	HexagonStrokeWidth      = 2
	ConnectionWidth         = 3
	TitleHexagonStrokeWidth = 4
	TitleConnectionWidth    = 6
	pathBuffer              = 4
)

// TODO meet corners of hexagon visually

func Hexagon(screen *ebiten.Image, hex *hexagon.Hex, borderColor color.RGBA) {
	vertices := hex.VertexCoordinates()
	vertices = append(vertices, vertices[0])
	for i := 0; i < hexagon.NumHexagonSides; i++ {
		vector.StrokeLine(screen,
			float32(vertices[i][0]),
			float32(vertices[i][1]),
			float32(vertices[i+1][0]),
			float32(vertices[i+1][1]),
			hex.EdgeWidth,
			borderColor,
			true)
	}
}

func TextHexagon(screen *ebiten.Image, hex *hexagon.TextHexagon, borderColor, connectionColor color.RGBA) {
	Hexagon(screen, hex.Hex, borderColor)
	x, y := hex.Center[0], hex.Center[1]
	fontSize := hex.TextSize
	face := hex.GetTextFace(fontSize)
	// todo refactor
	if hex.Str == "How to Play" {
		offset := text.Advance("to", face)
		drawOptions := &text.DrawOptions{}
		drawOptions.ColorScale.ScaleWithColor(connectionColor)
		// to
		drawOptions.GeoM.Translate(x-(offset/2), y-(fontSize*2/3))
		text.Draw(screen, "to", face, drawOptions)
		// how
		previousXOffset := x - (offset / 2)
		offset = text.Advance("How", face)
		drawOptions.GeoM.Translate(x-(offset/2)-previousXOffset, -fontSize)
		text.Draw(screen, "How", face, drawOptions)
		// play
		previousXOffset = x - (offset / 2)
		offset = text.Advance("Play", face)
		drawOptions.GeoM.Translate(x-(offset/2)-previousXOffset, fontSize*2)
		text.Draw(screen, "Play", face, drawOptions)
	} else {
		offset := text.Advance(hex.Str, face)
		drawOptions := &text.DrawOptions{}
		drawOptions.GeoM.Translate(x-(offset/2), y-(fontSize*2/3))
		drawOptions.ColorScale.ScaleWithColor(connectionColor)
		text.Draw(screen, hex.Str, face, drawOptions)
	}
}

func HexagonConnections(screen *ebiten.Image, hex *hexagon.Hex, connectionColor color.RGBA, theme *color2.Theme) {
	if len(hex.Connections) == 0 {
		return
	}
	for _, connection := range hex.Connections {
		HexagonConnection(screen, hex, connection, connectionColor, theme.BackgroundColor)
	}
}

func HexagonConnection(screen *ebiten.Image, hex *hexagon.Hex, connection hexagon.Connection, connectionColor, backgroundColor color.RGBA) {
	sideA := connection[0]
	sideB := connection[1]
	diff := math.Abs(float64(sideA - sideB))
	// Straight Across
	if diff == 3 {
		LineConnection(screen, hex, connection, hex.ConnectionWidth+pathBuffer, backgroundColor)
		LineConnection(screen, hex, connection, hex.ConnectionWidth, connectionColor)
		return
	}
	// Large Curve
	if diff == 2 {
		centerSide := (sideA + sideB) / 2
		angleToSide := math.Pi*2/3 + math.Pi*1/3*float64(centerSide)
		LargeCurveConnection(screen, hex, angleToSide, hex.ConnectionWidth+pathBuffer, backgroundColor)
		LargeCurveConnection(screen, hex, angleToSide, hex.ConnectionWidth, connectionColor)
		return
	}
	if diff == 4 {
		oppositeCenterSide := (sideA + sideB) / 2
		angleToSide := -math.Pi*1/3 + math.Pi*1/3*float64(oppositeCenterSide)
		LargeCurveConnection(screen, hex, angleToSide, hex.ConnectionWidth+pathBuffer, backgroundColor)
		LargeCurveConnection(screen, hex, angleToSide, hex.ConnectionWidth, connectionColor)
		return
	}
	// Small Curve
	if diff == 1 {
		SmallCurveConnection(screen, hex, min(sideA, sideB), hex.ConnectionWidth+pathBuffer, backgroundColor)
		SmallCurveConnection(screen, hex, min(sideA, sideB), hex.ConnectionWidth, connectionColor)
		return
	}
	if diff == 5 {
		SmallCurveConnection(screen, hex, 5, hex.ConnectionWidth+pathBuffer, backgroundColor)
		SmallCurveConnection(screen, hex, 5, hex.ConnectionWidth, connectionColor)
		return
	}
}

func LineConnection(screen *ebiten.Image, hex *hexagon.Hex, connection hexagon.Connection, strokeWidth float32, connectionColor color.RGBA) {
	sideA := connection[0]
	sideB := connection[1]
	sides := hex.HexagonSideCoordinates()
	x1 := float32(sides[sideA][0])
	y1 := float32(sides[sideA][1])
	x2 := float32(sides[sideB][0])
	y2 := float32(sides[sideB][1])
	vector.StrokeLine(screen, x1, y1, x2, y2, strokeWidth, connectionColor, true)
}

func LargeCurveConnection(screen *ebiten.Image, hex *hexagon.Hex, angleToSide float64, strokeWidth float32, connectionColor color.Color) {
	x := float32(hex.Center[0] - math.Cos(angleToSide)*hex.SideRadius*2)
	y := float32(hex.Center[1] - math.Sin(angleToSide)*hex.SideRadius*2)
	radius := float32(hex.VertexRadius+hex.SideRadius/2) + hex.EdgeWidth
	startAngle := float32(angleToSide - math.Pi/6)
	endAngle := float32(angleToSide + math.Pi/6)
	vector.StrokePartialCircle(screen, x, y, radius, startAngle, endAngle, strokeWidth, connectionColor, true)
}

func SmallCurveConnection(screen *ebiten.Image, hex *hexagon.Hex, vertex int, strokeWidth float32, connectionColor color.RGBA) {
	vertices := hex.VertexCoordinates()
	x := float32(vertices[vertex][0])
	y := float32(vertices[vertex][1])
	adjustor := math.Pi / 3 * float32(vertex)
	vector.StrokePartialCircle(screen, x, y, float32(hex.VertexRadius/2), math.Pi/2+adjustor, -math.Pi*5/6+adjustor, strokeWidth, connectionColor, true)
}

func Loops(screen *ebiten.Image, loops []hexagon.Loop, completedLoopColor, backgroundColor color.RGBA) {
	for _, loop := range loops {
		for _, hexConnection := range loop {
			HexagonConnection(screen, hexConnection.Hex, hexConnection.Connection, completedLoopColor, backgroundColor)
		}
	}
}
