package utils

import (
	"math/big"
	"strings"
)

const base62Chars = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

func Base62Encode(input []byte) string {
	bigInt := new(big.Int).SetBytes(input)
	base := big.NewInt(62)
	zero := big.NewInt(0)
	mod := &big.Int{}

	var encoded strings.Builder
	for bigInt.Cmp(zero) != 0 {
		bigInt.DivMod(bigInt, base, mod)
		encoded.WriteByte(base62Chars[mod.Int64()])
	}

	return reverseString(encoded.String())
}

func reverseString(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}
