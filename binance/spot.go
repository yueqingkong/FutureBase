package binance

import (
	"fmt"
	"github.com/yueqingkong/FutureBase/base"
	"github.com/yueqingkong/FutureBase/util"
	"log"
	"sort"
	"strconv"
	"time"
)

type Spot struct {
	Url       string
	ApiKey    string
	SecretKey string
}

var (
	spot *Spot
)

func NewSpot() *Spot {
	if spot == nil {
		spot = new(Spot)
		spot.Url = "https://api.binance.com"
	}
	return spot
}

func (self *Spot) InitKeys(keys ...string) {
	self.ApiKey = keys[0]
	self.SecretKey = keys[1]
}

func (self *Spot) Symbol(symbol base.SYMBOL) string {
	var v string
	switch symbol {
	case base.BTCUSDT:
		v = "BTCUSDT"
		break
	case base.ETHUSDT:
		v = "ETHUSDT"
		break
	case base.LTCUSDT:
		v = "LTCUSDT"
		break
	case base.EOSUSDT:
		v = "EOSUSDT"
		break
	case base.BTMUSDT:
		v = "BTMUSDT"
		break
	}
	return v
}

func (self *Spot) Period(period base.PERIOD) string {
	var v string
	switch period {
	case base.MIN_1:
		v = "1m"
		break
	case base.MIN_5:
		v = "5m"
		break
	case base.MIN_15:
		v = "15m"
		break
	case base.MIN_30:
		v = "30m"
		break
	case base.HOUR_1:
		v = "1h"
		break
	case base.HOUR_4:
		v = "4h"
		break
	case base.DAY_1:
		v = "1d"
		break
	case base.WEEK_1:
		v = "1w"
		break
	case base.MONTH_1:
		v = "1m"
		break
	}
	return v
}

func (self *Spot) Depth(depth base.DEPTH) string {
	var v string
	switch depth {
	case base.DEPTH_0:
		v = "step0"
		break
	case base.DEPTH_10:
		v = "step1"
		break
	case base.DEPTH_100:
		v = "step2"
		break
	case base.DEPTH_1000:
		v = "step3"
		break
	case base.DEPTH_10000:
		v = "step4"
		break
	case base.DEPTH_100000:
		v = "step5"
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

// api验证时使用
func (self *Spot) header(param *map[string]string) map[string]string {
	keys := make([]string, 0)
	for k := range *param {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var index int
	var requestBody string
	for _, value := range keys {
		requestBody += fmt.Sprintf("%s=%s", value, (*param)[value])
		if index != len(*param)-1 {
			requestBody += "&"
		}
		index++
	}

	signer := util.HmacSha256Signer(requestBody, self.SecretKey)
	(*param)["signature"] = signer

	headerMap := make(map[string]string)
	headerMap["X-MBX-APIKEY"] = self.ApiKey
	return headerMap
}

// 测试能否联通
func (self *Spot) Ping() interface{} {
	var api = "/api/v1/ping"
	var url = self.Url + api
	var p interface{}
	util.Get(url, nil, &p)
	return p
}

// 获取服务器时间
func (self *Spot) Servertime() Time {
	var api = "/api/v1/time"
	var url = self.Url + api
	var time Time
	util.Get(url, nil, &time)
	return time
}

// 交易规范信息
func (self *Spot) Exchangeinfo() ExchangeInfo {
	var api = "/api/v1/exchangeInfo"
	var url = self.Url + api
	var info ExchangeInfo
	util.Get(url, nil, &info)
	return info
}

// 深度信息
func (self *Spot) Depths(symbol base.SYMBOL) Depth {
	var api = fmt.Sprintf("/api/v1/depth?symbol=%s", self.Symbol(symbol))
	var url = self.Url + api
	var depth Depth
	util.Get(url, nil, &depth)
	return depth
}

//  近期成交(归集)
func (self *Spot) Trades(symbol base.SYMBOL) []Trade {
	var api = fmt.Sprintf("/api/v1/trades?symbol=%s", self.Symbol(symbol))
	var url = self.Url + api
	var trade []Trade
	util.Get(url, nil, &trade)
	return trade
}

// K线数据
// 返回的数据是包含开始时间的那条K线的,返回默认默认条数，避免起始时间内没有数据返回空
func (self *Spot) KLines(symbol base.SYMBOL, period base.PERIOD, st time.Time) KLine {
	var start string
	//var end string

	if st.IsZero() {
		oldStart, err := time.Parse("2006-01-02 15:04:05", "2017-08-17 00:00:00")
		if err != nil {
			log.Print(err)
		}
		start = strconv.FormatInt(oldStart.UnixNano()/1e6, 10)
		st = oldStart
	} else {
		start = strconv.FormatInt(st.UnixNano()/1e6, 10)
	}

	var api = fmt.Sprintf("/api/v1/klines?symbol=%s&interval=%s&startTime=%s", self.Symbol(symbol), self.Period(period), start)
	var url = self.Url + api

	var kline KLine
	util.Get(url, nil, &kline)
	return kline
}

// 当前平均价格
func (self *Spot) AveragePrices(symbol base.SYMBOL) AveragePrice {
	var api = fmt.Sprintf("/api/v3/avgPrice?symbol=%s", self.Symbol(symbol))
	var url = self.Url + api
	var price AveragePrice
	util.Get(url, nil, &price)
	return price
}

// 24hr价格变动情况
func (self *Spot) DayHrs(symbol base.SYMBOL) DayHR {
	var api = fmt.Sprintf("/api/v1/ticker/24hr?symbol=%s", self.Symbol(symbol))
	var url = self.Url + api
	var hr DayHR
	util.Get(url, nil, &hr)
	return hr
}

// 获取交易对最新价格
func (self *Spot) Price(contract base.CONTRACT_TYPE, symbol base.SYMBOL) float32 {
	var api = fmt.Sprintf("/api/v3/ticker/price?symbol=%s", self.Symbol(symbol))
	var url = self.Url + api

	var p Price
	err := util.Get(url, nil, &p)

	var price float32
	if err == nil {
		price = util.StringToFloat32(p.Price)
	}
	return price
}

// 最优挂单接口
func (self *Spot) Tickers(symbol base.SYMBOL) Ticker {
	var api = fmt.Sprintf("/api/v3/ticker/bookTicker?symbol=%s", self.Symbol(symbol))
	var url = self.Url + api
	var ticker Ticker
	util.Get(url, nil, &ticker)
	return ticker
}

// 市场价下单
// side (SELL BUY)
// ty (LIMIT MARKET)
func (self *Spot) Order(symbol base.SYMBOL, side string, amount float32) Order {
	var api = "/api/v3/order"
	var url = self.Url + api

	params := make(map[string]string)
	params["symbol"] = self.Symbol(symbol)
	params["side"] = side
	params["type"] = "MARKET"
	params["quantity"] = util.Float32ToString(util.FloatDeceimal(amount))
	params["timestamp"] = util.UnixMillis(time.Now())

	headerMap := self.header(&params)

	var order Order
	util.PostForm(url, headerMap, params, &order)
	return order
}

// 账户信息
func (self *Spot) Account() Account {
	var api = "/api/v3/account"
	var url = self.Url + api

	params := make(map[string]string)
	params["timestamp"] = util.UnixMillis(time.Now())

	headerMap := self.header(&params)
	url = fmt.Sprintf("%s?timestamp=%s&signature=%s", url, params["timestamp"], params["signature"])

	var account Account
	util.Get(url, headerMap, &account)
	return account
}
