package plat

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"
)

var (
	// 获取Token 交易对
	_symbolMap = map[string]string{
		"btc": "BTC-USDT",
		"eth": "ETHUSDT",
		"ltc": "LTCUSDT",
		"eos": "EOSUSDT",

		//"ae":  "AEBTC",
		//"bts": "BTSBTC",
		//"snt": "SNTBTC",
	}
)

const (
	CONTENT_TYPE = "Content-Type"
	ACCEPT       = "Accept"
	COOKIE       = "Cookie"
	LOCALE       = "locale="

	OK_ACCESS_KEY        = "OK-ACCESS-KEY"
	OK_ACCESS_SIGN       = "OK-ACCESS-SIGN"
	OK_ACCESS_TIMESTAMP  = "OK-ACCESS-TIMESTAMP"
	OK_ACCESS_PASSPHRASE = "OK-ACCESS-PASSPHRASE"

	APPLICATION_JSON      = "application/json"
	APPLICATION_JSON_UTF8 = "application/json; charset=UTF-8"

	ENGLISH = "en_US"
)

type OKexFuture struct {
}

func NewOKexFuture() OKexFuture {
	return OKexFuture{}
}


/**
 * 初始化 Key
 */
func (future OKexFuture) InitKeys(api string, secret string, phrase string) {
	ApiKey = api
	SecretKey = secret
	Passphrase = phrase
}

func (future OKexFuture) header(request string, path string, body interface{}) map[string]string {
	var paramString string
	if body != nil {
		bodyBytes, _ := json.Marshal(body)
		paramString = string(bodyBytes)
	}

	timestamp := IsoTime(time.Now())
	comnination := timestamp + strings.ToUpper(request) + path + paramString
	sign, err := HmacSha256Base64Signer(comnination, SecretKey)
	if err != nil {
		log.Print("签名失败", err)
	}

	var headerMap = make(map[string]string, 0)
	headerMap[ACCEPT] = APPLICATION_JSON
	headerMap[CONTENT_TYPE] = APPLICATION_JSON_UTF8
	headerMap[COOKIE] = LOCALE + ENGLISH
	headerMap[OK_ACCESS_KEY] = ApiKey
	headerMap[OK_ACCESS_SIGN] = sign
	headerMap[OK_ACCESS_TIMESTAMP] = timestamp
	headerMap[OK_ACCESS_PASSPHRASE] = Passphrase
	return headerMap
}

// 合约持仓信息
func (future OKexFuture) Position() FuturePosition {
	var api = "/api/futures/v3/position"
	var url = okApi + api

	var position FuturePosition
	Get(url, future.header("get", api, nil), &position)
	return position
}

// 单个合约持仓信息
func (future OKexFuture) InstrumenPosition(instrumenid string) (FutureInstrumentPosition, error) {
	var api = fmt.Sprintf("/api/futures/v3/%s/position", instrumenid)
	var url = okApi + api

	var position FutureInstrumentPosition
	err := Get(url, future.header("get", api, nil), &position)
	return position, err
}

// 所有币种合约账户信息
func (future OKexFuture) Account() (FutureAccount, error) {
	var api = "/api/futures/v3/accounts"
	var url = okApi + api

	var account FutureAccount
	err := Get(url, future.header("get", api, nil), &account)
	return account, err
}

// 单个币种合约账户信息
func (future OKexFuture) SymbolAccount(symbol string) (FutureSymbolAccount, error) {
	var api = fmt.Sprintf("/api/futures/v3/accounts/%s", symbol)
	var url = okApi + api

	var account FutureSymbolAccount
	err := Get(url, future.header("get", api, nil), &account)
	return account, err
}

// 获取合约币种杠杆倍数
func (future OKexFuture) Times(symbol string) FutureTimes {
	var api = fmt.Sprintf("/api/futures/v3/accounts/%s/leverage", symbol)
	var url = okApi + api

	var times FutureTimes
	Get(url, future.header("get", api, nil), &times)
	return times
}

/**
 * 设定合约币种杠杆倍数
 * direct (short|long)
 * time 倍数(10|20)
 */
func (future OKexFuture) SetTimes(symbol string, id string, direct string, time int32) FutureTimes {
	var api = fmt.Sprintf("/api/futures/v3/accounts/%s/leverage", symbol)
	var url = okApi + api

	param := make(map[string]interface{}, 0)
	param["leverage"] = time
	param["instrument_id"] = id
	param["direction"] = direct
	param["currency"] = symbol

	var times FutureTimes
	Post(url, future.header("post", api, param), param, &times)
	return times
}

