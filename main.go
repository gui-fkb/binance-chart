package main

import (
	"image/color"
	"log"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	screenWidth  = 800
	screenHeight = 600
	barWidth     = 24
)

var (
	green = color.RGBA{0, 255, 0, 255} // Bullish
	red   = color.RGBA{255, 0, 0, 255} // Bearish
)

func main() {
	// Start WebSocket goroutine
	go startWebsocket()

	// Run Ebiten game loop
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("BTC/USDT Candlestick Chart")

	game := &Game{scaleFactor: 1, offsetX: 0, offsetY: 0}
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}

type Game struct {
	scaleFactor float64
	offsetX     int
	offsetY     int
	lastMouseX  int
	lastMouseY  int
	dragging    bool
}

func (g *Game) Update() error {
	// Zoom in/out using mouse wheel
	_, dy := ebiten.Wheel()
	if dy > 0 {
		g.scaleFactor *= 1.1 // Zoom in
	} else if dy < 0 {
		g.scaleFactor /= 1.1 // Zoom out
	}

	// Panning using arrow keys
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		g.offsetX += 5 // Reduced speed
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		g.offsetX -= 5
	}
	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		g.offsetY += 5
	}
	if ebiten.IsKeyPressed(ebiten.KeyDown) {
		g.offsetY -= 5
	}

	// Mouse dragging
	x, y := ebiten.CursorPosition()
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		if !g.dragging {
			// Start dragging, store initial position
			g.dragging = true
			g.lastMouseX = x
			g.lastMouseY = y
		} else {
			// Calculate difference and apply a damping factor
			g.offsetX += (x - g.lastMouseX) / 2 // Reduce sensitivity
			g.offsetY += (y - g.lastMouseY) / 2
			g.lastMouseX = x
			g.lastMouseY = y
		}
	} else {
		g.dragging = false
	}

	if ebiten.IsKeyPressed(ebiten.KeyR) { // Press 'R' to reset
		g.scaleFactor = 1.0
		g.offsetX = 0
		g.offsetY = 0
	}

	return nil
}
func (g *Game) Draw(screen *ebiten.Image) {
	if len(Candlesticks) == 0 {
		return
	}

	// Find price range
	minPrice, maxPrice := Atof(Candlesticks[0].LowPrice), Atof(Candlesticks[0].HighPrice)
	for _, c := range Candlesticks {
		cLow := Atof(c.LowPrice)
		cHigh := Atof(c.HighPrice)

		if cLow < minPrice {
			minPrice = cLow
		}
		if cHigh > maxPrice {
			maxPrice = cHigh
		}
	}

	priceRange := maxPrice - minPrice
	if priceRange == 0 {
		priceRange = 1 // Avoid division by zero
	}

	// Calculate spacing
	xOffset := screenWidth - len(Candlesticks)*barWidth + g.offsetX

	var minCandleHeight = 1.0 * g.scaleFactor

	for i, c := range Candlesticks {
		x := xOffset + i*barWidth

		scaleFactor := (screenHeight * 0.3) * g.scaleFactor // Apply zoom
		margin := (screenHeight * 0.1) * g.scaleFactor      // Scale margin

		yOpen := (screenHeight - margin - ((Atof(c.OpenPrice)-minPrice)/priceRange)*scaleFactor) + float64(g.offsetY)
		yClose := (screenHeight - margin - ((Atof(c.ClosePrice)-minPrice)/priceRange)*scaleFactor) + float64(g.offsetY)
		yHigh := (screenHeight - margin - ((Atof(c.HighPrice)-minPrice)/priceRange)*scaleFactor) + float64(g.offsetY)
		yLow := (screenHeight - margin - ((Atof(c.LowPrice)-minPrice)/priceRange)*scaleFactor) + float64(g.offsetY)

		if abs(yOpen-yClose) < minCandleHeight {
			if yOpen > yClose {
				yOpen = yClose + minCandleHeight
			} else {
				yClose = yOpen + minCandleHeight
			}
		}

		// Determine color
		candleColor := red
		if Atof(c.ClosePrice) > Atof(c.OpenPrice) {
			candleColor = green
		}

		// Draw candlestick body
		vector.DrawFilledRect(screen, float32(x), float32(yOpen), float32(barWidth), float32(yClose-yOpen), candleColor, true)

		// Draw wick
		vector.StrokeLine(screen, float32(x+barWidth/2), float32(yHigh), float32(x+barWidth/2), float32(yLow), 1, candleColor, true)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

// string to float
func Atof(s string) float64 {
	f, _ := strconv.ParseFloat(s, 64)
	return f
}

func Atoi(s string) int {
	i, _ := strconv.Atoi(s)
	return i
}

func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}
