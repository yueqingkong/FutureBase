package okex

import (
	"encoding/json"
	"fmt"
	"github.com/yueqingkong/FutureBase/base"
	"github.com/yueqingkong/FutureBase/util"
	"log"
	"strings"
	"time"
)

type Spot struct {
	Url       string
	Apikey    string
	SecretKey string
	PhraseKey string
}

var (
	ok *Spot
)

func NewSpot() *Spot {
	if ok == nil {
		ok = new(Spot)
		ok.Url = "https://www.okex.com"
	}
	return ok
}

// 初始化 Key
// api string, secret string, phrase string
func (ok *Spot) InitKeys(keys ...string) {
	ok.Apikey = keys[0]
	ok.SecretKey = keys[1]
	ok.PhraseKey = keys[2]
}

func (ok *Spot) Symbol(symbol base.SYMBOL) string {
	var v string
	switch symbol {
	case base.BTCUSDT:
		v = "BTC-USDT"
		break
	case base.ETHUSDT:
		v = "ETH-USDT"
		break
	case base.LTCUSDT:
		v = "LTC-USDT"
		break
	case base.EOSUSDT:
		v = "EOS-USDT"
		break
	case base.BTMUSDT:
		v = "BTM-USDT"
		break
	}
	return v
}

func (self *Spot) Period(period base.PERIOD) string {
	var v string
	switch period {
	case base.MIN_1:
		v = "60"
		break
	case base.MIN_3:
		v = "180"
		break
	case base.MIN_5:
		v = "300"
		break
	case base.MIN_15:
		v = "900"
		break
	case base.MIN_30:
		v = "1800"
		break
	case base.HOUR_1:
		v = "3600"
		break
	case base.HOUR_2:
		v = "7200"
		break
	case base.HOUR_4:
		v = "14400"
		break
	case base.HOUR_6:
		v = "21600"
		break
	case base.HOUR_12:
		v = "43200"
		break
	case base.DAY_1:
		v = "86400"
		break
	case base.WEEK_1:
		v = "604800"
		break
	}
	return v
}

func (self *Spot) Depth(depth base.DEPTH) string {
	var v string
	switch depth {
	case base.DEPTH_0:
		v = "0.1"
		break
	case base.DEPTH_10:
		v = "0.01"
		break
	case base.DEPTH_100:
		v = "0.001"
		break
	case base.DEPTH_1000:
		v = "0.0001"
		break
	case base.DEPTH_10000:
		v = "0.00001"
		break
	case base.DEPTH_100000:
		v = "0.000001"
		break
	}
	return v
}

// 订单类型
func (self *Spot) OrderType(order base.ORDER_TYPE) string {
	var v string
	switch order {
	case base.MARKET_BUY:
		v = "buy-market"
		break
	}
	return v
}

func (self *Spot) header(request string, path string, body interface{}) map[string]string {
	var paramString string
	if body != nil {
		bodyBytes, _ := json.Marshal(body)
		paramString = string(bodyBytes)
	}

	timestamp := util.IsoTime(time.Now())
	comnination := timestamp + strings.ToUpper(request) + path + paramString
	sign, err := util.HmacSha256Base64Signer(comnination, ok.SecretKey)
	if err != nil {
		log.Print("签名失败", err)
	}

	var headerMap = make(map[string]string, 0)
	headerMap[ACCEPT] = APPLICATION_JSON
	headerMap[CONTENT_TYPE] = APPLICATION_JSON_UTF8
	headerMap[COOKIE] = LOCALE + ENGLISH
	headerMap[OK_ACCESS_KEY] = ok.Apikey
	headerMap[OK_ACCESS_SIGN] = sign
	headerMap[OK_ACCESS_TIMESTAMP] = timestamp
	headerMap[OK_ACCESS_PASSPHRASE] = ok.PhraseKey
	return headerMap
}

// 获取币对信息
func (self *Spot) Instruments() Instrument {
	var api = "/api/spot/v3/instruments"
	var url = ok.Url + api
	var inst Instrument
	util.Get(url, nil, &inst)
	return inst
}

// 获取深度数据
func (self *Spot) Depths(symbol base.SYMBOL) Depth {
	var api = fmt.Sprintf("/api/spot/v3/instruments/%s/book", ok.Symbol(symbol))
	var url = ok.Url + api
	var depth Depth
	util.Get(url, nil, &depth)
	return depth
}

// 获取全部ticker信息
func (self *Spot) Tickers() []SpotTicker {
	var api = "/api/spot/v3/instruments/ticker"
	var url = ok.Url + api
	var tickers []SpotTicker
	util.Get(url, nil, &tickers)
	return tickers
}

// 获取币对的最新成交价、买一价、卖一价和24小时交易量的快照信息。
func (self *Spot) Ticker(symbol base.SYMBOL) SpotTicker {
	var api = fmt.Sprintf("/api/spot/v3/instruments/%s/ticker", self.Symbol(symbol))
	var url = ok.Url + api
	log.Println(url)

	var ticker SpotTicker
	util.Get(url, nil, &ticker)
	return ticker
}

