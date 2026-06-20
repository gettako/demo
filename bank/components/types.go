package components

import (
	"time"
)

type Transaction struct {
	ID        string    `json:"id"`
	Dest      string    `json:"dest"`
	Amount    int       `json:"amount"`
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
}
