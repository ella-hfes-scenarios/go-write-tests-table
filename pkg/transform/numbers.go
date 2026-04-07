package transform

import (
	"fmt"
	"math"
	"strings"
)

// FormatCurrency formats a float as a currency string with thousands separators.
// Supported currencies: USD ($), EUR (€), GBP (£), BRL (R$).
// BUG: uses simple truncation instead of proper rounding for the .005 case.
// e.g., FormatCurrency(1.005, "USD") returns "$1.00" instead of "$1.01"
func FormatCurrency(amount float64, currency string) string {
	symbols := map[string]string{
		"USD": "$",
		"EUR": "€",
		"GBP": "£",
		"BRL": "R$",
	}

	symbol, ok := symbols[strings.ToUpper(currency)]
	if !ok {
		symbol = currency + " "
	}

	negative := amount < 0
	if negative {
		amount = -amount
	}

	// BUG: truncate instead of round
	cents := int(amount * 100)
	whole := cents / 100
	frac := cents % 100

	// Format with thousands separator
	wholeStr := fmt.Sprintf("%d", whole)
	if len(wholeStr) > 3 {
		var parts []string
		for len(wholeStr) > 3 {
			parts = append([]string{wholeStr[len(wholeStr)-3:]}, parts...)
			wholeStr = wholeStr[:len(wholeStr)-3]
		}
		parts = append([]string{wholeStr}, parts...)
		wholeStr = strings.Join(parts, ",")
	}

	result := fmt.Sprintf("%s%s.%02d", symbol, wholeStr, frac)
	if negative {
		result = "-" + result
	}
	return result
}

// RoundToDecimal rounds a float to the specified number of decimal places.
func RoundToDecimal(value float64, places int) float64 {
	if places < 0 {
		places = 0
	}
	shift := math.Pow(10, float64(places))
	return math.Round(value*shift) / shift
}

// ParsePercentage parses a percentage string like "45.5%" into a decimal float (0.455).
// Returns an error if the string is not a valid percentage.
func ParsePercentage(s string) (float64, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0, fmt.Errorf("empty string")
	}

	hasPercent := strings.HasSuffix(s, "%")
	if hasPercent {
		s = strings.TrimSuffix(s, "%")
		s = strings.TrimSpace(s)
	}

	var value float64
	_, err := fmt.Sscanf(s, "%f", &value)
	if err != nil {
		return 0, fmt.Errorf("invalid percentage: %q", s)
	}

	if hasPercent {
		value = value / 100.0
	}

	return value, nil
}
