package thorchain

import (
	"fmt"
	"testing"
	. "gopkg.in/check.v1"
	"github.com/zlyzol/xchaingo/common"
)
/*
func TestHttpGet(t *testing.T) {
	address := "tthor1fs5jqvwp9u05802vfsru8zndmq5ucanrw8gg96"
	c := NewHttpClient("https://testnet.thornode.thorchain.info")
	_, _, err := c.Get("/bank/balances/" + string(address))
	assert.NoError(t, err, "NewHttpClient")
}
*/
//func TestSetParameters(t *testing.T) {

// Hook up gocheck into the "go test" runner.
func TestAll(t *testing.T) { TestingT(t) }

type TestSuite struct{}
var _ = Suite(&TestSuite{})

func (s *TestSuite) TestSetParameters(c *C) {
	t, err := NewThorchainClient(testnet, "flip portion grant body mad mountain infant edit pig execute tired ridge")
	c.Assert(err, IsNil)
	c.Log("NewThorchainClient OK")
	x := t.GetXChainClient()

	x.SetNetwork(mainnet)
	c.Check(x.GetNetwork(), Equals, mainnet)
	x.SetNetwork(testnet)
	c.Check(x.GetNetwork(), Equals, testnet)
	c.Log("SetNetwork OK")

	c.Check(x.GetExplorerUrl(), Equals, "https://viewblock.io/thorchain")
	x.SetNetwork(testnet)
	a := "tthor1249ujrfl6pnhzxarwhqxpfu3k53hrndax2xs59"
	c.Check(x.GetExplorerAddressUrl(a), Equals, "https://viewblock.io/thorchain/address/"+a+"?network=testnet")
	tx := "D820AD9F3C1580EF8DD5DEB0909DC33520EF6793DA59EFCD666C6A5E8E6CA1DA"
	c.Check(x.GetExplorerTxUrl(tx), Equals, "https://viewblock.io/thorchain/tx/"+tx+"?network=testnet")
	c.Log("GetExplorerUrl & co OK")

	c.Check(x.GetFees(nil), Equals, defaultFees)
}

func (s *TestSuite) TestValidateAddress(c *C) {
	t, err := NewThorchainClient(testnet, "flip portion grant body mad mountain infant edit pig execute tired ridge")
	c.Assert(err, IsNil)
	x := t.GetXChainClient()
	a := common.Address("tthor1249ujrfl6pnhzxarwhqxpfu3k53hrndax2xs59")
	c.Check(x.ValidateAddress(a), Equals, true)
}

func (s *TestSuite) TestGetTestnetAddress(c *C) {
	t, err := NewThorchainClient(testnet, "penalty guard brown luxury move bar wrong hero trick update grow bitter") // random passphrase
	c.Assert(err, IsNil)
	x := t.GetXChainClient()
	address1, err := x.SetPhrase("flip portion grant body mad mountain infant edit pig execute tired ridge")
	c.Assert(err, IsNil)
	address2, err := x.GetAddress()
	c.Assert(err, IsNil)
	c.Check(string(address1), Equals, "tthor1249ujrfl6pnhzxarwhqxpfu3k53hrndax2xs59")
	c.Check(string(address1), Equals, string(address2))
}
func (s *TestSuite) TestGetMainnetAddress(c *C) {
	t, err := NewThorchainClient(mainnet, "penalty guard brown luxury move bar wrong hero trick update grow bitter") // random passphrase
	c.Assert(err, IsNil)
	x := t.GetXChainClient()
	address1, err := x.SetPhrase("flip portion grant body mad mountain infant edit pig execute tired ridge")
	c.Assert(err, IsNil)
	address2, err := x.GetAddress()
	c.Assert(err, IsNil)
	c.Check(string(address1), Equals, "thor1249ujrfl6pnhzxarwhqxpfu3k53hrndazahqdq")
	c.Check(string(address1), Equals, string(address2))
}
func (s *TestSuite) TestGetMainnetBalance(c *C) {
	t, err := NewThorchainClient(mainnet, "flip portion grant body mad mountain infant edit pig execute tired ridge") // random passphrase
	c.Assert(err, IsNil)
	x := t.GetXChainClient()
	address, err := x.GetAddress()
	c.Assert(err, IsNil)
	balances, err := x.GetBalance(address, nil)
	c.Assert(err, IsNil)
	c.Log(fmt.Sprintf("Mainnet balances: %+v", balances))
	c.Assert(len(balances) > 0, Equals, true)
	c.Check(balances[0].Asset.String(), Equals, common.RuneAsset().String())
	c.Check(balances[0].Amount.String(), Equals, common.NewUintFromFloat(1).String())
}
func (s *TestSuite) TestGetTestnetBalance(c *C) {
	t, err := NewThorchainClient(testnet, "flip portion grant body mad mountain infant edit pig execute tired ridge") // random passphrase
	c.Assert(err, IsNil)
	x := t.GetXChainClient()
	address, err := x.GetAddress()
	c.Assert(err, IsNil)
	balances, err := x.GetBalance(address, nil)
	c.Assert(err, IsNil)
	c.Log(fmt.Sprintf("Testnet balances: %+v", balances))
	c.Assert(len(balances) > 0, Equals, true)
	c.Check(balances[0].Asset.String(), Equals, common.RuneAsset().String())
}
func (s *TestSuite) xTestTransfer(c *C) {
	t, err := NewThorchainClient(testnet, "flip portion grant body mad mountain infant edit pig execute tired ridge") // random passphrase
	c.Assert(err, IsNil)
	x := t.GetXChainClient()
	params := common.TxParams{
		Asset:		AssetRune,
		Amount:		common.NewUintFromFloat(0.1),
		Recipient:	"tthor1fnzhjcnqn33fahdywf7t5azcapjry83r3h5j3g",
	}
	hash, err := x.Transfer(params)
	c.Assert(err, IsNil)
	c.Assert(len(hash) > 0, Equals, true)
}
func (s *TestSuite) TestPurgeClient(c *C) {
	t, err := NewThorchainClient(testnet, "flip portion grant body mad mountain infant edit pig execute tired ridge") // random passphrase
	c.Assert(err, IsNil)
	x := t.GetXChainClient()
	params := common.TxParams{
		Asset:		AssetRune,
		Amount:		common.NewUintFromFloat(0.1),
		Recipient:	"tthor1fnzhjcnqn33fahdywf7t5azcapjry83r3h5j3g",
	}
	x.PurgeClient()
	_, err = x.Transfer(params)
	c.Assert(err, NotNil)
}
