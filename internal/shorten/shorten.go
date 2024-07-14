package shorten

import (
	"crypto/sha256"
	"fmt"
	"math/big"

	"github.com/itchyny/base58-go"
)

// sha256Of функция хеширует входящую строку
func sha256Of(input string) []byte {
	algorithm := sha256.New()
	algorithm.Write([]byte(input))
	return algorithm.Sum(nil)
}

// base58Encoded функция шифрует массив байт
func base58Encoded(bytes []byte) (string, error) {
	encoding := base58.BitcoinEncoding
	encoded, _ := encoding.Encode(bytes)
	return string(encoded), nil
}

// Shorten функция преобразует оригинальный URL в сокращенный
func Shorten(initialURL string) string {
	urlHashBytes := sha256Of(initialURL)
	generatedNumber := new(big.Int).SetBytes(urlHashBytes).Uint64()

	finalString, err := base58Encoded([]byte(fmt.Sprintf("%d", generatedNumber)))
	if err != nil {
		return ""
	}
	return finalString[:8]
}
