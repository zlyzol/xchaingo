package thorchain

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/zlyzol/xchaingo/common"
	"gitlab.com/thorchain/thornode/bifrost/thorclient"
	//	"gitlab.com/thorchain/thornode/x/thorchain/types"
	//	hd "github.com/cosmos/cosmos-sdk/crypto/hd"
	//	cKeys "github.com/cosmos/cosmos-sdk/crypto/keyring"
)

const (
	decimal = 8
	defaultGas = 2000000
	maxTxCount = 100

	chainId		= "thorchain"
	coinType	= 931
	chainStr	= common.THORChain

	// thornode port: 1317
	// midgard port: 8080
	// tendermint testnet port :26657
	// tendermint mainnet port: 27147

	// urls
	protocol			= "http://"
	defaultExplorerUrl	= "https://viewblock.io/thorchain"
	defaultTestnetNode	= "18.159.71.230:1317" // "testnet.thornode.thorchain.info"
	defaultTestnetRpc	= "18.159.71.230:26657" // "testnet.rpc.thorchain.info"
	defaultMannetNode	= "18.214.28.114:1317" // "thornode.thorchain.info"
	defaultMannetRpc	= "18.214.28.114:26657" // "rpc.thorchain.info"
	midgardPort			= ":8080"
	tendermintTestnetPort		= ":26657"
	tendermintMainnetPort		= ":27147"

	// thornode node url endpoints
	balanceEndpoint = "/bank/balances/"
	accountEndpoint = "/auth/accounts/"
	// explorer url endpoints
	txEndpoint 		= "/tx/"
	addressEndpoint = "/address/"
	explorerTestnetPostfix = "?network=testnet"
)
type network = common.Network
var testnet = common.Testnet
var mainnet = common.Mainnet
var tenderminPort = map[network]string{testnet: tendermintTestnetPort, mainnet: tendermintMainnetPort}
var defaultFees = common.SingleFee(defaultGas)
var AssetRune = common.RuneNative
var prefixes = map[network]Prefixes{
	mainnet: {
		AccAddr:  "thor",
		AccPub:   "thorpub",
		ValAddr:  "thorv",
		ValPub:   "thorvpub",
		ConsAddr: "thorc",
		ConsPub:  "thorcpub",
	},
	testnet: {
		AccAddr:  "tthor",
		AccPub:   "tthorpub",
		ValAddr:  "tthorv",
		ValPub:   "tthorvpub",
		ConsAddr: "tthorc",
		ConsPub:  "tthorcpub",
	},
}
var defaultClientUrl = ClientUrl{
	testnet: NodeUrl{
		Node: defaultTestnetNode,
		Rpc:  defaultTestnetRpc,
	},
	mainnet: NodeUrl{
		Node: defaultMannetNode,
		Rpc:  defaultMannetRpc,
	},
}
var defaultExplorerUrls = ExplorerUrls{
	Root: ExplorerUrl{
		testnet: defaultExplorerUrl,
		mainnet: defaultExplorerUrl,
	},
	Tx: ExplorerUrl{
		testnet: defaultExplorerUrl + txEndpoint,
		mainnet: defaultExplorerUrl + txEndpoint,
	},
	Address: ExplorerUrl{
		testnet: defaultExplorerUrl + addressEndpoint,
		mainnet: defaultExplorerUrl + addressEndpoint,
	},
}

type NodeUrl struct {
	Node string
	Rpc  string
}
type ClientUrl map[network]NodeUrl
type ExplorerUrl map[network]string
type ExplorerUrls struct {
	Root    ExplorerUrl
	Tx      ExplorerUrl
	Address ExplorerUrl
}
type ThorchainClientParams struct {
	ClientUrl    ClientUrl
	ExplorerUrls ExplorerUrls
}
type DepositParams struct {
	WalletIndex uint64
	Asset       common.Asset
	Amount      common.BaseAmount
	Memo        string
}
type Prefixes struct {
	AccAddr  string
	AccPub   string
	ValAddr  string
	ValPub   string
	ConsAddr string
	ConsPub  string
}

type ThorchainClient interface {
	Deposit(params DepositParams) (common.TxHash, error)
	GetXChainClient() common.XChainClient
}

type client struct {
	network				network
	clientUrl			ClientUrl
	explorerUrls		ExplorerUrls
	phrase				string
	bridge				*thorclient.ThorchainBridge
	walletIndex			uint32
}

func NewThorchainClient(network network, phrase string) (ThorchainClient, error) {
	getMetrics()
	c := client{
		network:		network,
		clientUrl:		defaultClientUrl,
		explorerUrls:	defaultExplorerUrls,
	}
	c.SetNetwork(network)
	_, err := c.SetPhrase(phrase)
	if err != nil {
		return nil, err
	}
	return &c, nil
}
func (c *client) GetXChainClient() common.XChainClient {
	return c
}
func (c *client) SetNetwork(network common.Network) {
    c.network = network
}
func (c *client) GetNetwork() common.Network {
	return c.network
}
func (c *client) GetExplorerUrl() string {
	return c.explorerUrls.Root[c.network]	
}
func (c *client) GetExplorerAddressUrl(address string) string {
	url := c.explorerUrls.Address[c.network] + address
	if c.network == testnet {
		url += explorerTestnetPostfix
	}
	return url
}
func (c *client) GetExplorerTxUrl(txID string) string {
	url := c.explorerUrls.Tx[c.network] + string(txID)
	if c.network == testnet {
		url += explorerTestnetPostfix
	}
	return url
}
func (c *client) ValidateAddress(address common.Address) bool {
	if !strings.HasPrefix(string(address), prefixes[c.network].AccAddr) {
		return false
	}
	return common.ValidateAddress(address)
}
func (c *client) GetAddress() (common.Address, error) {
	return c.getAddressForWalletIndex(0)
}
func (c *client) getAddressForWalletIndex(walletIndex int) (common.Address, error) {
	if c.phrase == "" {
		return "", fmt.Errorf("client not initialized")
	}
	accAddr := c.bridge.GetContext().GetFromAddress()
    return common.Address(accAddr.String()), nil
}
func (c *client) SetPhrase(phrase string) (common.Address, error) {
	return c.setPhraseForWalletIndex(phrase, 0)
}
func (c *client) setPhraseForWalletIndex(phrase string, walletIndex int) (common.Address, error) {
    if (c.phrase != phrase) {
		c.phrase = phrase
		c.updateBridge()
	}
	return c.getAddressForWalletIndex(walletIndex)
}

