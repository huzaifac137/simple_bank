package util

func IsValidCurrency(currency string) bool {
	switch currency {
	case "USD", "EUR", "CAD":
		return true
	}
	return false
}
