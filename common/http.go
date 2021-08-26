package common

import (
	//	"encoding/json"
	"fmt"
	"net/http"

	//	"net/url"
	"time"

	"gopkg.in/resty.v1"
	//	"github.com/zlyzol/xchaingo/binance-chain/go-sdk/types"
	//	"github.com/zlyzol/xchaingo/binance-chain/go-sdk/types/tx"
	//	"github.com/gorilla/websocket"
	//	"github.com/zlyzol/xchaingo/xclient"
)

const (
	MaxReadWaitTime = 30 * time.Second
)

type HttpClient struct {
	nodeUrl string
}

func NewHttpClient(nodeUrl string) *HttpClient {
	return &HttpClient{
		nodeUrl: nodeUrl,
	}
}

func (c *HttpClient) Get(path string) ([]byte, int, error) {
	qp := map[string]string{}
	return c.GetQp(path, qp)
}

// original from binance api:
func (c *HttpClient) GetQp(path string, qp map[string]string) ([]byte, int, error) {
	resp, err := resty.R().SetQueryParams(qp).Get(c.nodeUrl + path)
	if err != nil {
		return nil, 0, err
	}
	if resp.StatusCode() >= http.StatusMultipleChoices || resp.StatusCode() < http.StatusOK {
		err = fmt.Errorf("bad response, status code %d, response: %s", resp.StatusCode(), string(resp.Body()))
	}
	return resp.Body(), resp.StatusCode(), err
}

// Post generic method
func (c *HttpClient) Post(path string, body interface{}, param map[string]string) ([]byte, error) {
	resp, err := resty.R().
		SetHeader("Content-Type", "text/plain").
		SetBody(body).
		SetQueryParams(param).
		Post(c.nodeUrl + path)

	if err != nil {
		return nil, err
	}
	if resp.StatusCode() >= http.StatusMultipleChoices {
		err = fmt.Errorf("bad response, status code %d, response: %s", resp.StatusCode(), string(resp.Body()))
	}
	return resp.Body(), err
}

