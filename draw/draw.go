package draw

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	color2 "github.com/tliddle1/hexloop/color"
	"github.com/tliddle1/hexloop/hexagon"
	"github.com/tliddle1/hexloop/vector"
)

const (
	hexagonStrokeWidth    = 2
	connectionStrokeWidth = 3
	pathBuffer            = 4
)

func Hexagon(screen *ebiten.Image, hex *hexagon.Hex, borderColor color.RGBA) {
	vertices := hex.VertexCoordinates()
	vertices = append(vertices, vertices[0])
	for i := 0; i < hexagon.NumHexagonSides; i++ {
		vector.StrokeLine(screen,
			float32(vertices[i][0]),
			float32(vertices[i][1]),
			float32(vertices[i+1][0]),
			float32(vertices[i+1][1]),
			hexagonStrokeWidth,
			borderColor,
			true)
	}
}

func HexagonConnections(screen *ebiten.Image, hex hexagon.Hex, theme *color2.Theme) {
	if len(hex.Connections) == 0 {
		return
	}
	for _, connection := range hex.Connections {
		HexagonConnection(screen, hex, connection, theme.ConnectionColor, theme.BackgroundColor)
	}
}

func HexagonConnection(screen *ebiten.Image, hex hexagon.Hex, connection hexagon.Connection, connectionColor, backgroundColor color.RGBA) {
	sideA := connection[0]
	sideB := connection[1]
	diff := math.Abs(float64(sideA - sideB))
	// Straight Across
	if diff == 3 {
		LineConnection(screen, hex, connection, connectionStrokeWidth+pathBuffer, backgroundColor)
		LineConnection(screen, hex, connection, connectionStrokeWidth, connectionColor)
		return
	}
	// Large Curve
	if diff == 2 {
		centerSide := (sideA + sideB) / 2
		angleToSide := math.Pi*2/3 + math.Pi*1/3*float64(centerSide)
		LargeCurveConnection(screen, hex, angleToSide, connectionStrokeWidth+pathBuffer, backgroundColor)
		LargeCurveConnection(screen, hex, angleToSide, connectionStrokeWidth, connectionColor)
		return
	}
	if diff == 4 {
		oppositeCenterSide := (sideA + sideB) / 2
		angleToSide := -math.Pi*1/3 + math.Pi*1/3*float64(oppositeCenterSide)
		LargeCurveConnection(screen, hex, angleToSide, connectionStrokeWidth+pathBuffer, backgroundColor)
		LargeCurveConnection(screen, hex, angleToSide, connectionStrokeWidth, connectionColor)
		return
	}
	// Small Curve
	if diff == 1 {
		SmallCurveConnection(screen, hex.VertexCoordinates(), min(sideA, sideB), connectionStrokeWidth+pathBuffer, backgroundColor)
		SmallCurveConnection(screen, hex.VertexCoordinates(), min(sideA, sideB), connectionStrokeWidth, connectionColor)
		return
	}
	if diff == 5 {
		SmallCurveConnection(screen, hex.VertexCoordinates(), 5, connectionStrokeWidth+pathBuffer, backgroundColor)
		SmallCurveConnection(screen, hex.VertexCoordinates(), 5, connectionStrokeWidth, connectionColor)
		return
	}
}

func LineConnection(screen *ebiten.Image, hex hexagon.Hex, connection hexagon.Connection, strokeWidth float32, connectionColor color.RGBA) {
	sideA := connection[0]
	sideB := connection[1]
	sides := hex.HexagonSideCoordinates()
	x1 := float32(sides[sideA][0])
	y1 := float32(sides[sideA][1])
	x2 := float32(sides[sideB][0])
	y2 := float32(sides[sideB][1])
	vector.StrokeLine(screen, x1, y1, x2, y2, strokeWidth, connectionColor, true)
}

func LargeCurveConnection(screen *ebiten.Image, hex hexagon.Hex, angleToSide float64, strokeWidth float32, connectionColor color.Color) {
	x := float32(hex.Center[0] - math.Cos(angleToSide)*hexagon.HexSideRadius*2)
	y := float32(hex.Center[1] - math.Sin(angleToSide)*hexagon.HexSideRadius*2)
	radius := float32(hexagon.HexVertexRadius + hexagon.HexSideRadius/2 + hexagonStrokeWidth)
	startAngle := float32(angleToSide - math.Pi/6)
	endAngle := float32(angleToSide + math.Pi/6)
	vector.StrokePartialCircle(screen, x, y, radius, startAngle, endAngle, strokeWidth, connectionColor, true)
}

func SmallCurveConnection(screen *ebiten.Image, vertices []hexagon.Coordinate, vertex int, strokeWidth float32, connectionColor color.RGBA) {
	x := float32(vertices[vertex][0])
	y := float32(vertices[vertex][1])
	adjustor := math.Pi / 3 * float32(vertex)
	vector.StrokePartialCircle(screen, x, y, hexagon.HexVertexRadius/2, math.Pi/2+adjustor, -math.Pi*5/6+adjustor, strokeWidth, connectionColor, true)
}

func Loops(screen *ebiten.Image, loops [][]hexagon.HexConnection, completedLoopColor, backgroundColor color.RGBA) {
	for _, loop := range loops {
		for _, hexConnection := range loop {
			HexagonConnection(screen, *hexConnection.Hex, hexConnection.Connection, completedLoopColor, backgroundColor)
		}
	}
}
