package util

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

const alphabet = "abcdefghhijklmnopqrstuvwxyz"

// RandomInt generates a random integer between min and max
func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

// RandomString geerates a random string of length n
func RandomString(n int) string {
	var sb strings.Builder
	k := len(alphabet)

	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}

// RandomOwnerName generates random account owner names for testing
func RandomOwnerName() string {
	return RandomString(6)
}

// RandomAmount generates random amount to be used as balance for testing
func RandomAmount() int64 {
	return RandomInt(0, 1000)
}

// RandomCurrency returns one of the provided values from the provided slice of currencies
func RandomCurrency() string {
	currencies := []string{"EUR", "USD", "GBP", "INR"}
	n := len(currencies)
	return currencies[rand.Intn(n)]
}

// RandomEmail generates a random email address
func RandomEmail() string {
	return fmt.Sprintf("%s@email.com", RandomString(6))
}
