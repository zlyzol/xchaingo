package main

import (
	"github.com/zlyzol/xchaingo/common"
	"github.com/zlyzol/xchaingo/thorchain"
	"github.com/zlyzol/xchaingo/binance"
	"github.com/zlyzol/xchaingo/ethereum"

)

func main() {
	thorchain.NewThorchainClient(common.Testnet, "flip portion grant body mad mountain infant edit pig execute tired ridge")
	binance.NewBinanceClient(common.Testnet, "flip portion grant body mad mountain infant edit pig execute tired ridge")
	ethereum.NewEthereumClient(common.Testnet, "flip portion grant body mad mountain infant edit pig execute tired ridge")
}