// 账单流水查询
func (future OKexFuture) Ledger(symbol string) (FutureLedger, error) {
	var api = fmt.Sprintf("/api/futures/v3/accounts/%s/ledger", symbol)
	var url = okApi + api

	var leger FutureLedger
	err := Post(url, future.header("get", api, nil), nil, &leger)
	return leger, err
}

/**
 * type(1:开多2:开空3:平多4:平空)
 * order_type: 0：普通委托（order type不填或填0都是普通委托） 1：只做Maker（Post only） 2：全部成交或立即取消（FOK） 3：立即成交并取消剩余（IOC）
 * match_price: 是否以对手价下单(0:不是 1:是)，默认为0，当取值为1时。price字段无效，当以对手价下单，order_type只能选择0:普通委托
 */
func (future OKexFuture) Order(symbol string, _type int32,ordertype int32, price float32, size int32,match_price int32) (FutureOrder, error) {
	var api = "/api/futures/v3/order"
	var url = okApi + api

	param := make(map[string]interface{}, 0)
	param["instrument_id"] = symbol
	param["type"] = _type
	param["order_type"] = ordertype
	param["price"] = price
	param["size"] = size
	param["match_price"] = match_price

	var order FutureOrder
	err := Post(url, future.header("post", api, param), param, &order)
	return order, err
}

// 撤销指定订单
func (future OKexFuture) CancelOrder(symbol string, orderid string) (FutureCancel, error) {
	var api = fmt.Sprintf("/api/futures/v3/cancel_order/%s/%s", symbol, orderid)
	var url = okApi + api

	param := make(map[string]interface{}, 0)
	param["instrument_id"] = symbol
	param["order_id"] = orderid

	var cancel FutureCancel
	err := Post(url, future.header("post", api, param), param, &cancel)
	return cancel, err
}

// 获取订单列表
// status (订单状态(-1.撤单成功；0:等待成交 1:部分成交 2:全部成交 6：未完成（等待成交+部分成交）7：已完成（撤单成功+全部成交))
func (future OKexFuture) List(symbol string, status int) FutureList {
	var api = fmt.Sprintf("/api/futures/v3/orders/%s", symbol)
	var url = okApi + api

	var list FutureList
	Get(url, future.header("get", api, nil), &list)
	return list
}

// 获取订单信息
func (future OKexFuture) OrderInfo(symbol string, orderid string) FutureOrderInfo {
	var api = fmt.Sprintf("/api/futures/v3/orders/%s/%s", symbol, orderid)
	var url = okApi + api

	var orderInfo FutureOrderInfo
	Get(url, future.header("get", api, nil), &orderInfo)
	return orderInfo
}

// 获取成交明细
func (future OKexFuture) Fills(symbol string, orderid string) FutureFill {
	var api = fmt.Sprintf("/api/futures/v3/fills?instrument_id=%s&order_id=%s", symbol, orderid)
	var url = okApi + api

	var fills FutureFill
	Get(url, future.header("get", api, nil), &fills)
	return fills
}

// 获取合约信息
func (future OKexFuture) Instruments() []FutureInstrument {
	var api = "/api/futures/v3/instruments"
	var url = okApi + api

	var instruments []FutureInstrument
	Get(url, future.header("get", api, nil), &instruments)
	return instruments
}

// 获取深度数据
func (future OKexFuture) Depth(symbol string) FutureDepth {
	var api = fmt.Sprintf("/api/futures/v3/instruments/%s/book", symbol)
	var url = okApi + api

	var depth FutureDepth
	Get(url, future.header("get", api, nil), &depth)
	return depth
}

// 获取全部ticker信息
func (future OKexFuture) TickerAll() FutureTickers {
	var api = "/api/futures/v3/instruments/ticker"
	var url = okApi + api

	var tickers FutureTickers
	Get(url, future.header("get", api, nil), &tickers)
	return tickers
}

// 获取某个ticker信息
func (future OKexFuture) Ticker(symbol string) FutureTicker {
	var api = fmt.Sprintf("/api/futures/v3/instruments/%s/ticker", symbol)
	var url = okApi + api

	var ticker FutureTicker
	Get(url, future.header("get", api, nil), &ticker)
	return ticker
}

