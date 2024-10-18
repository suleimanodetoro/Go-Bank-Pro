package util

// Constants for supported currencies
// These represent some of the world's most widely used and traded currencies
const (
	USD = "USD" // United States Dollar
	EUR = "EUR" // Euro
	JPY = "JPY" // Japanese Yen
	GBP = "GBP" // British Pound Sterling
	AUD = "AUD" // Australian Dollar
	CAD = "CAD" // Canadian Dollar
	CHF = "CHF" // Swiss Franc
	CNY = "CNY" // Chinese Yuan
	HKD = "HKD" // Hong Kong Dollar
	NZD = "NZD" // New Zealand Dollar
	SEK = "SEK" // Swedish Krona
	KRW = "KRW" // South Korean Won
	SGD = "SGD" // Singapore Dollar
	NOK = "NOK" // Norwegian Krone
	MXN = "MXN" // Mexican Peso
	INR = "INR" // Indian Rupee
	RUB = "RUB" // Russian Ruble
	ZAR = "ZAR" // South African Rand
	TRY = "TRY" // Turkish Lira
	BRL = "BRL" // Brazilian Real
)

// IsSupportedCurrency returns true if the currency is supported
func IsSupportedCurrency(currency string) bool {
	switch currency {
	case USD, EUR, JPY, GBP, AUD, CAD, CHF, CNY, HKD, NZD, SEK, KRW, SGD, NOK, MXN, INR, RUB, ZAR, TRY, BRL:
		return true
	}
	return false
}
