// exchange_job.go
// Description: A background job that fetches simulated exchange rates.

package services

import (
	"encoding/json"
	"math/rand"
	"time"

	"gettako.dev/tako/contracts"
)

// ExchangeProvider registers a background job to fetch exchange rates.
type ExchangeProvider struct{}

func NewExchangeProvider() *ExchangeProvider {
	return &ExchangeProvider{}
}

func (p *ExchangeProvider) Register(app contracts.Application) error {
	return nil
}

func (p *ExchangeProvider) Boot(app contracts.Application) error {
	// Schedule a recurring job to fetch exchange rates every 5 seconds
	_ = app.Context().Jobs().Every(5*time.Second, func() {
		// Simulate HTTP request latency
		time.Sleep(500 * time.Millisecond)

		// Simulate fetching from api.exchangerate.host or similar
		baseRate := 15000.0                         // Base IDR to USD
		fluctuation := (rand.Float64() * 200) - 100 // -100 to +100 IDR fluctuation

		rates := map[string]float64{
			"USD_IDR": baseRate + fluctuation,
			"EUR_IDR": (baseRate + fluctuation) * 1.08,
			"SGD_IDR": (baseRate + fluctuation) * 0.74,
		}

		ratesJSON, _ := json.Marshal(rates)

		// Store in Cache
		if cache := app.Cache(); cache != nil {
			_ = cache.Set("exchange_rates", string(ratesJSON), 2*time.Minute)
			app.Events().Dispatch("exchange_rates_updated", string(ratesJSON))
		}

		app.Logger().Info("Exchange rates updated successfully.")
	})

	// Trigger immediately on boot so UI has initial data
	go func() {
		if cache := app.Cache(); cache != nil {
			rates := map[string]float64{
				"USD_IDR": 15000.0,
				"EUR_IDR": 16200.0,
				"SGD_IDR": 11100.0,
			}
			ratesJSON, _ := json.Marshal(rates)
			_ = cache.Set("exchange_rates", string(ratesJSON), 2*time.Minute)
		}
	}()

	return nil
}
