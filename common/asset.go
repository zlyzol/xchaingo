package common

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

var (
	BNBAsset, _     = NewChainAsset(BNBChain, "BNB")
	BTCAsset, _     = NewChainAsset(BTCChain, "BTC")
	RuneNative, _   = NewChainAsset(THORChain, "RUNE")
	RuneBNBAsset, _ = NewChainAsset(BNBChain, "RUNE-B1A")
	EmptyAsset   	= Asset{}
)

type Asset struct {
	Chain  Chain  `json:"chain" mapstructure:"chain"`
	Symbol Symbol `json:"symbol" mapstructure:"symbol"`
	Ticker Ticker `json:"ticker" mapstructure:"ticker"`
}

func NewAsset(input string) (Asset, error) {
	input = strings.ToUpper(input)
	var err error
	asset := Asset{}
	parts := strings.Split(input, ".")
	var sym string
	if len(parts) == 1 { // TODO I really dont think we should default at all.
		asset.Chain = ""
		sym = parts[0]
	} else {
		asset.Chain = Chain(parts[0])
		if err != nil {
			return Asset{}, err
		}
		sym = parts[1]
	}

	asset.Symbol, err = NewSymbol(sym)
	if err != nil {
		return Asset{}, err
	}

	parts = strings.Split(sym, "-")
	asset.Ticker, err = NewTicker(parts[0])
	if err != nil {
		return Asset{}, err
	}
	if asset.IsEmpty() { panic("bad asset") }
	return asset, nil
}

func NewChainAsset(chain string, input string) (Asset, error) {
	input = strings.ToUpper(input)
	var err error
	asset := Asset{}
	parts := strings.Split(input, ".")
	var sym string
	if len(parts) == 1 {
		sym = parts[0]
	} else {
		if !strings.EqualFold(chain, parts[0]) {
			return EmptyAsset, fmt.Errorf("chain conflict, chain parameter is %s, chain in asset name is %s", chain, parts[0])
		}
		sym = parts[1]
	}
	asset.Chain = Chain(chain)

	asset.Symbol, err = NewSymbol(sym)
	if err != nil {
		return Asset{}, err
	}

	parts = strings.Split(sym, "-")
	asset.Ticker, err = NewTicker(parts[0])
	if err != nil {
		return Asset{}, err
	}
	if asset.IsEmpty() { panic("bad asset") }
	return asset, nil
}

func (a Asset) Equal(a2 Asset) bool {
	return a.Chain.Equal(a2.Chain) && a.Ticker.Equal(a2.Ticker)
}
func (a Asset) IsEmpty() bool {
	return a.Chain.IsEmpty() || a.Symbol.IsEmpty() || a.Ticker.IsEmpty()
}

func (a Asset) String() string {
	return fmt.Sprintf("%s.%s", a.Chain.String(), a.Symbol.String())
}

func (a Asset) ChainTickerString() string {
	return fmt.Sprintf("%s.%s", a.Chain.String(), a.Ticker.String())
}

func RuneAsset() Asset {
	if strings.EqualFold(os.Getenv("NATIVE"), "false") {
		return RuneBNBAsset
	}
	return RuneNative
}

func IsBNBAsset(a Asset) bool {
	return a.Equal(BNBAsset)
}

func IsRuneAsset(a Asset) bool {
	return a.Equal(RuneBNBAsset) || a.Equal(RuneNative)
}

func (a Asset) MarshalJSON() ([]byte, error) {
	return json.Marshal(a.String())
}

func (a *Asset) UnmarshalJSON(data []byte) error {
	var err error
	var assetStr string
	if err := json.Unmarshal(data, &assetStr); err != nil {
		return err
	}
	*a, err = NewAsset(assetStr)
	return err
}
