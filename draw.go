package main

import (
	colors "image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	color2 "github.com/tliddle1/game/color"
	"github.com/tliddle1/game/vector"
)

const (
	hexagonStrokeWidth    = 2
	connectionStrokeWidth = 3
	pathBuffer            = 4
)

func drawHexagon(screen *ebiten.Image, hex *Hex, color colors.RGBA) {
	vertices := hex.vertexCoordinates()
	vertices = append(vertices, vertices[0])
	for i := 0; i < numHexagonSides; i++ {
		vector.StrokeLine(screen,
			float32(vertices[i][0]),
			float32(vertices[i][1]),
			float32(vertices[i+1][0]),
			float32(vertices[i+1][1]),
			hexagonStrokeWidth,
			color,
			true)
	}
}

func drawHexagonConnections(screen *ebiten.Image, hex Hex, theme *color2.Theme) {
	if len(hex.connections) == 0 {
		return
	}
	for _, connection := range hex.connections {
		drawHexagonConnection(screen, hex, connection, theme.ConnectionColor, theme.BackgroundColor)
	}
}

//func drawPendingHexagonConnection(screen *ebiten.Image, hex Hex, nextConnections []Connection, theme *Theme) {
//	if len(hex.connections) != 0 {
//		return
//	}
//	for i, connection := range nextConnections {
//		drawHexagonConnection(screen, hex, connection, pendingColors[i])
//	}
//}

func drawHexagonConnection(screen *ebiten.Image, hex Hex, connection Connection, color, backgroundColor colors.RGBA) {
	sideA := connection[0]
	sideB := connection[1]
	diff := math.Abs(float64(sideA - sideB))
	// Straight Across
	if diff == 3 {
		drawLineConnection(screen, hex, connection, connectionStrokeWidth+pathBuffer, backgroundColor)
		drawLineConnection(screen, hex, connection, connectionStrokeWidth, color)
		return
	}
	// Large Curve
	if diff == 2 {
		centerSide := (sideA + sideB) / 2
		angleToSide := math.Pi*2/3 + math.Pi*1/3*float64(centerSide)
		drawLargeCurveConnection(screen, hex, angleToSide, connectionStrokeWidth+pathBuffer, backgroundColor)
		drawLargeCurveConnection(screen, hex, angleToSide, connectionStrokeWidth, color)
		return
	}
	if diff == 4 {
		oppositeCenterSide := (sideA + sideB) / 2
		angleToSide := -math.Pi*1/3 + math.Pi*1/3*float64(oppositeCenterSide)
		drawLargeCurveConnection(screen, hex, angleToSide, connectionStrokeWidth+pathBuffer, backgroundColor)
		drawLargeCurveConnection(screen, hex, angleToSide, connectionStrokeWidth, color)
		return
	}
	// Small Curve
	if diff == 1 {
		drawSmallCurveConnection(screen, hex.vertexCoordinates(), min(sideA, sideB), connectionStrokeWidth+pathBuffer, backgroundColor)
		drawSmallCurveConnection(screen, hex.vertexCoordinates(), min(sideA, sideB), connectionStrokeWidth, color)
		return
	}
	if diff == 5 {
		drawSmallCurveConnection(screen, hex.vertexCoordinates(), 5, connectionStrokeWidth+pathBuffer, backgroundColor)
		drawSmallCurveConnection(screen, hex.vertexCoordinates(), 5, connectionStrokeWidth, color)
		return
	}
}

func drawLineConnection(screen *ebiten.Image, hex Hex, connection Connection, strokeWidth float32, color colors.RGBA) {
	sideA := connection[0]
	sideB := connection[1]
	sides := hex.hexagonSideCoordinates()
	x1 := float32(sides[sideA][0])
	y1 := float32(sides[sideA][1])
	x2 := float32(sides[sideB][0])
	y2 := float32(sides[sideB][1])
	vector.StrokeLine(screen, x1, y1, x2, y2, strokeWidth, color, true)
}

func drawLargeCurveConnection(screen *ebiten.Image, hex Hex, angleToSide float64, strokeWidth float32, color colors.Color) {
	x := float32(hex.center[0] - math.Cos(angleToSide)*hexSideRadius*2)
	y := float32(hex.center[1] - math.Sin(angleToSide)*hexSideRadius*2)
	radius := float32(hexVertexRadius + hexSideRadius/2 + hexagonStrokeWidth)
	startAngle := float32(angleToSide - math.Pi/6)
	endAngle := float32(angleToSide + math.Pi/6)
	vector.StrokePartialCircle(screen, x, y, radius, startAngle, endAngle, strokeWidth, color, true)
}

func drawSmallCurveConnection(screen *ebiten.Image, vertices []Coordinate, vertex int, strokeWidth float32, color colors.RGBA) {
	x := float32(vertices[vertex][0])
	y := float32(vertices[vertex][1])
	adjustor := math.Pi / 3 * float32(vertex)
	vector.StrokePartialCircle(screen, x, y, hexVertexRadius/2, math.Pi/2+adjustor, -math.Pi*5/6+adjustor, strokeWidth, color, true)
}
