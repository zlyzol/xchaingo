package common

import (
  "strings"
)

type Chain string

const (
  BNBChain string = "BNB"
  BTCChain string = "BTC"
  ETHChain string = "ETH"
  THORChain string = "THOR"
  CosmosChain string = "GAIA"
  PolkadotChain string = "POLKA"
  BCHChain string = "BCH"
  LTCChain string = "LTC"
)

var Chains = []string{BNBChain, BTCChain, ETHChain, THORChain, CosmosChain, PolkadotChain, BCHChain, LTCChain}


/**
 * Type guard to check whether string  is based on type `Chain`
 *
 * @param {string} chain The chain string.
 * @returns {boolean} `true` or `false`
 */
func IsChain(chainId string) bool {
  for _, c := range Chains {
    if strings.EqualFold(chainId, c) {
      return true
    }
  }
  return false
}

/**
 * Convert chain to string.
 *
 * @param {Chain} chainId.
 * @returns {string} The string based on the given chain type.
 */
func ChainToString(chainId string) string {
  if !IsChain(chainId) {
    return "unknown chain"
  }
  return map[string]string{
    BNBChain: "Binance Chain",
    BTCChain: "Bitcoin",
    ETHChain: "Ethereum",
    THORChain: "Thorchain",
    CosmosChain: "Cosmos",
    PolkadotChain: "Polkadot",
    BCHChain: "Bitcoin Cash",
    LTCChain: "Litecoin",
  }[chainId]
}
func (c Chain) Equal(c2 Chain) bool {
	return strings.EqualFold(string(c), string(c))
}
func (c Chain) IsEmpty() bool {
	return strings.TrimSpace(string(c)) == ""
}

func (c Chain) String() string {
	return string(c)
}