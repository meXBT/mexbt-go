// meXBT API client
package mexbt

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	// CreateOrder() sides
	SIDE_BUY  = "buy"
	SIDE_SELL = "sell"

	// Product pairs
	BTCUSD = "BTCUSD"
	BTCMXN = "BTCMXN"

	// CreateOrder() types
	ORDER_TYPE_MARKET = 1
	ORDER_TYPE_LIMIT  = 0

	// Modify() modes
	MODIFY_BUMP    = 0
	MODIFY_EXECUTE = 1
)

var (
	// Set to true to connect to sandbox API
	Sandbox = false
)

type tickerRequest struct {
	Pair string `json:"productPair"`
}

type request struct {
}

type tradesByDateRequest struct {
	Ins       string `json:"ins"`
	StartDate int    `json:"startDate"`
	EndDate   int    `json:"endDate"`
}

type tradesRequest struct {
	Ins        string `json:"ins"`
	StartIndex int    `json:"startIndex"`
	Count      int    `json:"count"`
}

type accountTradesRequest struct {
	StartIndex int `json:"startIndex"`
	Count      int `json:"count"`
	signedRequest
}

type orderbookRequest struct {
	ProductPair string `json:"productPair"`
}

type signedRequest struct {
	ApiKey   string `json:"apiKey"`
	ApiNonce int64  `json:"apiNonce"`
	ApiSig   string `json:"apiSig"`
	Ins      string `json:"ins"`
}

type createOrderRequest struct {
	signedRequest
	Side      string  `json:"side"`
	OrderType int     `json:"orderType"`
	Qty       float64 `json:"qty,string"`
	Px        float64 `json:"px,string"`
}

type withdrawRequest struct {
	signedRequest
	Amount        float64 `json:"amount,string"`
	SendToAddress string  `json:"sendToAdress"`
}

type modifyOrderRequest struct {
	signedRequest
	ServerOrderId int64 `json:"serverOrderId"`
	ModifyAction  int   `json:"modifyAction"`
}

type cancelOrderRequest struct {
	signedRequest
	ServerOrderId int64 `json:"serverOrderId"`
	ModifyAction  int   `json:"modifyAction"`
}

type Result struct {
	IsAccepted   bool
	RejectReason string `json:",omitempty"`
}

type TickerResult struct {
	High                    float32
	Last                    float32
	Bid                     float32
	Volume                  float32
	Volume24h               float32
	Volume24hrProduct2      float32
	Low                     float32
	Ask                     float32
	Total24HrQtyTraded      float32
	Total24HrProduct2Traded float32
	Total24HrNumTrades      int
	SellOrderCount          int
	BuyOrderCount           int
	NumOfCreateOrders       int
	Result
}

type TradesByDateResult struct {
	DateTimeUtc int64
	Ins         string
	StartDate   int
	EndDate     int
	Trades      []Trade
	Result
}

type TradesResult struct {
	DateTimeUtc int64
	Ins         string
	StartIndex  int
	Count       int
	Trades      []Trade
	Result
}

type Trade struct {
	Tid int
	Order
	Unixtime              int
	UtcTicks              int64
	IncomingOrderSide     int
	IncomingServerOrderId int
	BookServerOrderId     int
}

type OrderbookResult struct {
	Bids []Order
	Asks []Order
	Result
}

type Order struct {
	Qty float64
	Px  float64 `json:",omitempty"`
}

type Config struct {
	ApiKey     string
	PrivateKey string
	UserId     string
}

type CreateOrderResult struct {
	ServerOrderId int64
	DateTimeUtc   int64
	Result
}

type BalanceResult struct {
	Currencies []BalanceEntry
	Result
}

type BalanceEntry struct {
	Name       string
	Balance    float64
	Hold       int
	TradeCount int
}

type UserInformationResult struct {
	UserInfoKVP []KeyValuePair
	Result
}

type KeyValuePair struct {
	Key   string
	Value string
}

type AccountOrdersResult struct {
	DateTimeUtc    int64
	OpenOrdersInfo []InsOrders
	Result
}

type InsOrders struct {
	Ins        string
	OpenOrders []OpenOrder
}

type OpenOrder struct {
	ServerOrderId int64
	AccountId     int
	Price         int
	QtyTotal      int
	QtyRemaining  int
	RecieveTime   int64
	Side          int
}

type DepositAddressesResult struct {
	Addresses []Address
	Result
}

type Address struct {
	Name           string
	DepositAddress string
}

type ProductPairsResult struct {
	ProductPairs []ProductPair
	Result
}

type ProductPair struct {
	Name                  string
	ProductPairCode       int
	Product1Label         string
	Product1DecimalPlaces int
	Product2Label         string
	Product2DecimalPlaces int
}

type ModifyResult struct {
	ServerOrderId int64
	DateTimeUtc   int64
	Result
}

func getApiURL(method string) string {
	api := strings.SplitN(method, "/", 2)
	if Sandbox && api[0] != "public" {
		return "https://" + api[0] + "-api-sandbox.mexbt.com/v1/" + api[1]
	} else {
		return "https://" + api[0] + "-api.mexbt.com/v1/" + api[1]
	}
}

func (c *Config) makeRequest(ins string) (r signedRequest) {
	r.Ins = ins
	r.ApiNonce = time.Now().UnixNano() / int64(1000000)
	r.ApiKey = c.ApiKey
	hasher := hmac.New(sha256.New, []byte(c.PrivateKey))

	hasher.Write([]byte(strconv.FormatInt(r.ApiNonce, 10) + c.UserId + c.ApiKey))
	r.ApiSig = strings.ToUpper(hex.EncodeToString(hasher.Sum(nil)))

	return
}

