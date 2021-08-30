module github.com/zlyzol/xchaingo

go 1.15

require (
	github.com/binance-chain/ledger-cosmos-go v0.9.9
	github.com/btcsuite/btcd v0.21.0-beta
	github.com/btcsuite/btcutil v1.0.2
	github.com/cosmos/cosmos-sdk v0.42.1
	github.com/cosmos/go-bip39 v1.0.0
	github.com/gorilla/websocket v1.4.2
	github.com/pkg/errors v0.9.1
	github.com/stretchr/testify v1.7.0
	github.com/tendermint/btcd v0.1.1
	github.com/tendermint/go-amino v0.16.0
	github.com/tendermint/tendermint v0.34.8
	gitlab.com/thorchain/thornode v0.64.0
	golang.org/x/crypto v0.0.0-20210711020723-a769d52b0f97
	gopkg.in/check.v1 v1.0.0-20200902074654-038fdea0a05b
	gopkg.in/resty.v1 v1.12.0

)

replace (
	github.com/agl/ed25519 => github.com/binance-chain/edwards25519 v0.0.0-20200305024217-f36fc4b53d43
	github.com/binance-chain/tss-lib => gitlab.com/thorchain/tss/tss-lib v0.0.0-20201118045712-70b2cb4bf916
	github.com/gogo/protobuf => github.com/regen-network/protobuf v1.3.2-alpha.regen.4
	github.com/tendermint/go-amino => github.com/binance-chain/bnc-go-amino v0.14.1-binance.1
	github.com/zondax/ledger-go => github.com/binance-chain/ledger-go v0.9.1
)
