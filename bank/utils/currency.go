package utils

import (
	"fmt"
	"strings"
)

// FormatIDR formats a float64 into Indonesian Rupiah format (e.g. Rp 123.456,30).
func FormatIDR(amount float64) string {
	// Format to 2 decimal places with standard dot separator first
	str := fmt.Sprintf("%.2f", amount)
	parts := strings.Split(str, ".")

	integerPart := parts[0]
	decimalPart := "00"
	if len(parts) > 1 {
		decimalPart = parts[1]
	}

	// Add thousands separators (dot)
	var formattedInt strings.Builder
	runes := []rune(integerPart)
	for i, r := range runes {
		if i > 0 && (len(runes)-i)%3 == 0 && r != '-' {
			formattedInt.WriteRune('.')
		}
		formattedInt.WriteRune(r)
	}

	return fmt.Sprintf("Rp %s,%s", formattedInt.String(), decimalPart)
}
