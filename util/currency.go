package util

// Constants for all supported currencies
const (
	USD = "USD"
	EUR = "EUR"
	GBP = "GBP"
	INR = "INR"
)

// IsSupportedCurrency checks is a currency is supported or not
func IsSupportedCurrency(currency string) bool {
	switch currency {
	case USD, EUR, GBP, INR:
		return true
	}
	return false
}
