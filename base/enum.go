package base

type SYMBOL int          // 交易对
type PERIOD int          // k线时间
type DEPTH int           // 深度
type ORDER int           // 下单
type ORDER_TYPE int      // 订单类型
type CONTRACT_PERIOD int // 合约交割周期类型

const (
	BTCUSDT SYMBOL = iota
	ETHUSDT
	LTCUSDT
	EOSUSDT
	BTMUSDT
	BTCUSD
	ETHUSD
	LTCUSD
	EOSUSD

	MIN_1 PERIOD = iota
	MIN_3
	MIN_5
	MIN_15
	MIN_30
	HOUR_1
	HOUR_2
	HOUR_4
	HOUR_6
	HOUR_12
	DAY_1
	WEEK_1
	MONTH_1
	YEAR_1

	DEPTH_0 DEPTH = iota
	DEPTH_10
	DEPTH_100
	DEPTH_1000
	DEPTH_10000
	DEPTH_100000

	// 合约下单
	BUY_LONG   ORDER = iota // 开多
	BUY_SHORT               // 开空
	SELL_LONG               // 平多
	SELL_SHORT              // 平空

	MARKET_BUY        ORDER_TYPE = iota // 市价买单
	MARKET_SELL                         // 市价卖单
	LIMIT_BUY                           // 限价买单
	LIMIT_SELL                          // 限价卖单
	IOC_BUY                             // IOC买单(立刻执行或取消)
	IOC_SELL                            // IOC卖单
	MARKET_LIMIT_BUY                    // 市场限价买单
	MARKET_LIMIT_SELL                   // 市场限价卖单
	STOP_LIMIT_BUY
	STOP_LIMIT_SELL

	NONE     CONTRACT_PERIOD = iota // 现货|交割合约
	WEEK                            // 当周
	NEX_WEEK                        // 次周
	QUARTER                         // 季度
	SWAP                            // 永续
)

func Period(period PERIOD) string {
	var s string

	switch period {
	case MIN_1:
		s = "1m"
		break
	case MIN_5:
		s = "5m"
		break
	case MIN_30:
		s = "30m"
		break
	case HOUR_1:
		s = "1h"
		break
	case HOUR_2:
		s = "2h"
		break
	case HOUR_4:
		s = "4h"
		break
	case HOUR_6:
		s = "6h"
		break
	case HOUR_12:
		s = "12h"
		break
	case DAY_1:
		s = "1d"
		break
	}
	return s
}
