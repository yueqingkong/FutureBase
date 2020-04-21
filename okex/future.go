package okex

import (
	"encoding/json"
	"fmt"
	"gitee.com/shieldpu_futures/FutureBase/base"
	"gitee.com/shieldpu_futures/FutureBase/util"
	"log"
	"strings"
	"time"
)

// 底层简单的Api
type Future struct {
	Url       string
	Apikey    string
	SecretKey string
	PhraseKey string
}

var (
	future *Future
)

func NewFuture() *Future {
	if future == nil {
		future = new(Future)
		future.Url = "https://www.okex.com"
	}
	return future
}

// 初始化 Key
// api string, secret string, phrase string
func (self *Future) InitKeys(keys ...string) {
	self.Apikey = keys[0]
	self.SecretKey = keys[1]
	self.PhraseKey = keys[2]
}

func (self *Future) Period(period base.PERIOD) string {
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

func (self *Future) Depth(depth base.DEPTH) string {
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
func (self *Future) OrderType(order base.ORDER_TYPE) string {
	var v string
	switch order {
	case base.MARKET_BUY:
		v = "buy-market"
		break
	}
	return v
}

func (self Future) header(request string, path string, body interface{}) map[string]string {
	var paramString string
	if body != nil {
		bodyBytes, _ := json.Marshal(body)
		paramString = string(bodyBytes)
	}

	timestamp := util.IsoTime(time.Now())
	comnination := timestamp + strings.ToUpper(request) + path + paramString
	sign, err := util.HmacSha256Base64Signer(comnination, self.SecretKey)
	if err != nil {
		log.Print("签名失败", err)
	}

	var headerMap = make(map[string]string, 0)
	headerMap[ACCEPT] = APPLICATION_JSON
	headerMap[CONTENT_TYPE] = APPLICATION_JSON_UTF8
	headerMap[COOKIE] = LOCALE + ENGLISH
	headerMap[OK_ACCESS_KEY] = self.Apikey
	headerMap[OK_ACCESS_SIGN] = sign
	headerMap[OK_ACCESS_TIMESTAMP] = timestamp
	headerMap[OK_ACCESS_PASSPHRASE] = self.PhraseKey
	return headerMap
}

// 合约持仓信息
func (self Future) Position() FuturePosition {
	var api = "/api/futures/v3/position"
	var url = self.Url + api

	var position FuturePosition
	util.Get(url, self.header("get", api, nil), &position)
	return position
}

// 单个合约持仓信息
func (self Future) InstrumenPosition(instrumenid string) (FutureInstrumentPosition, error) {
	var api = fmt.Sprintf("/api/futures/v3/%s/position", instrumenid)
	var url = self.Url + api

	var position FutureInstrumentPosition
	err := util.Get(url, self.header("get", api, nil), &position)
	return position, err
}

// 所有币种合约账户信息
func (self Future) Account() (FutureAccount, error) {
	var api = "/api/futures/v3/accounts"
	var url = self.Url + api

	var account FutureAccount
	err := util.Get(url, self.header("get", api, nil), &account)
	return account, err
}

// 单个币种合约账户信息
func (self Future) SymbolAccount(symbol string) (FutureSymbolAccount, error) {
	var api = fmt.Sprintf("/api/futures/v3/accounts/%s", symbol)
	var url = self.Url + api

	var account FutureSymbolAccount
	err := util.Get(url, self.header("get", api, nil), &account)
	return account, err
}

// 获取合约币种杠杆倍数
func (self Future) Times(symbol string) FutureTimes {
	var api = fmt.Sprintf("/api/futures/v3/accounts/%s/leverage", symbol)
	var url = self.Url + api

	var times FutureTimes
	util.Get(url, self.header("get", api, nil), &times)
	return times
}

/**
 * 设定合约币种杠杆倍数
 * direct (short|long)
 * time 倍数(10|20)
 */
func (self Future) SetTimes(symbol string, id string, direct string, time int32) FutureTimes {
	var api = fmt.Sprintf("/api/futures/v3/accounts/%s/leverage", symbol)
	var url = self.Url + api

	param := make(map[string]interface{}, 0)
	param["leverage"] = time
	param["instrument_id"] = id
	param["direction"] = direct
	param["currency"] = symbol

	var times FutureTimes
	util.Post(url, self.header("post", api, param), param, &times)
	return times
}

// 账单流水查询
func (self Future) Ledger(symbol string) (FutureLedger, error) {
	var api = fmt.Sprintf("/api/futures/v3/accounts/%s/ledger", symbol)
	var url = self.Url + api

	var leger FutureLedger
	err := util.Post(url, self.header("get", api, nil), nil, &leger)
	return leger, err
}

/**
 * type(1:开多2:开空3:平多4:平空)
 * order_type: 0：普通委托（order type不填或填0都是普通委托） 1：只做Maker（Post only） 2：全部成交或立即取消（FOK） 3：立即成交并取消剩余（IOC）
 * match_price: 是否以对手价下单(0:不是 1:是)，默认为0，当取值为1时。price字段无效，当以对手价下单，order_type只能选择0:普通委托
 */
func (self Future) Order(instrumentid string, _type int32, ordertype int32, price float32, size int32, match_price int32) (FutureOrder, error) {
	var api = "/api/futures/v3/order"
	var url = self.Url + api

	param := make(map[string]interface{}, 0)
	param["instrument_id"] = instrumentid
	param["type"] = _type
	param["order_type"] = ordertype
	param["price"] = price
	param["size"] = size
	param["match_price"] = match_price

	var order FutureOrder
	err := util.Post(url, self.header("post", api, param), param, &order)
	return order, err
}

// 撤销指定订单
func (self Future) CancelOrder(symbol string, orderid string) (FutureCancel, error) {
	var api = fmt.Sprintf("/api/futures/v3/cancel_order/%s/%s", symbol, orderid)
	var url = self.Url + api

	param := make(map[string]interface{}, 0)
	param["instrument_id"] = symbol
	param["order_id"] = orderid

	var cancel FutureCancel
	err := util.Post(url, self.header("post", api, param), param, &cancel)
	return cancel, err
}

// 获取订单列表
// status (订单状态(-1.撤单成功；0:等待成交 1:部分成交 2:全部成交 6：未完成（等待成交+部分成交）7：已完成（撤单成功+全部成交))
func (self Future) List(symbol string, status int) FutureList {
	var api = fmt.Sprintf("/api/futures/v3/orders/%s", symbol)
	var url = self.Url + api

	var list FutureList
	util.Get(url, self.header("get", api, nil), &list)
	return list
}

// 获取订单信息
func (self Future) OrderInfo(symbol string, orderid string) FutureOrderInfo {
	var api = fmt.Sprintf("/api/futures/v3/orders/%s/%s", symbol, orderid)
	var url = self.Url + api

	var orderInfo FutureOrderInfo
	util.Get(url, self.header("get", api, nil), &orderInfo)
	return orderInfo
}

// 获取成交明细
func (self Future) Fills(symbol string, orderid string) FutureFill {
	var api = fmt.Sprintf("/api/futures/v3/fills?instrument_id=%s&order_id=%s", symbol, orderid)
	var url = self.Url + api

	var fills FutureFill
	util.Get(url, self.header("get", api, nil), &fills)
	return fills
}

// 获取合约信息
func (self Future) Instruments() []FutureInstrument {
	var api = "/api/futures/v3/instruments"
	var url = self.Url + api

	var instruments []FutureInstrument
	util.Get(url, self.header("get", api, nil), &instruments)
	return instruments
}

// 获取深度数据
func (self Future) Depths(symbol string) FutureDepth {
	var api = fmt.Sprintf("/api/futures/v3/instruments/%s/book", symbol)
	var url = self.Url + api

	var depth FutureDepth
	util.Get(url, self.header("get", api, nil), &depth)
	return depth
}

// 获取全部ticker信息
func (self Future) TickerAll() FutureTickers {
	var api = "/api/futures/v3/instruments/ticker"
	var url = self.Url + api

	var tickers FutureTickers
	util.Get(url, self.header("get", api, nil), &tickers)
	return tickers
}

// 获取某个ticker信息
func (self Future) Ticker(symbol string) FutureTicker {
	var api = fmt.Sprintf("/api/futures/v3/instruments/%s/ticker", symbol)
	var url = self.Url + api

	var ticker FutureTicker
	util.Get(url, self.header("get", api, nil), &ticker)
	return ticker
}

// 获取成交数据
func (self Future) Trades(symbol string, limit int32) ([]Trades, error) {
	var api = fmt.Sprintf("/api/futures/v3/instruments/%s/trades?limit=%d", symbol, limit)
	var url = self.Url + api

	var trades []Trades
	err := util.Get(url, self.header("get", api, nil), &trades)
	return trades, err
}

// 获取K线数据
// 合约历史记录不能回溯,只能拉取最近200条(cnmd,K线数据可能不完整)
// 60   180  300  900   1800  3600  7200  14400 21600 43200  86400 604800
// 1min 3min 5min 15min 30min 1hour 2hour 4hour 6hour 12hour 1day  1week
func (self Future) Candle(instrumentid string, interval int32, st time.Time) (FutureCandles, error) {
	var api string

	if st.IsZero() {
		var err error
		st, err = time.Parse("2006-01-02 15:04:05", "2017-08-17 00:00:00")
		if err != nil {
			log.Print(err)
		}
	}

	start := util.IsoTime(st)
	api = fmt.Sprintf("/api/futures/v3/instruments/%s/candles?start=%s&granularity=%d", instrumentid, start, interval)

	var url = self.Url + api

	var candles FutureCandles
	err := util.Get(url, self.header("get", api, nil), &candles)
	return candles, err
}

// 获取指数信息
func (okex Future) Index(symbol string) FutureIndex {
	var api = fmt.Sprintf("/api/futures/v3/instruments/%s/index", symbol)
	var url = okex.Url + api
	// log.Print(url)

	var index FutureIndex
	util.Get(url, okex.header("get", api, nil), &index)
	return index
}

// 获取法币汇率
func (self Future) Rate() FutureRate {
	var api = "/api/futures/v3/rate"
	var url = self.Url + api
	// log.Print(url)

	var rate FutureRate
	util.Get(url, self.header("get", api, nil), &rate)
	return rate
}

// 获取预估交割价
func (self Future) EstimatedPrice(symbol string) FutureEstimatedPrice {
	var api = fmt.Sprintf("/api/futures/v3/instruments/%s/estimated_price", symbol)
	var url = self.Url + api
	// log.Print(url)

	var price FutureEstimatedPrice
	util.Get(url, self.header("get", api, nil), &price)
	return price
}

// 获取平台总持仓量
func (self Future) OpenInterest(symbol string) FutureOpenInterest {
	var api = fmt.Sprintf("/api/futures/v3/instruments/%s/open_interest", symbol)
	var url = self.Url + api
	// log.Print(url)

	var interest FutureOpenInterest
	util.Get(url, self.header("get", api, nil), &interest)
	return interest
}

// 获取合约当前交易的最高买价和最低卖价
func (self *Future) Price(instrumentid string) float32 {
	var api = fmt.Sprintf("/api/futures/v3/instruments/%s/price_limit", instrumentid)
	var url = self.Url + api

	var price float32
	var limit FuturePriceLimit
	err := util.Get(url, self.header("get", api, nil), &limit)
	if err == nil {
		price = util.StringToFloat32(limit.Lowest)
	}

	return price
}

// 获取当前限价
func (self Future) MarkPrice(symbol string) FutureMarkPrice {
	var api = fmt.Sprintf("/api/futures/v3/instruments/%s/mark_price", symbol)
	var url = self.Url + api
	// log.Print(url)

	var price FutureMarkPrice
	util.Get(url, self.header("get", api, nil), &price)
	return price
}

// 获取爆仓单
func (self Future) Liquidation(symbol string, status int) FutureLiquidation {
	var api = fmt.Sprintf("/api/futures/v3/instruments/%s/liquidation?status=%d", symbol, status)
	var url = self.Url + api
	// log.Print(url)

	var liquidation FutureLiquidation
	util.Get(url, self.header("get", api, nil), &liquidation)
	return liquidation
}

// 获取合约挂单冻结数量
func (self Future) Holds(symbol string) FutureHold {
	var api = fmt.Sprintf("/api/futures/v3/accounts/%s/holds", symbol)
	var url = self.Url + api
	// log.Print(url)

	var hold FutureHold
	util.Get(url, self.header("get", api, nil), &hold)
	return hold
}

