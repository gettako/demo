// fraud_rpc.go
// Description: An RPC Service that acts as a Fraud Detection system.

package services

import (
	"context"
	"fmt"
	"time"

	"gettako.dev/tako/contracts"
)

// FraudProvider registers the Fraud Detection RPC handlers.
type FraudProvider struct{}

func NewFraudProvider() *FraudProvider {
	return &FraudProvider{}
}

func (p *FraudProvider) Register(app contracts.Application) error {
	return nil
}

func (p *FraudProvider) Boot(app contracts.Application) error {
	// Register the RPC handler
	app.RPC().Route("fraud.check").Handle(func(ctx context.Context, req contracts.RPCRequest) (contracts.RPCResponse, error) {
		amount, ok := req.Payload.(float64)
		if !ok {
			return contracts.RPCResponse{}, fmt.Errorf("invalid payload: expected float64")
		}

		// Simulate microservice network latency
		time.Sleep(800 * time.Millisecond)

		app.Logger().Info("Checking transaction fraud probability", "amount", amount)

		// Simple heuristic: amount > 50,000,000 is suspicious
		if amount > 50_000_000 {
			app.Logger().Warn("High-risk transaction blocked by Fraud Service", "amount", amount)
			return contracts.RPCResponse{Data: map[string]any{
				"safe":   false,
				"reason": "Amount exceeds normal daily limits. Marked as suspicious.",
			}}, nil
		}

		return contracts.RPCResponse{Data: map[string]any{
			"safe": true,
		}}, nil
	})

	return nil
}