// 获取币的当前最新价格
func (self *Spot) Price(contract base.CONTRACT_PERIOD, symbol base.SYMBOL) float32 {
	var api = fmt.Sprintf("/api/spot/v3/instruments/%s/ticker", self.Symbol(symbol))
	var url = ok.Url + api
	// log.Println(url)

	var ticker SpotTicker
	err := util.Get(url, nil, &ticker)

	var price float32
	if err == nil {
		price = util.StringToFloat32(ticker.Last)
	}
	return price
}

// 获取成交数据
func (self *Spot) Trades(symbol base.SYMBOL) []Trade {
	var api = fmt.Sprintf("/api/spot/v3/instruments/%s/trades", ok.Symbol(symbol))
	var url = ok.Url + api
	var trades []Trade
	util.Get(url, nil, &trades)
	return trades
}

// 获取币对的K线数据。K线数据按请求的粒度分组返回，k线数据最多可获取最近2000条
func (self *Spot) KLine(symbol base.SYMBOL, period base.PERIOD, st time.Time) SpotCandles {
	if st.IsZero() {
		var err error
		st, err = time.Parse("2006-01-02 15:04:05", "2017-08-17 00:00:00")
		if err != nil {
			log.Print(err)
		}
	}

	start := util.IsoTime(st)
	api := fmt.Sprintf("/api/spot/v3/instruments/%s/candles?granularity=%s&start=%s", self.Symbol(symbol), self.Period(period), start)
	var url = self.Url + api

	var candles SpotCandles
	err := util.Get(url, self.header("get", api, nil), &candles)
	if err != nil {
		log.Println(err)
	}
	return candles
}

// 币币账户信息
func (self *Spot) Accounts() []Account {
	var api = "/api/spot/v3/accounts"
	var url = ok.Url + api
	var accounts []Account
	util.Get(url, ok.header("get", api, nil), &accounts)
	return accounts
}

// 账单流水查询
func (self *Spot) Ledger(symbol string) []Ledger {
	var api = fmt.Sprintf("/api/spot/v3/accounts/%s/ledger", symbol)
	var url = ok.Url + api
	var ledgers []Ledger
	util.Get(url, ok.header("get", api, nil), &ledgers)
	return ledgers
}

// 下单
// type limit，market(默认是limit)，当以market（市价）下单，order_type只能选择0:普通委托
// side(buy | sell)
// order_type 0：普通委托（order type不填或填0都是普通委托） 1：只做Maker（Post only） 2：全部成交或立即取消（FOK） 3：立即成交并取消剩余（IOC）
// buy:  买入金额
// sell: 卖出数量
func (self *Spot) Order(symbol base.SYMBOL, side string, buy float32, sell float32) Order {
	var api = "/api/spot/v3/orders"
	var url = ok.Url + api

	param := make(map[string]interface{}, 0)
	param["type"] = "market" //默认市价单
	param["side"] = side
	param["instrument_id"] = ok.Symbol(symbol)
	param["order_type"] = 0
	param["margin_trading"] = 1

	if param["type"] == "market" { //市价单参数
		param["notional"] = buy
		param["size"] = sell
	} else if param["type"] == "limit" { //限价单单参数

	}

	var order Order
	util.Post(url, ok.header("post", api, param), param, &order)
	return order
}

// 撤销指定订单
func (self *Spot) CancelOrder(symbol base.SYMBOL, orderId string) CancelOrder {
	var api = fmt.Sprintf("/api/spot/v3/cancel_orders/%s", orderId)
	var url = ok.Url + api

	param := make(map[string]interface{}, 0)
	param["instrument_id"] = ok.Symbol(symbol)
	param["order_id"] = orderId

	var cancel CancelOrder
	util.Post(url, ok.header("post", api, param), param, &cancel)
	return cancel
}

// 获取订单列表
// state 订单状态("-2":失败,"-1":撤单成功,"0":等待成交 ,"1":部分成交, "2":完全成交,"3":下单中,"4":撤单中,"6": 未完成（等待成交+部分成交），"7":已完成（撤单成功+全部成交））
func (self *Spot) OrderList(symbol base.SYMBOL, state int32) []OrderItem {
	var api = fmt.Sprintf("/api/spot/v3/orders?instrument_id=%s&&state=%d", ok.Symbol(symbol), state)
	var url = ok.Url + api

	var item []OrderItem
	util.Get(url, ok.header("get", api, nil), &item)
	return item
}

// 获取所有未成交订单
func (self *Spot) OrderPending() []OrderPending {
	var api = "/api/spot/v3/orders_pending"
	var url = ok.Url + api

	var pending []OrderPending
	util.Get(url, ok.header("get", api, nil), &pending)
	return pending
}

// 获取订单信息
func (self *Spot) OrderInfo(symbol base.SYMBOL, orderId string) OrderInfo {
	var api = fmt.Sprintf("/api/spot/v3/orders/%s?instrument_id=%s", orderId, ok.Symbol(symbol))
	var url = ok.Url + api

	var info OrderInfo
	util.Get(url, ok.header("get", api, nil), &info)
	return info
}

// 获取成交明细
func (self *Spot) Deals(symbol base.SYMBOL, orderId string) Deal {
	var api = fmt.Sprintf("/api/spot/v3/fills?order_id=%s&instrument_id=%s", orderId, ok.Symbol(symbol))
	var url = ok.Url + api

	var deal Deal
	util.Get(url, ok.header("get", api, nil), &deal)
	return deal
}
