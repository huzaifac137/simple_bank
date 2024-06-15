package util

import (
	"math/rand"
	"strings"
)

var randomInt = rand.New(rand.NewSource(42))

const alphabet = "abcdefghijklmnopqrstuvwxyz"

func RandomString(n int) string {
	var sb strings.Builder
	k := len(alphabet)
	for i := 0; i < n; i++ {
		sb.WriteByte(alphabet[randomInt.Intn(k)])
	}
	return sb.String()
}

func RandomOwner() string {
	return RandomString(6)
}

func RandomMoney() int {
	return randomInt.Intn(1000)
}

func RandomCurrencies() string {
	currencies := []string{"USD", "EUR", "CAD"}
	k := len(currencies)
	return currencies[randomInt.Intn(k)]
}
