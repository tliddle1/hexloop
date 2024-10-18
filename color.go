package main

import (
	"image/color"
	"strconv"
)

type Theme struct {
	BackgroundColor        color.RGBA
	HexBorderColor         color.RGBA
	ConnectionColor        color.RGBA
	PendingConnectionColor color.RGBA
}

func NewBeeTheme() *Theme {
	return &Theme{
		BackgroundColor:        color.RGBA{R: 251, G: 217, B: 100, A: 255}, // Yellow
		HexBorderColor:         color.RGBA{R: 254, G: 237, B: 161, A: 255}, // Beige
		ConnectionColor:        color.RGBA{R: 133, G: 77, B: 13, A: 255},   // Brown
		PendingConnectionColor: color.RGBA{R: 150, G: 123, B: 182, A: 255}, // Lavender
	}
}
func NewBlueTheme() *Theme {
	return &Theme{
		BackgroundColor:        HexToRGB("#53687E"),
		HexBorderColor:         HexToRGB("#3A4454"),
		ConnectionColor:        HexToRGB("#C2B2B4"),
		PendingConnectionColor: HexToRGB("#6B4E71"),
	}
}
func NewDefaultTheme() *Theme {
	return NewBeeTheme()
	//return &Theme{
	//	BackgroundColor:        HexToRGB(""),
	//	HexBorderColor:         HexToRGB(""),
	//	ConnectionColor:        HexToRGB(""),
	//	PendingConnectionColor: HexToRGB(""),
	//}
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
