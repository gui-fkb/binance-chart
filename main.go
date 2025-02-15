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
	barWidth     = 16
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

	game := &Game{}
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}

type Game struct{}

func (g *Game) Update() error {
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
	xOffset := screenWidth - len(Candlesticks)*barWidth

	const minCandleHeight = 1

	for i, c := range Candlesticks {
		x := xOffset + i*barWidth

		scaleFactor := screenHeight * 0.3 // Reduzir um pouco a escala
		margin := screenHeight * 0.1      // Adicionar uma margem no topo e base

		yOpen := int(screenHeight - margin - ((Atof(c.OpenPrice)-minPrice)/priceRange)*scaleFactor)
		yClose := int(screenHeight - margin - ((Atof(c.ClosePrice)-minPrice)/priceRange)*scaleFactor)
		yHigh := int(screenHeight - margin - ((Atof(c.HighPrice)-minPrice)/priceRange)*scaleFactor)
		yLow := int(screenHeight - margin - ((Atof(c.LowPrice)-minPrice)/priceRange)*scaleFactor)

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

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
