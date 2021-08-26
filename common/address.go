package common

import (
	"strings"
	"github.com/btcsuite/btcutil/bech32"
)

func ValidateAddress(address Address) bool {
	// Check bech32 addresses, would succeed any string bech32 encoded
	decodedAddress, decodedBytes, err := bech32.Decode(string(address))
	if err != nil {
		return false
	}
	_, _ = decodedAddress, decodedBytes
	decoded, err := bech32.Encode(decodedAddress, decodedBytes)
	if err != nil {
		return false
	}
	if !strings.EqualFold(decoded, string(address)) {
		return false
	}
	return true
}
