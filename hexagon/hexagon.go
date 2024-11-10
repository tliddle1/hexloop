package hexagon

import (
	"bytes"
	"log"
	"math"

	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

const (
	sqrt3               = 1.7320508075688772
	HexVertexRadiusTest = 60
	HexSideRadiusTest   = HexVertexRadiusTest * sqrt3 / 2
	HexVertexRadius     = 30 // Distance from center of hexagon to vertex
	HexSideRadius       = HexVertexRadius * sqrt3 / 2
	NumHexagonSides     = 6
)

type Connection [2]int
type Coordinate [2]float64

// TODO hover over the connection to make it glow
// maybe hover over the hex to make its three Connections glow?
// maybe black if it hits a wall?

// Hex represents a hexagonal tile
type Hex struct {
	Col, Row                   int // Column and Row
	VertexRadius, SideRadius   float64
	EdgeWidth, ConnectionWidth float32
	Center                     Coordinate // center of the hex
	Connections                []Connection
	Hovered                    bool
}

func NewHex(col, row int, originX, originY, hexVertexRadius float64, edgeWidth, connectionWidth float32) *Hex {
	hexSideRadius := hexVertexRadius * sqrt3 / 2
	x := float64(col) * hexSideRadius
	y := float64(row) * hexVertexRadius * 3
	// Offset odd rows to create staggered effect
	if col%2 != 0 {
		y += hexVertexRadius * 1.5
	}
	return &Hex{
		Col:             col,
		Row:             row,
		VertexRadius:    hexVertexRadius,
		EdgeWidth:       edgeWidth,
		ConnectionWidth: connectionWidth,
		SideRadius:      hexSideRadius,
		Center:          Coordinate{x + originX, y + originY},
	}
}

func (this *Hex) Empty() bool {
	return len(this.Connections) == 0
}

// VertexCoordinates returns a slice of coordinates of each vertex of the hexagon
func (this *Hex) VertexCoordinates() []Coordinate {
	vertices := make([]Coordinate, NumHexagonSides)
	for i := 0; i < NumHexagonSides; i++ {
		angle := math.Pi/3*float64(i) - math.Pi/6
		x := getXCoordinateFromPolar(this.Center[0], this.VertexRadius, angle)
		y := getYCoordinateFromPolar(this.Center[1], this.VertexRadius, angle)
		vertices[i] = Coordinate{x, y}
	}
	return vertices
}

// HexagonSideCoordinates returns a slice of coordinates of the midpoint of each side of the hexagon
func (this *Hex) HexagonSideCoordinates() []Coordinate {
	sides := make([]Coordinate, NumHexagonSides)
	for i := 0; i < NumHexagonSides; i++ {
		angle := math.Pi / 3 * float64(i-1)
		x := getXCoordinateFromPolar(this.Center[0], this.SideRadius, angle)
		y := getYCoordinateFromPolar(this.Center[1], this.SideRadius, angle)
		sides[i] = Coordinate{x, y}
	}
	return sides
}

// PointInHexagon checks if a point is inside the hexagon
func (this *Hex) PointInHexagon(px, py float64) bool {
	buffer := .1 // prevents two hexagons being selected at once
	dx := math.Abs(px-this.Center[0]) / this.VertexRadius
	dy := math.Abs(py-this.Center[1]) / this.VertexRadius
	return dx <= 1.0-buffer &&
		dy <= math.Sqrt(3.0)/2.0-buffer &&
		dx+dy/math.Sqrt(3.0) <= 1.0-buffer
}

// ConnectedSide returns the other half of the connection
// (e.g. if [0,2] is one of the Connections and 2 is passed in then 0 is returned)
func (this *Hex) ConnectedSide(side int) (connectedSide int) {
	if len(this.Connections) == 0 {
		panic("why are you looking for the connectedSide of an empty hexagon?")
	}
	for _, connection := range this.Connections {
		if connection[0] == side {
			return connection[1]
		}
		if connection[1] == side {
			return connection[0]
		}
	}
	return -1
}

func (this *Hex) Equals(hex *Hex) bool {
	if this.Row != hex.Row || this.Col != hex.Col {
		return false
	}
	return true
}

func (this *Hex) Reset() {
	this.Hovered = false
	this.Connections = nil
}

func getXCoordinateFromPolar(centerX, radius, angle float64) float64 {
	return centerX + radius*math.Cos(angle)
}

func getYCoordinateFromPolar(centerY, radius, angle float64) float64 {
	return centerY + radius*math.Sin(angle)
}

type HexConnection struct {
	Hex        *Hex
	Connection Connection
}

type Loop []HexConnection

type TextHexagon struct {
	*Hex
	Str      string
	TextSize float64
}

func (this *TextHexagon) GetTextFace(textSize float64) text.Face {
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

func NewTextHexagon(col, row int, originX, originY, hexVertexRadius float64, edgeWidth, connectionWidth float32, str string, testSize float64) *TextHexagon {
	return &TextHexagon{
		Hex:      NewHex(col, row, originX, originY, hexVertexRadius, edgeWidth, connectionWidth),
		Str:      str,
		TextSize: testSize,
	}
}
