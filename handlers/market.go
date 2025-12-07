package handlers

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"time"
)

type MarketData struct {
	Index   string  `json:"index"`
	Value   float64 `json:"value"`
	Change  float64 `json:"change"`
	Percent float64 `json:"percent"`
	Up      bool    `json:"up"`
}

func MarketDataHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Simulate real-time data
		// In a real app, this would fetch from an external API
		indices := []string{"SENSEX", "NIFTY", "NASDAQ", "DOW JONES"}
		rand.Seed(time.Now().UnixNano())
		idx := indices[rand.Intn(len(indices))]

		baseVal := 50000.0
		if idx == "NIFTY" {
			baseVal = 18000.0
		} else if idx == "NASDAQ" {
			baseVal = 14000.0
		} else if idx == "DOW JONES" {
			baseVal = 34000.0
		}

		change := (rand.Float64() * 200) - 100 // -100 to +100
		value := baseVal + change
		percent := (change / baseVal) * 100

		data := MarketData{
			Index:   idx,
			Value:   value,
			Change:  change,
			Percent: percent,
			Up:      change >= 0,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(data)
	}
}
