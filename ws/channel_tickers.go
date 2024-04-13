package ws

import (
	"errors"
)

var _ Channel = TickersChannel{}

// TickersChannel
//
// https://www.okx.com/docs-v5/en/#order-book-trading-market-data-ws-tickers-channel
type TickersChannel struct {
	// Channel name
	Channel string `json:"channel,omitempty"`
	// Instrument ID
	InstId string `json:"instId"`
}

// GetChannel implement Channel interface
func (t TickersChannel) GetChannel() (string, error) {
	return t.Channel, nil
}

func (c TickersChannel) validate() error {
	if c.Channel == "" {
		return errors.New("channel cann't be empty")
	}
	if c.InstId == "" {
		return errors.New("instId cann't be emtpy")
	}
	return nil
}

// TickersData
type TickersData struct {
	// Instrument type
	InstType string `json:"instType"`
	// Instrument ID
	InstID string `json:"instId"`
	// Last traded price
	Last string `json:"last"`
	// Last traded size
	LastSz string `json:"lastSz"`
	// Best ask price
	AskPx string `json:"askPx"`
	// Best ask size
	AskSz string `json:"askSz"`
	// Best bid price
	BidPx string `json:"bidPx"`
	// Best bid size
	BidSz string `json:"bidSz"`
	// Open price in the past 24 hours
	Open24H string `json:"open24h"`
	// Highest price in the past 24 hours
	High24H string `json:"high24h"`
	// Lowest price in the past 24 hours
	Low24H string `json:"low24h"`
	// Open price in the UTC 0
	SodUtc0 string `json:"sodUtc0"`
	// Open price in the UTC 8
	SodUtc8 string `json:"sodUtc8"`
	// 24h trading volume, with a unit of currency.
	// If it is a derivatives contract, the value is the number of base currency.
	// If it is SPOT/MARGIN, the value is the quantity in quote currency.
	VolCcy24H string `json:"volCcy24h"`
	// 24h trading volume, with a unit of contract.
	// If it is a derivatives contract, the value is the number of contracts.
	// If it is SPOT/MARGIN, the value is the quantity in base currency.
	Vol24H string `json:"vol24h"`
	// Ticker data generation time, Unix timestamp format in milliseconds, e.g. 1597026383085
	Ts string `json:"ts"`
}
