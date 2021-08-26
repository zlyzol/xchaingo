package binance

import (
	"fmt"
	"testing"
	. "gopkg.in/check.v1"
	"github.com/zlyzol/xchaingo/common"
)

// Hook up gocheck into the "go test" runner.
func TestAll(t *testing.T) { TestingT(t) }

type TestSuite struct{}
var _ = Suite(&TestSuite{})

func (s *TestSuite) TestSetParameters(c *C) {
	t, err := NewBinanceClient(testnet, "flip portion grant body mad mountain infant edit pig execute tired ridge")
	c.Assert(err, IsNil)
	c.Log("NewBinanceClient OK")
	x := t.GetXChainClient()

	x.SetNetwork(mainnet)
	c.Check(x.GetNetwork(), Equals, mainnet)
	x.SetNetwork(testnet)
	c.Check(x.GetNetwork(), Equals, testnet)
	c.Log("SetNetwork OK")

	c.Check(x.GetExplorerUrl(), Equals, "https://testnet-explorer.binance.org")
	x.SetNetwork(testnet)
	a := "tbnb174ht75un03ddmsumaas9z2tgmv723el5wax0qr"
	c.Check(x.GetExplorerAddressUrl(a), Equals, "https://testnet-explorer.binance.org/address/"+a)
	tx := "D820AD9F3C1580EF8DD5DEB0909DC33520EF6793DA59EFCD666C6A5E8E6CA1DA"
	c.Check(x.GetExplorerTxUrl(tx), Equals, "https://testnet-explorer.binance.org/tx/"+tx)
	c.Log("GetExplorerUrl & co OK")

	c.Check(x.GetFees(nil), Equals, defaultFees)
}

func (s *TestSuite) TestValidateAddress(c *C) {
	t, err := NewBinanceClient(testnet, "flip portion grant body mad mountain infant edit pig execute tired ridge")
	c.Assert(err, IsNil)
	x := t.GetXChainClient()
	a := common.Address("tbnb174ht75un03ddmsumaas9z2tgmv723el5wax0qr")
	c.Check(x.ValidateAddress(a), Equals, true)
}

func (s *TestSuite) xTestGetTestnetAddress(c *C) {
	t, err := NewBinanceClient(testnet, "penalty guard brown luxury move bar wrong hero trick update grow bitter") // random passphrase
	c.Assert(err, IsNil)
	x := t.GetXChainClient()
	address1, err := x.SetPhrase("flip portion grant body mad mountain infant edit pig execute tired ridge")
	c.Assert(err, IsNil)
	address2, err := x.GetAddress()
	c.Assert(err, IsNil)
	c.Check(string(address1), Equals, "tbnb174ht75un03ddmsumaas9z2tgmv723el5wax0qr")
	c.Check(string(address1), Equals, string(address2))
}
func (s *TestSuite) TestGetMainnetAddress(c *C) {
	t, err := NewBinanceClient(mainnet, "penalty guard brown luxury move bar wrong hero trick update grow bitter") // random passphrase
	c.Assert(err, IsNil)
	x := t.GetXChainClient()
	address1, err := x.SetPhrase("flip portion grant body mad mountain infant edit pig execute tired ridge")
	c.Assert(err, IsNil)
	address2, err := x.GetAddress()
	c.Assert(err, IsNil)
	c.Check(string(address1), Equals, "bnb174ht75un03ddmsumaas9z2tgmv723el5qg0tqj")
	c.Check(string(address1), Equals, string(address2))
}
func (s *TestSuite) TestGetMainnetBalance(c *C) {
	t, err := NewBinanceClient(mainnet, "flip portion grant body mad mountain infant edit pig execute tired ridge") // random passphrase
	c.Assert(err, IsNil)
	x := t.GetXChainClient()
	address, err := x.GetAddress()
	c.Assert(err, IsNil)
	balances, err := x.GetBalance(address, nil)
	c.Log(fmt.Sprintf("Mainnet balances: %+v", balances))
	c.Assert(len(balances) > 0, Equals, true)
	c.Check(balances[0].Asset.String(), Equals, AssetBNB.String())
	c.Check(balances[0].Amount.String(), Equals, common.NewUint(5000000).String())
}
func (s *TestSuite) TestGetTestnetBalance(c *C) {
	t, err := NewBinanceClient(testnet, "flip portion grant body mad mountain infant edit pig execute tired ridge") // random passphrase
	c.Assert(err, IsNil)
	x := t.GetXChainClient()
	address, err := x.GetAddress()
	c.Assert(err, IsNil)
	balances, err := x.GetBalance(address, nil)
	c.Log(fmt.Sprintf("Testnet balances: %+v", balances))
	c.Assert(len(balances) > 0, Equals, true)
	c.Check(balances[0].Asset.String(), Equals, AssetBNB.String())
}
func (s *TestSuite) TestTransfer(c *C) {
	t, err := NewBinanceClient(testnet, "flip portion grant body mad mountain infant edit pig execute tired ridge") // random passphrase
	c.Assert(err, IsNil)
	x := t.GetXChainClient()
	params := common.TxParams{
		Asset:		AssetBNB,
		Amount:		common.NewUintFromFloat(0.01),
		Recipient:	"tbnb12c9lnfdryjdc9fg88a4lj5kd98agjcttk7jpa7",
		Memo:		"xchaingo test memo",
	}
	hash, err := x.Transfer(params)
	c.Assert(err, IsNil)
	c.Assert(len(hash) > 0, Equals, true)
}
func (s *TestSuite) TestPurgeClient(c *C) {
	t, err := NewBinanceClient(testnet, "flip portion grant body mad mountain infant edit pig execute tired ridge") // random passphrase
	c.Assert(err, IsNil)
	x := t.GetXChainClient()
	params := common.TxParams{
		Asset:		AssetBNB,
		Amount:		common.NewUintFromFloat(0.01),
		Recipient:	"tbnb12c9lnfdryjdc9fg88a4lj5kd98agjcttk7jpa7",
		Memo:		"xchaingo test memo",
	}
	x.PurgeClient()
	_, err = x.Transfer(params)
	c.Assert(err, NotNil)
}
