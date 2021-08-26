package common

import (
	"time"
)

type Date time.Time
type BaseAmount = Uint
func AssetToString(asset Asset) string {return Asset(asset).String()}
func AmountFromString(s string) BaseAmount{return BaseAmount(NewUintFromString(s))}

type Address string
type Network uint8
const (
  Mainnet Network = iota
  Testnet
)
type Balance struct {
  Asset		Asset
  Amount	BaseAmount
}
type Balances []Balance
type TxType uint8
const (
  Transfer TxType = iota
  Unknown
)
type TxHash string
type TxTo struct {
  To		  Address // address
  Amount	BaseAmount // amount
}
type TxFrom struct {
  From		Address	//		Address | TxHash // address or tx id
  Amount	BaseAmount // amount
}
type Tx struct {
  Asset	Asset // asset
  From	[]TxFrom // list of "from" txs. BNC will have one `TxFrom` only, `BTC` might have many transactions going "in" (based on UTXO)
  To	  []TxTo // list of "to" transactions. BNC will have one `TxTo` only, `BTC` might have many transactions going "out" (based on UTXO)
  Date	Date // timestamp of tx
  Type	TxType // type
  Hash	string // Tx hash
}
type TxsPage struct {
  Total	uint64
  Txs	[]Tx
}
type TxHistoryParams struct {
  Address	Address // Address to get history for
  Offset	uint64 // Optional Offset
  Limit		uint64 // Optional Limit of transactions
  StartTime	Date // Optional start time
  Asset		string // Optional asset. Result transactions will be filtered by this asset
}
type TxParams struct {
  WalletIndex uint64 // send from this HD index
  Asset       Asset
  Amount      BaseAmount
  Recipient   Address
  Memo        string // optional memo to pass
}
type FeeRate uint64
type FeeRates struct {
	Average	FeeRate
	Fast	FeeRate
	Fastest	FeeRate
}
type FeeType uint8
const (
  FlatFee FeeType = iota
  PerByte
)
type Fee BaseAmount
type FeeOptions struct {
	Average	Fee
	Fast	Fee
	Fastest	Fee
}
type Fees struct {
	FeeOptions  FeeOptions
  Type        FeeType
}
type FeesWithRates struct {
	Rates	FeeRates
	Fees	Fees
}
// FeesParams:
// In most cases, clients don't expect any paramter in `getFees`
// but in some cases, they do (e.g. in xchain-ethereum).
// To workaround this, we just define an interface param for now.
// If needed, any client can extend `FeeParams` to add more  (Check `xchain-ethereum` as an example)
type FeesParams interface{}

type XChainClient interface {
	SetNetwork(net Network)
	GetNetwork() Network

	GetExplorerUrl() string
	GetExplorerAddressUrl(address string) string
	GetExplorerTxUrl(txID string) string

	ValidateAddress(address Address) bool
	GetAddress() (Address, error)

	SetPhrase(phrase string) (Address, error)
	GetBalance(address Address, assets []Asset) (Balances, error)
	GetTransactions(params TxHistoryParams) (TxsPage, error)
	GetTransactionData(txId string, assetAddress Address) (Tx, error)
	GetFees(feeParams interface{}) Fees // nil or FeeParams
	Transfer(params TxParams) (TxHash, error)
	PurgeClient()
}
