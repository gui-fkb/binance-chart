package main

import (
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/fatih/color"
	"github.com/gorilla/websocket"
)

const wsBaseUrl = "wss://stream.binance.com:9443"

type BinanceStream struct {
	Stream string      `json:"stream"`
	Data   BinanceData `json:"data"`
}

type BinanceData struct {
	EventType string       `json:"e"`
	EventTime int64        `json:"E"`
	Symbol    string       `json:"s"`
	KLine     BinanceKline `json:"k"`
}

type BinanceKline struct {
	StartTime           int64  `json:"t"` // Kline start time
	CloseTime           int64  `json:"T"` // Kline close time
	Symbol              string `json:"s"` // Symbol
	Interval            string `json:"i"` // Interval
	FirstTradeID        int    `json:"f"` // First trade ID
	LastTradeID         int    `json:"L"` // Last trade ID
	OpenPrice           string `json:"o"` // Open price
	ClosePrice          string `json:"c"` // Close price
	HighPrice           string `json:"h"` // High price
	LowPrice            string `json:"l"` // Low price
	BaseAssetVolume     string `json:"v"` // Base asset volume
	NumberOfTrades      int    `json:"n"` // Number of trades
	IsKlineClosed       bool   `json:"x"` // Is this kline closed?
	QuoteAssetVolume    string `json:"q"` // Quote asset volume
	TakerBuyBaseVolume  string `json:"V"` // Taker buy base asset volume
	TakerBuyQuoteVolume string `json:"Q"` // Taker buy quote asset volume
	Ignore              string `json:"B"` // Ignore
}

var Candlesticks []BinanceKline

func startWebsocket() {
	conn, _, err := websocket.DefaultDialer.Dial(fmt.Sprintf("%s/stream?streams=btcusdt@kline_1m", wsBaseUrl), nil)
	if err != nil {
		slog.Error("could not connect to the websocket server")
	}
	defer conn.Close()

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			slog.Error("could not read message from the websocket server")
		}

		//fmt.Println(string(msg))

		var stream BinanceStream
		json.Unmarshal(msg, &stream)

		if stream.Data.EventType == "kline" {
			fmt.Println(string(msg))
			processKline(stream.Data.KLine)
		}
	}

}

func processKline(kline BinanceKline) {
	klen := len(Candlesticks)
	fmt.Println("Candlesticks Length: ", klen)

	if klen == 0 || Candlesticks[klen-1].IsKlineClosed {
		Candlesticks = append(Candlesticks, kline)
		return
	}

	Candlesticks[klen-1] = kline

	if kline.IsKlineClosed {
		color.Red("Candle Closed")
	}
}
