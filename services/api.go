package services

import (
	"crypto/rand"
	"fmt"
)

// APIClient demonstrates a Transient service in the IoC container.
// A new instance is created every time it is resolved via app.Make().
type APIClient struct {
	ID string
}

// NewAPIClient creates a new mock API client with a unique identifier.
func NewAPIClient() *APIClient {
	b := make([]byte, 4)
	rand.Read(b)
	return &APIClient{
		ID: fmt.Sprintf("CLIENT-%x", b),
	}
}

// FetchData simulates an API call.
func (a *APIClient) FetchData() string {
	return "Data from " + a.ID
}
