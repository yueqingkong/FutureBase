package huobi

import (
	"fmt"
	"gitee.com/shieldpu_futures/FutureBase/base"
	"gitee.com/shieldpu_futures/FutureBase/util"
	"log"
	"net/url"
	"time"
)

type Spot struct {
	Url       string
	Accesskey string
	SecretKey string
}

var (
	spot *Spot
)

func NewSpot() *Spot {
	if spot == nil {
		spot = new(Spot)
		spot.Url = "https://api.huobi.pro"
	}
	return spot
}

// 初始化 Key
// api string, secret string, phrase string
func (self *Spot) InitKeys(keys ...string) {
	self.Accesskey = keys[0]
	self.SecretKey = keys[1]
}

func (self *Spot) Symbol(symbol base.SYMBOL) string {
	var v string
	switch symbol {
	case base.BTCUSDT:
		v = "btcusdt"
		break
	case base.ETHUSDT:
		v = "ethusdt"
		break
	case base.LTCUSDT:
		v = "ltcusdt"
		break
	case base.EOSUSDT:
		v = "eosusdt"
		break
	case base.BTMUSDT:
		v = "btmusdt"
		break
	}
	return v
}

func (self *Spot) Period(period base.PERIOD) string {
	var v string
	switch period {
	case base.MIN_1:
		v = "1min"
		break
	case base.MIN_5:
		v = "5min"
		break
	case base.MIN_15:
		v = "15min"
		break
	case base.MIN_30:
		v = "30min"
		break
	case base.HOUR_1:
		v = "60min"
		break
	case base.HOUR_4:
		v = "4hour"
		break
	case base.DAY_1:
		v = "1day"
		break
	case base.WEEK_1:
		v = "1week"
		break
	case base.MONTH_1:
		v = "1mon"
		break
	case base.YEAR_1:
		v = "1year"
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

//
func (self *Spot) Header() map[string]string {
	headers := make(map[string]string)
	headers["User-Agent"] = "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/39.0.2171.71 Safari/537.36"
	headers["Content-Type"] = "application/json"
	headers["Accept-Language"] = "zh-cn"
	return headers
}

// 参数签名
func (self *Spot) Signature(method string, request string) string {
	params := &url.Values{}

	timestamp := time.Now().UTC().Format("2006-01-02T15:04:05")
	params.Set(ACCESS_KEY, self.Accesskey)
	params.Set(SIGNNATURE_METHOD, "HmacSHA256")
	params.Set(SIGNNATURE_VERSION, "2")
	params.Set(TIMESTAMP, timestamp)

	encode := params.Encode()
	payload := fmt.Sprintf("%s\n%s\n%s\n%s", method, "api.huobi.pro", request, encode)
	sign, err := util.HmacSha256Base64Signer(payload, self.SecretKey)
	if err != nil {
		log.Println("[(self *Spot) header] ", payload, err)
	}
	params.Set(SIGNATURE, sign)

	signature := params.Encode()
	return signature
}

func (self *Spot) MethodToString(get bool) string {
	var method string
	if get {
		method = "GET"
	} else {
		method = "POST"
	}
	return method
}

// 获取所有交易对
func (self *Spot) Symbols() []SpotSymbol {
	var api = "/v1/common/symbols"
	var url = self.Url + api

	type Result struct {
		Data   []SpotSymbol `json:"data"`
		Status string       `json:"status"`
	}

	var result Result
	var symbols []SpotSymbol
	err := util.Get(url, nil, &result)

	if err == nil {
		symbols = result.Data
	}
	return symbols
}

// 获取所有币种
func (self *Spot) Currencys() []string {
	var api = "/v1/common/currencys"
	var url = self.Url + api

	type Result struct {
		Data   []string `json:"data"`
		Status string   `json:"status"`
	}

	var result Result
	var currency []string
	err := util.Get(url, nil, &result)

	if err == nil {
		currency = result.Data
	}
	return currency
}

// 查询各币种及其所在区块链的静态参考信息（公共数据）
func (self *Spot) Currencies(currency string) []Currency {
	var api = "/v2/reference/currencies"
	var url = self.Url + api

	type Result struct {
		Data []Currency `json:"data"`
		Code int        `json:"status"`
	}

	var result Result
	var currencies []Currency
	err := util.Get(url, nil, &result)

	if err == nil && result.Code == 2000 {
		currencies = result.Data
	}
	return currencies
}

// 返回当前的系统时间，时间是调整为北京时间的时间戳，单位毫秒
func (self *Spot) Timestamp() int32 {
	var api = "/v1/common/timestamp"
	var url = self.Url + api

	type Result struct {
		Status string `json:"status"`
		Data   int32  `json:"data"`
	}

	var result Result
	var timestamp int32
	err := util.Get(url, nil, &result)

	if err == nil && result.Status == "ok" {
		timestamp = result.Data
	}
	return timestamp
}

// 历史K线数据
func (self *Spot) Kline(symbol base.SYMBOL, period base.PERIOD) []KLine {
	var api = fmt.Sprintf("/market/history/kline?symbol=%s&period=%s",
		self.Symbol(symbol), self.Period(period))

	var url = self.Url + api
	log.Println(url)

	type Result struct {
		Status string  `json:"status"`
		Data   []KLine `json:"data"`
		Ts     int     `json:"ts"`
		Ch     string  `json:"ch"`
	}

	var result Result
	var klines []KLine
	err := util.Get(url, nil, &result)

	if err == nil && result.Status == "ok" {
		klines = result.Data
	}
	return klines
}

// 此接口获取ticker信息同时提供最近24小时的交易聚合信息
func (self *Spot) Merged(symbol base.SYMBOL) SpotMergin {
	var api = fmt.Sprintf("/market/detail/merged?symbol=%s", self.Symbol(symbol))
	var url = self.Url + api
	log.Println(url)

	type Result struct {
		Status string     `json:"status"`
		Tick   SpotMergin `json:"tick"`
		Ch     string     `json:"ch"`
		Ts     int        `json:"ts"`
	}

	var result Result
	var mergin SpotMergin
	err := util.Get(url, nil, &result)

	if err == nil && result.Status == "ok" {
		mergin = result.Tick
	}
	return mergin
}

// 获得所有交易对的 tickers，数据取值时间区间为24小时滚动。
func (self *Spot) Tick() []SpotTick {
	var api = "/market/tickers"
	var url = self.Url + api
	log.Println(url)

	type Result struct {
		Status string     `json:"status"`
		Data   []SpotTick `json:"data"`
		Ch     string     `json:"ch"`
		Ts     int        `json:"ts"`
	}

	var result Result
	var tick []SpotTick
	err := util.Get(url, nil, &result)

	if err == nil && result.Status == "ok" {
		tick = result.Data
	}
	return tick
}

// 此接口返回指定交易对的当前市场深度数据。
func (self *Spot) Depths(symbol base.SYMBOL, depth int32, ty base.DEPTH) SpotDepth {
	var api = fmt.Sprintf("/market/depth?symbol=%s&depth=%d&type=%s", self.Symbol(symbol), depth, self.Depth(ty))
	var url = self.Url + api
	log.Println(url)

	type Result struct {
		Status string    `json:"status"`
		Data   SpotDepth `json:"tick"`
		Ch     string    `json:"ch"`
		Ts     int       `json:"ts"`
	}

	var result Result
	var back SpotDepth
	err := util.Get(url, nil, &result)

	if err == nil && result.Status == "ok" {
		back = result.Data
	}
	return back
}

// 此接口返回指定交易对最新的一个交易记录。
func (self *Spot) MarketTrade(symbol base.SYMBOL) SpotTrade {
	var api = fmt.Sprintf("/market/trade?symbol=%s", self.Symbol(symbol))
	var url = self.Url + api
	log.Println(url)

	type Result struct {
		Status string    `json:"status"`
		Data   SpotTrade `json:"tick"`
		Ch     string    `json:"ch"`
		Ts     int       `json:"ts"`
	}

	var result Result
	var data SpotTrade
	err := util.Get(url, nil, &result)

	if err == nil && result.Status == "ok" {
		data = result.Data
	}
	return data
}

// 此接口返回指定交易对最新的一个交易记录。
func (self *Spot) Price(contract base.CONTRACT_TYPE, symbol base.SYMBOL) float32 {
	var api = fmt.Sprintf("/market/trade?symbol=%s", self.Symbol(symbol))
	var url = self.Url + api
	// log.Println(url)

	type Result struct {
		Status string    `json:"status"`
		Data   SpotTrade `json:"tick"`
		Ch     string    `json:"ch"`
		Ts     int       `json:"ts"`
	}

	var result Result
	var data SpotTrade
	err := util.Get(url, nil, &result)

	var price float32
	if err == nil && result.Status == "ok" {
		data = result.Data
		price = data.Record[0].Price
	}
	return price
}

// 此接口返回指定交易对近期的所有交易记录。
func (self *Spot) HistoryTrade(symbol base.SYMBOL, size int32) []SpotTrade {
	var api = fmt.Sprintf("/market/history/trade?symbol=%s&size=%d", self.Symbol(symbol), size)
	var url = self.Url + api
	log.Println(url)

	type Result struct {
		Status string      `json:"status"`
		Data   []SpotTrade `json:"data"`
		Ch     string      `json:"ch"`
		Ts     int         `json:"ts"`
	}

	var result Result
	var data []SpotTrade
	err := util.Get(url, nil, &result)

	if err == nil && result.Status == "ok" {
		data = result.Data
	}
	return data
}

// 此接口返回最近24小时的行情数据汇总。
func (self *Spot) MarketDetail(symbol base.SYMBOL) SpotDetail {
	var api = fmt.Sprintf("/market/detail?symbol=%s", self.Symbol(symbol))
	var url = self.Url + api
	log.Println(url)

	type Result struct {
		Status string     `json:"status"`
		Data   SpotDetail `json:"tick"`
		Ch     string     `json:"ch"`
		Ts     int        `json:"ts"`
	}

	var result Result
	var data SpotDetail
	err := util.Get(url, nil, &result)

	if err == nil && result.Status == "ok" {
		data = result.Data
	}
	return data
}

// 查询当前用户的所有账户 ID account-id 及其相关信息
func (self *Spot) Accounts() []SpotAccount {
	var api = "/v1/account/accounts"

	signature := self.Signature(self.MethodToString(true), api)
	var url = self.Url + api + "?" + signature

	type Result struct {
		Status string        `json:"status"`
		Data   []SpotAccount `json:"data"`
		Ch     string        `json:"ch"`
		Ts     int           `json:"ts"`
	}

	var result Result
	var data []SpotAccount
	err := util.Get(url, nil, &result)

	if err == nil && result.Status == "ok" {
		data = result.Data
	}
	return data
}

// 查询指定账户的余额，支持以下账户：
// spot：现货账户， margin：逐仓杠杆账户，otc：OTC 账户，point：点卡账户，super-margin：全仓杠杆账户
func (self *Spot) Balance(accountid string) SpotBalance {
	var api = fmt.Sprintf("/v1/account/accounts/%s/balance", accountid)

	signature := self.Signature(self.MethodToString(true), api)
	var url = self.Url + api + "?" + signature

	type Result struct {
		Status string      `json:"status"`
		Data   SpotBalance `json:"data"`
		Ch     string      `json:"ch"`
		Ts     int         `json:"ts"`
	}

	var result Result
	var data SpotBalance
	err := util.Get(url, nil, &result)

	if err == nil && result.Status == "ok" {
		data = result.Data
	}
	return data
}

// 该节点基于用户账户ID返回账户流水。
func (self *Spot) AccountHistory(accountid string) SpotBalance {
	var api = "/v1/account/history"

	signature := self.Signature(self.MethodToString(true), api)
	var url = self.Url + api + "?account_id" + accountid + "&" + signature

	type Result struct {
		Status string      `json:"status"`
		Data   SpotBalance `json:"data"`
		Ch     string      `json:"ch"`
		Ts     int         `json:"ts"`
	}

	var result Result
	var data SpotBalance
	err := util.Get(url, nil, &result)

	if err == nil && result.Status == "ok" {
		data = result.Data
	}
	return data
}

// 母账户查询其下所有子账号的各币种汇总余额
func (self *Spot) AggregateBalance() []SpotAggregateBalance {
	var api = "/v1/subuser/aggregate-balance"

	signature := self.Signature(self.MethodToString(true), api)
	var url = self.Url + api + "?" + signature

	type Result struct {
		Status string                 `json:"status"`
		Data   []SpotAggregateBalance `json:"data"`
		Ch     string                 `json:"ch"`
		Ts     int                    `json:"ts"`
	}

	var result Result
	var data []SpotAggregateBalance
	err := util.Get(url, nil, &result)

	if err == nil && result.Status == "ok" {
		data = result.Data
	}
	return data
}

// 母账户查询其下所有子账号的各币种汇总余额
// type buy-market：市价买, sell-market：市价卖, buy-limit：限价买, sell-limit：限价卖, buy-ioc：IOC买单, sell-ioc：IOC卖单, buy-limit-maker, sell-limit-maker
func (self *Spot) Order(accountid string, symbol base.SYMBOL, ty string, amount float32) string {
	var api = "/v1/order/orders/place"

	signature := self.Signature(self.MethodToString(false), api)
	var url = self.Url + api + "?" + signature

	paramsMap := make(map[string]interface{})
	paramsMap["account-id"] = accountid
	paramsMap["symbol"] = self.Symbol(symbol)
	paramsMap["type"] = ty
	paramsMap["amount"] = util.Float32ToString(amount)

	type Result struct {
		Status string `json:"status"`
		Data   string `json:"data"`
		Ch     string `json:"ch"`
		Ts     int    `json:"ts"`
	}

	var result Result
	var data string
	err := util.Post(url, self.Header(), paramsMap, &result)

	if err == nil && result.Status == "ok" {
		data = result.Data
	}
	return data
}
