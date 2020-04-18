package plat

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"
)

var (
	okApi = "https://www.okex.com"

	// 获取Token 交易对
	symbolMap = map[string]string{
		"btc": "BTC-USDT",
		"eth": "ETH-USDT",
		"ltc": "LTC-USDT",
		"eos": "EOS-USDT",
	}

	// 申请的平台key
	ApiKey     = ""
	SecretKey  = ""
	Passphrase = ""
)

type OK struct {
}

func NewOKex() OK {
	return OK{}
}

/**
 * 初始化 Key
 */
func (ok OK) InitKeys(api string, secret string, phrase string) {
	ApiKey = api
	SecretKey = secret
	Passphrase = phrase
}

func (ok OK) header(request string, path string, body interface{}) map[string]string {
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

/**
 * 获取币对信息
 */
func (ok OK) Instruments() Instrument {
	var api = "/api/spot/v3/instruments"
	var url = okApi + api
	var inst Instrument
	Get(url, nil, &inst)
	return inst
}

/**
 * 获取深度数据
 */
func (ok OK) Depths(symbol string) Depth {
	var api = fmt.Sprintf("/api/spot/v3/instruments/%s/book", symbolMap[symbol])
	var url = okApi + api
	var depth Depth
	Get(url, nil, &depth)
	return depth
}

/**
 * 获取全部ticker信息
 */
func (ok OK) Tickers() []Ticker {
	var api = "/api/spot/v3/instruments/ticker"
	var url = okApi + api
	var tickers []Ticker
	Get(url, nil, &tickers)
	return tickers
}

/**
 * 获取某个ticker信息
 */
func (ok OK) TickerOne(id string) Ticker {
	var api = fmt.Sprintf("/api/spot/v3/instruments/%s/ticker", id)
	var url = okApi + api
	var ticker Ticker
	Get(url, nil, &ticker)
	return ticker
}

/**
 * 获取成交数据
 */
func (ok OK) Trades(symbol string) []Trade {
	var api = fmt.Sprintf("/api/spot/v3/instruments/%s/trades", symbolMap[symbol])
	var url = okApi + api
	var trades []Trade
	Get(url, nil, &trades)
	return trades
}

/**
 * 币币账户信息
 */
func (ok OK) Accounts() []Account {
	var api = "/api/spot/v3/accounts"
	var url = okApi + api
	var accounts []Account
	Get(url, ok.header("get", api, nil), &accounts)
	return accounts
}

/**
 * 账单流水查询
 */
func (ok OK) Ledger(symbol string) []Ledger {
	var api = fmt.Sprintf("/api/spot/v3/accounts/%s/ledger", symbol)
	var url = okApi + api
	var ledgers []Ledger
	Get(url, ok.header("get", api, nil), &ledgers)
	return ledgers
}

/**
 * 下单
 * type limit，market(默认是limit)，当以market（市价）下单，order_type只能选择0:普通委托
 * side(buy | sell)
 * order_type 0：普通委托（order type不填或填0都是普通委托） 1：只做Maker（Post only） 2：全部成交或立即取消（FOK） 3：立即成交并取消剩余（IOC）
 * buy:  买入金额
 * sell: 卖出数量
 */
func (ok OK) Order(symbol string, side string, buy float32, sell float32) Order {
	var api = "/api/spot/v3/orders"
	var url = okApi + api

	param := make(map[string]interface{}, 0)
	param["type"] = "market" //默认市价单
	param["side"] = side
	param["instrument_id"] = symbolMap[symbol]
	param["order_type"] = 0
	param["margin_trading"] = 1

	if param["type"] == "market" { //市价单参数
		param["notional"] = buy
		param["size"] = sell
	} else if param["type"] == "limit" { //限价单单参数

	}

	var order Order
	Post(url, ok.header("post", api, param), param, &order)
	return order
}

/**
 * 撤销指定订单
 */
func (ok OK) CancelOrder(symbol string, orderId string) CancelOrder {
	var api = fmt.Sprintf("/api/spot/v3/cancel_orders/%s", orderId)
	var url = okApi + api

	param := make(map[string]interface{}, 0)
	param["instrument_id"] = symbolMap[symbol]
	param["order_id"] = orderId

	var cancel CancelOrder
	Post(url, ok.header("post", api, param), param, &cancel)
	return cancel
}

/**
 * 获取订单列表
 * state 订单状态("-2":失败,"-1":撤单成功,"0":等待成交 ,"1":部分成交, "2":完全成交,"3":下单中,"4":撤单中,"6": 未完成（等待成交+部分成交），"7":已完成（撤单成功+全部成交））
 */
func (ok OK) OrderList(symbol string, state int32) []OrderItem {
	var api = fmt.Sprintf("/api/spot/v3/orders?instrument_id=%s&&state=%d", symbolMap[symbol], state)
	var url = okApi + api

	var item []OrderItem
	Get(url, ok.header("get", api, nil), &item)
	return item
}

/**
 * 获取所有未成交订单
 */
func (ok OK) OrderPending() []OrderPending {
	var api = "/api/spot/v3/orders_pending"
	var url = okApi + api

	var pending []OrderPending
	Get(url, ok.header("get", api, nil), &pending)
	return pending
}

/**
 * 获取订单信息
 */
func (ok OK) OrderInfo(symbol string, orderId string) OrderInfo {
	var api = fmt.Sprintf("/api/spot/v3/orders/%s?instrument_id=%s", orderId, symbolMap[symbol])
	var url = okApi + api

	var info OrderInfo
	Get(url, ok.header("get", api, nil), &info)
	return info
}

/**
 * 获取成交明细
 */
func (ok OK) Deals(symbol string, orderId string) Deal {
	var api = fmt.Sprintf("/api/spot/v3/fills?order_id=%s&instrument_id=%s", orderId, symbolMap[symbol])
	var url = okApi + api

	var deal Deal
	Get(url, ok.header("get", api, nil), &deal)
	return deal
}
