package thorchain

import (
	"time"
	"sync"

	"github.com/zlyzol/xchaingo/common"

	hd "github.com/cosmos/cosmos-sdk/crypto/hd"
	cKeys "github.com/cosmos/cosmos-sdk/crypto/keyring"
	"gitlab.com/thorchain/thornode/bifrost/thorclient"
	thorcommon "gitlab.com/thorchain/thornode/common"
	"gitlab.com/thorchain/thornode/common/cosmos"
	"gitlab.com/thorchain/thornode/x/thorchain"
	"gitlab.com/thorchain/thornode/bifrost/config"
	tmetrics "gitlab.com/thorchain/thornode/bifrost/metrics"
	"gitlab.com/thorchain/thornode/x/thorchain/types"
)

var (
	once sync.Once
	metrics *tmetrics.Metrics
)

func (c *client) updateBridge() error {
	cfg, _, kb, err := c.prepareBridge()
	if err != nil {
		return err
	}
	c.bridge, err = thorclient.NewThorchainBridge(cfg, metrics, thorclient.NewKeysWithKeybase(kb, cfg.SignerName, cfg.SignerPasswd))
	return err
}
// NewFundraiserParams creates a BIP 44 parameter object from the params:
// m / 44' / coinType' / account' / 0 / address_index
// The fixed parameters (purpose', coin_type', and change) are determined by what was used in the fundraiser.
//func NewFundraiserParams(account, coinType, addressIdx uint32) *BIP44Params {
//	return NewParams(44, coinType, account, false, addressIdx)
//}

func (c *client) prepareBridge() (config.ClientConfiguration, cKeys.Info, cKeys.Keyring, error) {
	cfg := config.ClientConfiguration{
		ChainID:         thorcommon.Chain(chainId),
		ChainHost:       c.getNode(),
		ChainRPC:        c.getRpc(),
		SignerName:      "",//"thorchain",
		SignerPasswd:    "",//"*Ivranka95123*",
		ChainHomeFolder: "",
	}
	thorchain.SetupConfigForTest()
	kb := cKeys.NewInMemory()
	params := *hd.NewFundraiserParams(0, coinType, c.walletIndex)
	hdPath := params.String()
	c.updateConfig(hdPath)
	info, err := kb.NewAccount(cfg.SignerName, c.phrase, cfg.SignerPasswd, hdPath, hd.Secp256k1)
	return cfg, info, kb, err
}
func (c *client) updateConfig(hdPath string) {
	p := prefixes[c.network]
	config := cosmos.GetConfig()
	config.SetBech32PrefixForAccount(p.AccAddr, p.AccPub)
	config.SetBech32PrefixForValidator(p.ValAddr, p.ValPub)
	config.SetBech32PrefixForConsensusNode(p.ConsAddr, p.ConsPub)
	config.SetCoinType(coinType)
	config.SetFullFundraiserPath(hdPath)
	//config.Seal()
}

func getMetrics() *tmetrics.Metrics {
	once.Do(func() { // <-- atomic, does not allow repeating
		var err error
		metrics, err = tmetrics.NewMetrics(config.MetricsConfiguration{
			Enabled:      false,
			ListenPort:   9000,
			ReadTimeout:  time.Second,
			WriteTimeout: time.Second,
			Chains:       thorcommon.Chains{thorcommon.Chain(common.THORChain)},
		})
		if err != nil {
			panic(err)
		}
	})
	return metrics
}

func bridgeTransfer(bridge *thorclient.ThorchainBridge, params common.TxParams) (common.TxHash, error) {
	if params.WalletIndex != 0 {
		panic("Not supported Walletindex other than 0")
	}
	cosmosCoins := cosmos.Coins{
		cosmos.NewCoin(params.Asset.String(), cosmos.NewInt(params.Amount.Fx8Int64())),
	}
	addr, err := cosmos.AccAddressFromBech32(string(params.Recipient))
	msg := types.NewMsgSend(bridge.GetContext().GetFromAddress(), addr, cosmosCoins)
	txID, err := bridge.Broadcast(msg)
	return common.TxHash(txID), err
}

func bridgeDeposit(bridge *thorclient.ThorchainBridge, params DepositParams) (common.TxHash, error) {
	if params.WalletIndex != 0 {
		panic("Not supported Walletindex other than 0")
	}
	commonCoins := thorcommon.Coins{
		thorcommon.NewCoin(thorcommon.RuneAsset(), cosmos.NewUint(uint64(params.Amount.Fx8Int64()))),
	};
	msg := types.NewMsgDeposit(commonCoins, params.Memo, bridge.GetContext().GetFromAddress())
	//"SWAP:BNB.BNB:tbnb12c9lnfdryjdc9fg88a4lj5kd98agjcttk7jpa7"
	txID, err := bridge.Broadcast(msg)
	return common.TxHash(txID), err
}

/*
func getAddressFromMnemonic(mnemonic string) string {
	kb := cKeys.NewInMemory()
	params := *hd.NewFundraiserParams(0, THORChainCoinType, 0)
	hdPath := params.String()
	info, err := kb.NewAccount(cfg.SignerName, mnemonic, cfg.SignerPasswd, hdPath, hd.Secp256k1)
	//tthor1fnzhjcnqn33fahdywf7t5azcapjry83r3h5j3g
	if err != nil {
		panic(err)
	}
	s := info.GetAddress().String()
	return s
}
*/