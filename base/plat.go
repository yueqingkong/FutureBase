package base

import "time"

//  底层平台通用 api
type PlatApi interface {
	Depth(depth DEPTH) string // 深度

	OrderType(order ORDER_TYPE) string // 订单类型
}

// (中间层)平台扩展接口 上层访问底层Api入口
type PlatBase interface {
	InitKeys([]string) // 账户key

	Symbol(SYMBOL) string // 交易对转普通类型(btc,eth,eos...)

	ORDER(ORDER) int32 // 下单类型(开多 1，开空 2，平多 3，平空 4)

	Price(period CONTRACT_PERIOD, symbol SYMBOL) float32 // 当前价格

	GetInstrument(period CONTRACT_PERIOD, symbol SYMBOL) string // 初始化 合约id

	Instrument(period CONTRACT_PERIOD, symbol SYMBOL) (string, time.Time) // 强制请求 合约id

	Delivery(period CONTRACT_PERIOD, symbol SYMBOL) (time.Time, time.Time) // 合约交割时间

	KLine(period CONTRACT_PERIOD, symbol SYMBOL, interval PERIOD, st time.Time) ([][]interface{}, error) // Kline

	Order(period CONTRACT_PERIOD, symbol SYMBOL, _type ORDER, price float32, unit int32) bool // 下单
}
