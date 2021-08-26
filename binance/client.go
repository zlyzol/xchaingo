package binance

import (
	"fmt"
	"strings"

	"github.com/zlyzol/xchaingo/common"

	bin "github.com/zlyzol/xchaingo/binance-chain/go-sdk/client"
	"github.com/zlyzol/xchaingo/binance-chain/go-sdk/common/types"
	"github.com/zlyzol/xchaingo/binance-chain/go-sdk/keys"

	"github.com/zlyzol/xchaingo/binance-chain/go-sdk/client/transaction"
	"github.com/zlyzol/xchaingo/binance-chain/go-sdk/types/msg"
)
var key		keys.KeyManager

type SingleAndMultiFees struct {
	Single	common.Fees
	Multi	common.Fees
}
type MultiTransfer struct {
	To		common.Address
	Coins	[]common.Balance
}
type MultiSendParams struct {
	WalletIndex		uint64
	Transactions	[]MultiTransfer
	memoM		string
}

var (
	decimal = 8
	defaultGas = 2000000
	maxTxCount = 100

	// urls
	defaultTestnetExplorerUrl	= "https://testnet-explorer.binance.org"
	defaultMainnetExplorerUrl	= "https://explorer.binance.org"
	defaultTestnetNode	= "testnet-dex.binance.org"
	defaultMainnetNode	= "dex.binance.org"

	// thornode node url endpoints
	balanceEndpoint = "/bank/balances/"
	accountEndpoint = "/auth/accounts/"
	// explorer url endpoints
	txEndpoint 		= "/tx/"
	addressEndpoint = "/address/"

	addrPrefix		= map[network]string {testnet: "tbnb", mainnet: "bnb"}

	chain			= common.BNBChain
)
type network = common.Network
var testnet = common.Testnet
var mainnet = common.Mainnet
var AssetBNB = common.BNBAsset
var defaultFees = common.SingleFee(common.Fee(common.NewUint(uint64(defaultGas))))

type BinanceClient interface {
//	GetBncClient() BncClient
	GetXChainClient() common.XChainClient
	GetMultiSendFees() common.Fees
	GetSingleAndMultiFees() SingleAndMultiFees
	MultiSend(params MultiSendParams) common.TxHash
}

type client struct {
	network		network
	node		string
	dex			bin.DexClient
	key			keys.KeyManager
	wallet		common.Address
	walletIndex	uint32
}
func NewBinanceClient(network network, phrase string) (BinanceClient, error) {
	c := client{}
	c.SetNetwork(network)
	_, err := c.SetPhrase(phrase)
	if err != nil {
		return nil, err
	}
	c.dex, err = bin.NewDexClient(c.node, getBinNet(c.network), c.key)
	return &c, nil
}
func getBinNet(network network) types.ChainNetwork {
	if network == testnet {
		return types.TestNetwork
	}
	return types.ProdNetwork
}
func (c *client) GetXChainClient() common.XChainClient {
	return c
}
func (c *client) SetNetwork(network common.Network) {
	c.network = network
	if c.network == testnet {
		c.node = defaultTestnetNode
	} else {
		c.node = defaultMainnetNode
	}
	types.Network = getBinNet(c.network)
}
func (c *client) GetNetwork() common.Network {
	return c.network
}
func (c *client) GetExplorerUrl() string {
	if c.network == testnet {
		return defaultTestnetExplorerUrl
	}
	return defaultMainnetExplorerUrl
}
func (c *client) GetExplorerAddressUrl(address string) string {
	return c.GetExplorerUrl() + addressEndpoint + address
}
func (c *client) GetExplorerTxUrl(txID string) string {
	return c.GetExplorerUrl() + txEndpoint + txID
}
func (c *client) ValidateAddress(address common.Address) bool {
	if !strings.HasPrefix(string(address), addrPrefix[c.network]) {
		return false
	}
	return common.ValidateAddress(address)
}
func (c *client) GetAddress() (common.Address, error) {
	return c.wallet, nil
}
func (c *client) SetPhrase(phrase string) (common.Address, error) {
	//hdPath := fmt.Sprintf("0'/0/%d", c.walletIndex)
	var err error
//	c.key, err = keys.NewMnemonicPathKeyManager(phrase, hdPath)
	c.key, err = keys.NewKeyStoreKeyManager("/Users/zolo/Documents/CryptoWallety/BNB/test-chain/xchain-go/w1/keystore.txt", "*Ivranka95123*")
	if err != nil {
		return "", err
	}
	bech32 := c.key.GetAddr()
	c.wallet = common.Address(bech32.String())
	return c.wallet, err
}
func (c *client) GetBalance(address common.Address, assets []common.Asset) (common.Balances, error) {
	balance, err := c.dex.GetAccount(string(c.wallet))
	if err != nil {
		return nil, err
	}
	xbal := make(common.Balances, 0, len(balance.Balances))
	for _, bb := range balance.Balances {
		xb := common.Balance{
			Asset:	newAsset(bb.Symbol),
			Amount:	common.BaseAmount(bb.Free.ToInt64()),
		}
		xbal = append(xbal, xb)
	}
	return xbal, nil
}
func (c *client) GetTransactions(params common.TxHistoryParams) (common.TxsPage, error) {
	panic("not implemented")
}
func (c *client) GetTransactionData(txId string, assetAddress common.Address) (common.Tx, error) {
	panic("not implemented")
}
func (c *client) GetFees(feeParams interface{}) common.Fees {
	return defaultFees
}
func (c *client) Transfer(params common.TxParams) (common.TxHash, error) {
	if !c.initialized() {
		return "", fmt.Errorf("BinanceClient struct not initialized")
	}
	if params.WalletIndex != 0 {
		panic("Not supported Walletindex other than 0")
	}
	toAddr, err := types.AccAddressFromBech32(string(params.Recipient))
	if err != nil {
		return "", err
	}
	msgs := []msg.Transfer{{
		ToAddr: toAddr,
		Coins:  types.Coins{types.Coin{Denom: params.Asset.Symbol.String(), Amount: int64(params.Amount.Fx8Uint64())}},
	}}
	res, err := c.dex.SendToken(msgs, true, transaction.WithMemo(params.Memo))
	//resSend, err := acc.bdex.dex.SendToken(msgs, true, transaction.WithMemo(depositAddress.Memo), transaction.WithAcNumAndSequence(acc.number, acc.sequence))
	if err != nil {
		return "", err
	}
	if !res.Ok {
		return common.TxHash(res.Hash), fmt.Errorf("SendToken result Ok = false, log: %s", res.Log)
	}
	return common.TxHash(res.Hash), nil
}
func (c *client) PurgeClient() {
	c.node = ""
	c.dex = nil
	c.key = nil
	c.wallet = ""
	c.walletIndex = 0
}

  // BinanceClient interface
func (c *client) GetMultiSendFees() common.Fees {
  	panic("not implemented")
}
func (c *client) GetSingleAndMultiFees() SingleAndMultiFees {
	panic("not implemented")
}
func (c *client) MultiSend(params MultiSendParams) common.TxHash {
	panic("not implemented")
}

// helper functions
func newAsset(s string) common.Asset {
	asset, err := common.NewChainAsset(chain, s)
	if err != nil {
		panic(err)
	}
	return asset
}
func (c *client) initialized() bool {
	return c.node != "" && c.dex != nil && c.key != nil && c.wallet != ""
}
