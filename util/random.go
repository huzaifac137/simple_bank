package util

import (
	"fmt"
	"math/rand"
	"strings"
)

var randomInt = rand.New(rand.NewSource(1000))

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
	return RandomString(15)
}

func RandomMoney() int {
	return randomInt.Intn(1000)
}

func RandomCurrencies() string {
	currencies := []string{"USD", "EUR", "CAD"}
	k := len(currencies)
	return currencies[randomInt.Intn(k)]
}

func RandomEmail() string {
	return fmt.Sprintf("%s@email.com", RandomString(15))
}
