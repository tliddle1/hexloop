package color

import (
	"image/color"
	"strconv"
)

type Theme struct {
	BackgroundColor         color.RGBA
	HexBorderColor          color.RGBA
	ConnectionColor         color.RGBA
	PendingHexBorderColor   color.RGBA
	PendingConnectionColors []color.RGBA
}

func NewBeeTheme() *Theme {
	return &Theme{
		BackgroundColor:       color.RGBA{R: 251, G: 217, B: 100, A: 255}, // Bee Yellow
		HexBorderColor:        color.RGBA{R: 254, G: 237, B: 161, A: 255}, // Beige
		ConnectionColor:       color.RGBA{R: 133, G: 77, B: 13, A: 255},   // Brown
		PendingHexBorderColor: color.RGBA{R: 150, G: 123, B: 182, A: 255}, // Lavender
		PendingConnectionColors: []color.RGBA{
			{R: 0, G: 255, B: 0, A: 255},   // Green
			{R: 0, G: 0, B: 255, A: 255},   // Blue
			{R: 255, G: 0, B: 255, A: 255}, // Magenta
		},
	}
}
func NewBlueTheme() *Theme {
	return &Theme{
		BackgroundColor:       HexToRGB("#53687E"), // Payne's Gray
		HexBorderColor:        HexToRGB("#3A4454"), // Charcoal
		ConnectionColor:       HexToRGB("#C2B2B4"), // French Gray
		PendingHexBorderColor: HexToRGB("#F5DDDD"), // Misty Rose
		PendingConnectionColors: []color.RGBA{
			{R: 0, G: 255, B: 0, A: 255},   // Green
			{R: 255, G: 165, B: 0, A: 255}, // Orange
			{R: 255, G: 0, B: 255, A: 255}, // Magenta
		},
	}
}
func NewDefaultTheme() *Theme {
	return NewBeeTheme()
}

// HexToRGB converts a 6-digit hex color code to RGB values.
func HexToRGB(hexCode string) color.RGBA {
	white := color.RGBA{R: 255, G: 255, B: 255, A: 255}
	if len(hexCode) > 6 {
		hexCode = hexCode[len(hexCode)-6:]
	}
	if len(hexCode) != 6 {
		return white
	}

	// Convert the hex string to its RGB components
	r, _ := strconv.ParseInt(hexCode[0:2], 16, 32)
	g, _ := strconv.ParseInt(hexCode[2:4], 16, 32)
	b, _ := strconv.ParseInt(hexCode[4:6], 16, 32)

	return color.RGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: 255}
}
