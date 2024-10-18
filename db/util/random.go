package util

import (
	"math/rand"
	"strings"
	"time"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"

// create a package-level random generator that can be reused
var rng = rand.New(rand.NewSource(time.Now().UnixNano()))

// RandomInt generates a random int between min and max
func RandomInt(min, max int64) int64 {
	return min + rng.Int63n(max-min+1)
}

// Generate random string of n length
func RandomString(n int) string {
	var sb strings.Builder
	k := len(alphabet)

	for i := 0; i < n; i++ {
		c := alphabet[rng.Intn(k)]
		sb.WriteByte(c)
	}
	return sb.String()
}

// Generate random user name
func RandomOwner() string {
	return RandomString(6)
}

// Create random amount of money
func RandomMoney() int64 {
	return RandomInt(0, 100)
}

// RandomCurrency returns a random supported currency from the list of constants
func RandomCurrency() string {
	// Create a slice of the supported currencies using the constants
	currencies := []string{
		USD, EUR, JPY, GBP, AUD, CAD, CHF, CNY, HKD, NZD,
		SEK, KRW, SGD, NOK, MXN, INR, RUB, ZAR, TRY, BRL,
	}
	// Generate a random index to select a currency
	n := len(currencies)
	return currencies[rng.Intn(n)]
}