/*
func (c *thorClient) GetAccount() (acc types.Account, err error) {
	{self.server}/bank/balances/{address}
}

func (c *thorClient) GetBalances(addr types.AccAddress) ([]types.TokenBalance, error) {
	account, err := c.GetAccount(addr)
	if err != nil {
		return nil, err
	}
	if account == nil {
		return []types.TokenBalance{}, nil
	}
	coins := account.GetCoins()

	symbs := make([]string, 0, len(coins))
	bals := make([]types.TokenBalance, 0, len(coins))
	for _, coin := range coins {
		symbs = append(symbs, coin.Denom)
		// count locked and frozen coins
		var locked, frozen int64
		nacc := account.(types.NamedAccount)
		if nacc != nil {
			locked = nacc.GetLockedCoins().AmountOf(coin.Denom)
			frozen = nacc.GetFrozenCoins().AmountOf(coin.Denom)
		}
		bals = append(bals, types.TokenBalance{
			Symbol: coin.Denom,
			Free:   types.Fixed8(coins.AmountOf(coin.Denom)),
			Locked: types.Fixed8(locked),
			Frozen: types.Fixed8(frozen),
		})
	}
	return bals, nil
}

func (c *thorClient) GetBalance(addr types.AccAddress, symbol string) (*types.TokenBalance, error) {
	if err := ValidateSymbol(symbol); err != nil {
		return nil, err
	}
	exist := c.existsCC(symbol)
	if !exist {
		return nil, errors.New("symbol not found")
	}
	acc, err := c.GetAccount(addr)
	if err != nil {
		return nil, err
	}
	if acc == nil {
		return &types.TokenBalance{
			Symbol: symbol,
			Free:   types.Fixed8Zero,
			Locked: types.Fixed8Zero,
			Frozen: types.Fixed8Zero,
		}, nil
	}
	var locked, frozen int64
	nacc := acc.(types.NamedAccount)
	if nacc != nil {
		locked = nacc.GetLockedCoins().AmountOf(symbol)
		frozen = nacc.GetFrozenCoins().AmountOf(symbol)
	}
	return &types.TokenBalance{
		Symbol: symbol,
		Free:   types.Fixed8(nacc.GetCoins().AmountOf(symbol)),
		Locked: types.Fixed8(locked),
		Frozen: types.Fixed8(frozen),
	}, nil
}

func (c *thorClient) GetFee() ([]types.FeeParam, error) {
	rawFee, err := c.ABCIQuery(fmt.Sprintf("%s/fees", ParamABCIPrefix), nil)
	if err != nil {
		return nil, err
	}
	if !rawFee.Response.IsOK() {
		return nil, fmt.Errorf(rawFee.Response.Log)
	}
	var fees []types.FeeParam
	err = c.cdc.UnmarshalBinaryLengthPrefixed(rawFee.Response.GetValue(), &fees)
	return fees, err
}

// GetTx returns transaction details
func (c *thorClient) GetTx(txHash string) (*tx.TxResult, error) {
	if txHash == "" {
		return nil, fmt.Errorf("Invalid tx hash %s ", txHash)
	}

	qp := map[string]string{}
	resp, _, err := c.Get("/tx/"+txHash, qp)
	if err != nil {
		return nil, err
	}

	var txResult tx.TxResult
	if err := json.Unmarshal(resp, &txResult); err != nil {
		return nil, err
	}

	return &txResult, nil
}

// PostTx returns transaction details
func (c *thorClient) PostTx(hexTx []byte, param map[string]string) ([]tx.TxCommitResult, error) {
	if len(hexTx) == 0 {
		return nil, fmt.Errorf("Invalid tx  %s", hexTx)
	}

	body := hexTx
	resp, err := c.Post("/broadcast", body, param)
	if err != nil {
		return nil, err
	}
	txResult := make([]tx.TxCommitResult, 0)
	if err := json.Unmarshal(resp, &txResult); err != nil {
		return nil, err
	}

	return txResult, nil
}

func (c *thorClient) WsGet(path string, constructMsg func([]byte) (interface{}, error), closeCh <-chan struct{}) (<-chan interface{}, error) {
	u := url.URL{Scheme: types.DefaultWSSchema, Host: c.baseUrl, Path: fmt.Sprintf("%s/%s", types.DefaultWSPrefix, path)}
	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		return nil, err
	}
	conn.SetPingHandler(nil)
	conn.SetPongHandler(
		func(string) error {
			conn.SetReadDeadline(time.Now().Add(MaxReadWaitTime))
			return nil
		})
	messages := make(chan interface{}, 0)
	finish := make(chan struct{}, 0)
	keepAliveCh := time.NewTicker(30 * time.Minute)
	pingTicker := time.NewTicker(10 * time.Second)
	go func() {
		defer conn.Close()
		defer close(messages)
		defer keepAliveCh.Stop()
		defer pingTicker.Stop()
		select {
		case <-closeCh:
			return
		case <-finish:
			return
		}
	}()
	go func() {
		writeMsg := func(m interface{}) bool {
			select {
			case <-closeCh:
				// already closed by user
				return true
			default:
			}
			messages <- m
			return false
		}
		for {
			select {
			case <-closeCh:
				conn.WriteControl(websocket.CloseMessage, nil, time.Now().Add(time.Second))
				return
			case <-keepAliveCh.C:
				conn.WriteJSON(&struct {
					Method string
				}{"keepAlive"})
			case <-pingTicker.C:
				conn.WriteControl(websocket.PingMessage, nil, time.Now().Add(time.Second))
			default:
				response := WSResponse{}
				err := conn.ReadJSON(&response)
				if err != nil {
					if closed := writeMsg(err); !closed {
						close(finish)
					}
					return
				}
				bz, err := json.Marshal(response.Data)
				if err != nil {
					if closed := writeMsg(err); !closed {
						close(finish)
					}
					return
				}
				msg, err := constructMsg(bz)
				if err != nil {
					if closed := writeMsg(err); !closed {
						close(finish)
					}
					return
				}
				//Todo delete condition when ws do not return account and order in the same time.
				if msg != nil {
					if closed := writeMsg(msg); closed {
						return
					}
				}
			}
		}
	}()
	return messages, nil
}

type WSResponse struct {
	Stream string
	Data   interface{}
}
*/
