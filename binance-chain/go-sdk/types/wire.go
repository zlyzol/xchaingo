package types

import (
	"github.com/tendermint/go-amino"
//	ctypes "github.com/tendermint/tendermint/rpc/core/types"
	ntypes "github.com/zlyzol/xchaingo/binance-chain/go-sdk/common/types"
	"github.com/zlyzol/xchaingo/binance-chain/go-sdk/types/tx"
)

func NewCodec() *amino.Codec {
	cdc := amino.NewCodec()
	//ctypes.RegisterAmino(cdc)
	ntypes.RegisterWire(cdc)
	tx.RegisterCodec(cdc)
	return cdc
}
