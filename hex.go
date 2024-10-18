package main

import "math"

const (
	sqrt3           = 1.7320508075688772
	hexVertexRadius = 30 // Distance from center of hexagon to vertex
	hexSideRadius   = hexVertexRadius * sqrt3 / 2
	numHexagonSides = 6
)

type Connection [2]int
type Coordinate [2]float64

// TODO Make connections a different color
// maybe black if it hits a wall?
// maybe you hover the connection to make it glow?

// Hex represents a hexagonal tile
type Hex struct {
	col, row    int        // Column and Row
	center      Coordinate // center of the hex
	connections []Connection
	hovered     bool
}

func (this *Hex) available() bool {
	return len(this.connections) == 0
}

// vertexCoordinates returns a slice of coordinates of each vertex of the hexagon
func (this *Hex) vertexCoordinates() []Coordinate {
	vertices := make([]Coordinate, numHexagonSides)
	for i := 0; i < numHexagonSides; i++ {
		angle := math.Pi/3*float64(i) - math.Pi/6
		x := getXCoordinateFromPolar(this.center[0], hexVertexRadius, angle)
		y := getYCoordinateFromPolar(this.center[1], hexVertexRadius, angle)
		vertices[i] = Coordinate{x, y}
	}
	return vertices
}

// hexagonSideCoordinates returns a slice of coordinates of the midpoint of each side of the hexagon
func (this *Hex) hexagonSideCoordinates() []Coordinate {
	sides := make([]Coordinate, numHexagonSides)
	for i := 0; i < numHexagonSides; i++ {
		angle := math.Pi / 3 * float64(i-1)
		x := getXCoordinateFromPolar(this.center[0], hexSideRadius, angle)
		y := getYCoordinateFromPolar(this.center[1], hexSideRadius, angle)
		sides[i] = Coordinate{x, y}
	}
	return sides
}

// pointInHexagon checks if a point is inside the hexagon
func (this *Hex) pointInHexagon(px, py, radius float64) bool {
	buffer := .1 // prevents two hexagons being selected at once
	dx := math.Abs(px-this.center[0]) / radius
	dy := math.Abs(py-this.center[1]) / radius
	return dx <= 1.0-buffer &&
		dy <= math.Sqrt(3.0)/2.0-buffer &&
		dx+dy/math.Sqrt(3.0) <= 1.0-buffer
}

func getXCoordinateFromPolar(centerX, radius, angle float64) float64 {
	return centerX + radius*math.Cos(angle)
}

func getYCoordinateFromPolar(centerY, radius, angle float64) float64 {
	return centerY + radius*math.Sin(angle)
}