// 获取成交数据
func (future OKexFuture) Trades(symbol string, limit int32) ([]Trades, error) {
	var api = fmt.Sprintf("/api/futures/v3/instruments/%s/trades?limit=%d", symbol, limit)
	var url = okApi + api

	var trades []Trades
	err := Get(url, future.header("get", api, nil), &trades)
	return trades, err
}

// 获取K线数据
// 合约历史记录不能回溯,只能拉取最近200条(cnmd,K线数据可能不完整)
// 60   180  300  900   1800  3600  7200  14400 21600 43200  86400 604800
// 1min 3min 5min 15min 30min 1hour 2hour 4hour 6hour 12hour 1day  1week
func (future OKexFuture) Candle(instrumentId string, interval string, st time.Time) (FutureCandles, error) {
	var api string

	var gran int32
	if interval == "5m" {
		gran = 300
	} else if interval == "30m" {
		gran = 1800
	} else if interval == "1h" {
		gran = 3600
	} else if interval == "2h" {
		gran = 7200
	} else if interval == "4h" {
		gran = 14400
	} else if interval == "6h" {
		gran = 21600
	} else if interval == "12h" {
		gran = 43200
	} else if interval == "1d" {
		gran = 86400
	}

	if st.IsZero() {
		var err error
		st, err = time.Parse("2006-01-02 15:04:05", "2017-08-17 00:00:00")
		if err != nil {
			log.Print(err)
		}
	}
	start := IsoTime(st)
	api = fmt.Sprintf("/api/futures/v3/instruments/%s/candles?start=%s&granularity=%d", instrumentId, start, gran)

	var url = okApi + api

	var candles FutureCandles
	err := Get(url, future.header("get", api, nil), &candles)
	return candles, err
}

// 获取指数信息
func (future OKexFuture) Index(symbol string) FutureIndex {
	var api = fmt.Sprintf("/api/futures/v3/instruments/%s/index", symbol)
	var url = okApi + api
	// log.Print(url)

	var index FutureIndex
	Get(url, future.header("get", api, nil), &index)
	return index
}

// 获取法币汇率
func (future OKexFuture) Rate() FutureRate {
	var api = "/api/futures/v3/rate"
	var url = okApi + api
	// log.Print(url)

	var rate FutureRate
	Get(url, future.header("get", api, nil), &rate)
	return rate
}

// 获取预估交割价
func (future OKexFuture) EstimatedPrice(symbol string) FutureEstimatedPrice {
	var api = fmt.Sprintf("/api/futures/v3/instruments/%s/estimated_price", symbol)
	var url = okApi + api
	// log.Print(url)

	var price FutureEstimatedPrice
	Get(url, future.header("get", api, nil), &price)
	return price
}

// 获取平台总持仓量
func (future OKexFuture) OpenInterest(symbol string) FutureOpenInterest {
	var api = fmt.Sprintf("/api/futures/v3/instruments/%s/open_interest", symbol)
	var url = okApi + api
	// log.Print(url)

	var interest FutureOpenInterest
	Get(url, future.header("get", api, nil), &interest)
	return interest
}

// 获取当前限价
func (future OKexFuture) PriceLimit(symbol string) (FuturePriceLimit, error) {
	var api = fmt.Sprintf("/api/futures/v3/instruments/%s/price_limit", symbol)
	var url = okApi + api
	// log.Print(url)

	var limit FuturePriceLimit
	err := Get(url, future.header("get", api, nil), &limit)
	return limit, err
}

// 获取当前限价
func (future OKexFuture) MarkPrice(symbol string) FutureMarkPrice {
	var api = fmt.Sprintf("/api/futures/v3/instruments/%s/mark_price", symbol)
	var url = okApi + api
	// log.Print(url)

	var price FutureMarkPrice
	Get(url, future.header("get", api, nil), &price)
	return price
}

// 获取爆仓单
func (future OKexFuture) Liquidation(symbol string, status int) FutureLiquidation {
	var api = fmt.Sprintf("/api/futures/v3/instruments/%s/liquidation?status=%d", symbol, status)
	var url = okApi + api
	// log.Print(url)

	var liquidation FutureLiquidation
	Get(url, future.header("get", api, nil), &liquidation)
	return liquidation
}

// 获取合约挂单冻结数量
func (future OKexFuture) Holds(symbol string) FutureHold {
	var api = fmt.Sprintf("/api/futures/v3/accounts/%s/holds", symbol)
	var url = okApi + api
	// log.Print(url)

	var hold FutureHold
	Get(url, future.header("get", api, nil), &hold)
	return hold
}