// Balance Account definition
type THORBalanceAccount struct {
	Height	string				`json:"height"`
	Result	[]THORTokenBalance	`json:"result"`
}
type THORTokenBalance struct {
	Denom	string	`json:"denom"`
	Amount	string	`json:"amount"`
}
// Account definition
type THORAccount struct {
	Height	string          	`json:"height"`
	Result	THORAccountResult `json:"result"`
}
type THORAccountResult struct {
	Type	string				`json:"type"`
	Value	THORAccountValue	`json:"value"`
}
type THORAccountValue struct {
	Address			string			`json:"address"`
	PublicKey		THORTxPubKey	`json:"public_key"`
	AccountNumber	string			`json:"account_number"`
	Sequence		string			`json:"sequence"`
}
type THORTxPubKey struct {
	Type  string `json:"type"`  // "tendermint/PubKeySecp256k1",
	Value string `json:"value"` //base64_pubkey
}
	
func (c *client) GetBalance(address common.Address, assets []common.Asset) (common.Balances, error) {
	saddress := string(address)
	if saddress == "" {
		if c.bridge == nil {
			return nil, common.AddressMissingError
		}
		saddress = c.bridge.GetContext().GetFromAddress().String()
	}
	h := common.NewHttpClient(c.getFullNode())
	resp, code, err := h.Get(balanceEndpoint + saddress)
	if err != nil {
		if code == http.StatusNotFound {
			return common.Balances{}, nil
		}
		return common.Balances{}, err
	}
	var tcbal THORBalanceAccount
	if err := json.Unmarshal(resp, &tcbal); err != nil {
		return nil, err
	}
	xbal := thorToXBalance(tcbal)
	return xbal, nil	
}
func thorToXBalance(tcbal THORBalanceAccount) common.Balances {
	xbal := make(common.Balances, 0, len(tcbal.Result))
	for _, tb := range tcbal.Result {
		xb := common.Balance{
			Asset:	newAsset(tb.Denom),
			Amount:	common.NewUintFromFx8String(tb.Amount),
		}
		xbal = append(xbal, xb)
	}
	return xbal
}
func (c *client) GetAccount(address string) (*THORAccount, error) {
	if address == "" {
		if c.bridge == nil {
			return nil, common.AddressMissingError
		}
		address = string(c.bridge.GetContext().GetFromAddress())
	}
	h := common.NewHttpClient(c.getFullNode())
	resp, code, err := h.Get(accountEndpoint + address)
	if err != nil {
		if code == http.StatusNotFound {
			return &THORAccount{}, nil
		}
		return nil, err
	}
	var account THORAccount
	if err := json.Unmarshal(resp, &account); err != nil {
		return nil, err
	}
	return &account, nil
}

func (c *client) GetTransactions(params common.TxHistoryParams) (common.TxsPage, error) {
	panic("not implemented")
}
func (c *client) GetTransactionData(txId string, assetAddress common.Address) (common.Tx, error) {
	panic("not implemented")
}
func (c *client) GetFees(feeParams interface{}) common.Fees {
	//params := feeParams.(common.FeesParams)
	return defaultFees
}
func (c *client) Transfer(params common.TxParams) (common.TxHash, error) {
	if !c.initialized() {
		return "", fmt.Errorf("ThorchainClient struct not initialized")
	}
	return bridgeTransfer(c.bridge, params)
}
func (c *client) PurgeClient() {
	c.phrase = ""
	c.bridge = nil
}
func (c *client) SetClientUrl(clientUrl ClientUrl) {
	c.clientUrl = clientUrl
}
func (c *client) GetClientUrl() NodeUrl {
	return c.clientUrl[c.network]
}
func (c *client) Deposit(params DepositParams) (common.TxHash, error) {
	if !c.initialized() {
		return "", fmt.Errorf("ThorchainClient struct not initialized")
	}
	return bridgeDeposit(c.bridge, params);
}
func (c *client) getNode() string {
	return c.clientUrl[c.network].Node
}
func (c *client) getFullNode() string {
	return protocol + c.clientUrl[c.network].Node
}
func (c *client) getRpc() string {
	return c.clientUrl[c.network].Rpc
}
func (c *client) initialized() bool {
	return c.phrase != "" && c.bridge != nil
}

func newAsset(s string) common.Asset {
	asset, err := common.NewChainAsset(chainStr, s)
	if err != nil {
		panic(err)
	}
	return asset
}
