package mexbt

import (
	"fmt"
	"os"
	"testing"
)

// Sanity check suite
// Just touches server's nipples and check if server likes it

// Sandbox keys
// Use my 1000 bitcoins as you wish
var config = Config{
	ApiKey:     "2f0873b7875e1bcdef5677ffcf600dc4",
	PrivateKey: "fd836ed10efbdd9617a7223e59b2fd8c",
	UserId:     "der@2-47.ru",
}

func TestMain(m *testing.M) {
	Sandbox = true
	os.Exit(m.Run())
}

func check(t *testing.T, ok bool, err error) {
	if err != nil || !ok {
		fmt.Printf("%+v\n", t)
		t.Fail()
	}
}

func TestTicker(t *testing.T) {

	res, e := Ticker(BTCMXN)
	check(t, res.IsAccepted, e)
}

func TestTradesByDate(t *testing.T) {

	res, e := TradesByDate(BTCMXN, 141530012, 1416559390)
	check(t, res.IsAccepted, e)
}

func TestTrades(t *testing.T) {
	res, e := Trades(BTCMXN, -1, 10)
	check(t, res.IsAccepted, e)
}

func TestMe(t *testing.T) {
	r, e := config.Me()
	check(t, r.IsAccepted, e)
}

func TestBalance(t *testing.T) {
	r, e := config.Balance()
	check(t, r.IsAccepted, e)
}

func TestAccountTrades(t *testing.T) {
	r, e := config.AccountTrades(BTCUSD, -1, 20)
	check(t, r.IsAccepted, e)
}

func TestDepositAdresses(t *testing.T) {
	r, e := config.DepositAddresses()
	check(t, r.IsAccepted, e)

	btc, found := r.Get("BTC")
	fmt.Println("BTC balance", btc)
	check(t, found, nil)
}

func TestAccountOrders(t *testing.T) {
	r, e := config.AccountOrders()
	check(t, r.IsAccepted, e)
}

func TestProductPairs(t *testing.T) {
	r, e := ProductPairs()
	check(t, r.IsAccepted, e)
}

func TestOrderbook(t *testing.T) {
	r, e := Orderbook(BTCUSD)
	check(t, r.IsAccepted, e)
}

func TestCreateModifyCancel(t *testing.T) {
	r, e := config.CreateLimitOrder(BTCMXN, SIDE_BUY, 12.42, 2.02)
	check(t, r.IsAccepted, e)
	orderId := r.ServerOrderId

	r2, e := config.MoveToTop(BTCMXN, orderId)
	check(t, r2.IsAccepted, e)
	orderId = r2.ServerOrderId

	r3, e := config.Cancel(BTCMXN, orderId)
	check(t, r3.IsAccepted, e)
}