func makeRequest(method string, r interface{}) ([]byte, error) {
	jsonStr, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}

	resp, err := http.Post(getApiURL(method), "application/json", bytes.NewBuffer(jsonStr))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func genericRequest(api string, request, result interface{}) error {

	body, err := makeRequest(api, request)
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, result)
	if err != nil {
		return err
	}

	return nil
}

func Ticker(pair string) (*TickerResult, error) {

	request := tickerRequest{Pair: pair}
	result := TickerResult{}

	err := genericRequest("public/ticker", &request, &result)
	return &result, err
}

func ProductPairs() (*ProductPairsResult, error) {

	request := request{}
	result := ProductPairsResult{}

	err := genericRequest("public/product-pairs", &request, &result)
	return &result, err
}

func TradesByDate(ins string, startDate int, endDate int) (*TradesByDateResult, error) {

	request := tradesByDateRequest{Ins: ins, StartDate: startDate, EndDate: endDate}
	result := TradesByDateResult{}
	err := genericRequest("public/trades-by-date", &request, &result)
	return &result, err
}

func Trades(ins string, startIndex int, count int) (*TradesResult, error) {
	request := tradesRequest{Ins: ins, StartIndex: startIndex, Count: count}
	result := TradesResult{}

	err := genericRequest("public/trades", &request, &result)
	return &result, err
}

func Orderbook(productPair string) (OrderbookResult, error) {

	request := orderbookRequest{ProductPair: productPair}
	result := OrderbookResult{}

	err := genericRequest("public/order-book", &request, &result)

	return result, err
}

func (c *Config) CreateMarketOrder(productPair string, side string, qty float64) (CreateOrderResult, error) {
	return c.CreateOrder(productPair, side, qty, 0, ORDER_TYPE_MARKET)
}

func (c *Config) CreateLimitOrder(productPair string, side string, qty float64, px float64) (CreateOrderResult, error) {
	return c.CreateOrder(productPair, side, qty, px, ORDER_TYPE_LIMIT)
}

func (c *Config) CreateOrder(productPair string, side string, qty float64, px float64, orderType int) (CreateOrderResult, error) {

	request := createOrderRequest{
		Qty: qty, Px: px, OrderType: orderType, Side: side,
		signedRequest: c.makeRequest(productPair),
	}
	result := CreateOrderResult{}

	err := genericRequest("private/orders/create", &request, &result)

	return result, err
}

func (c *Config) MoveToTop(ins string, serverOrderId int64) (ModifyResult, error) {
	return c.Modify(ins, serverOrderId, MODIFY_BUMP)
}

func (c *Config) ExecuteNow(ins string, serverOrderId int64) (ModifyResult, error) {
	return c.Modify(ins, serverOrderId, MODIFY_EXECUTE)
}

func (c *Config) Modify(ins string, serverOrderId int64, action int) (ModifyResult, error) {

	request := modifyOrderRequest{
		ServerOrderId: serverOrderId, ModifyAction: action,
		signedRequest: c.makeRequest(ins),
	}
	result := ModifyResult{}

	err := genericRequest("private/orders/modify", &request, &result)

	return result, err
}

func (c *Config) Cancel(ins string, serverOrderId int64) (ModifyResult, error) {

	request := cancelOrderRequest{
		ServerOrderId: serverOrderId,
		signedRequest: c.makeRequest(ins),
	}
	result := ModifyResult{}

	err := genericRequest("private/orders/cancel", &request, &result)

	return result, err
}

func (c *Config) CancelAll(ins string) (Result, error) {

	request := c.makeRequest(ins)
	result := Result{}

	err := genericRequest("private/orders/cancel-all", &request, &result)

	return result, err
}

func (c *Config) Balance() (BalanceResult, error) {

	request := c.makeRequest("")
	result := BalanceResult{}

	err := genericRequest("private/balance", &request, &result)

	return result, err
}

func (c *Config) Me() (UserInformationResult, error) {

	request := c.makeRequest("")
	result := UserInformationResult{}

	err := genericRequest("private/me", &request, &result)

	return result, err
}

func (c *Config) AccountTrades(ins string, startIndex int, count int) (*TradesResult, error) {

	request := accountTradesRequest{
		StartIndex: startIndex, Count: count,
		signedRequest: c.makeRequest(ins),
	}
	result := TradesResult{}

	err := genericRequest("private/trades", &request, &result)
	return &result, err
}

func (c *Config) AccountOrders() (*AccountOrdersResult, error) {

	request := c.makeRequest("")
	result := AccountOrdersResult{}

	err := genericRequest("private/orders", &request, &result)
	return &result, err
}

func (c *Config) DepositAddresses() (*DepositAddressesResult, error) {

	request := c.makeRequest("")
	result := DepositAddressesResult{}

	err := genericRequest("private/deposit-addresses", &request, &result)
	return &result, err
}

// TODO: check on non-sandbox API, because it does not answer this request
func (c *Config) Withdraw(ins string, amount float64, address string) (*Result, error) {

	request := withdrawRequest{
		Amount: amount, SendToAddress: address,
		signedRequest: c.makeRequest(""),
	}

	request.Ins = ins

	result := Result{}

	err := genericRequest("private/withdraw", &request, &result)
	return &result, err
}
