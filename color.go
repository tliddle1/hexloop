package main

import "image/color"

// TODO make game settings to pass these values in
// TODO different color themes
var (
	backgroundColor        = color.RGBA{R: 251, G: 217, B: 100, A: 255} // Yellow
	hexBorderColor         = color.RGBA{R: 254, G: 237, B: 161, A: 255} // Beige
	connectionColor        = color.RGBA{R: 133, G: 77, B: 13, A: 255}   // Brown
	pendingConnectionColor = color.RGBA{R: 150, G: 123, B: 182, A: 255} // Lavender
)
