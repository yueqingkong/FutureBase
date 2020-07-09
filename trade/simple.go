package trade

import (
	"fmt"
	"github.com/yueqingkong/FutureBase/base"
	"github.com/yueqingkong/FutureBase/orm"
	"github.com/yueqingkong/FutureBase/util"
	"log"
	"time"
)

// 测试策略
// 买入 价格高于昨日价格 + R
// 卖出 价格高于昨日价格 - R
// 止损 价格回撤0.02
type Simple struct {
	NAME   string
	PERIOD string
	LOSE   float32
}

func NewSimple() Simple {
	return Simple{
		NAME:   "simple",
		PERIOD: "12h",
		LOSE:   0.02,
	}
}

func (self Simple) Name() string {
	return self.NAME
}

// 根据价格确认是否交易
// 交易类型 交易价格
func (self Simple) Tick(plat base.PlatBase, contract base.CONTRACT_PERIOD, symbol base.SYMBOL, price float32, t time.Time) {
	s := plat.Symbol(symbol)

	redisKey := self.NAME + "-" + s
	channelInfo := "simple"
	saveValue := util.Float32ToString(price)

	xorm := orm.NewXOrm()
	lastCoin := xorm.Last(s, self.PERIOD)
	if lastCoin.Symbol == "" {
		return
	}

	ATR := lastCoin.High - lastCoin.Low
	lastBuy := lastCoin.Close + ATR
	lastSell := lastCoin.Close - ATR

	records := xorm.LastRecord(s, self.NAME, 1)
	if len(records) == 0 || records[0].Operation == 3 || records[0].Operation == 4 { // 未持有仓位
		if util.PriceEqual(price, lastBuy) {
			KeepToRedis(redisKey, saveValue)

			self.Buy(plat, contract, symbol, base.BUY_LONG, price, channelInfo, t)
		} else if util.PriceEqual(price, lastSell) {
			KeepToRedis(redisKey, saveValue)

			self.Buy(plat, contract, symbol, base.BUY_SHORT, price, channelInfo, t)
		}
	} else {
		lastRecord := records[0]
		lastPrice := lastRecord.Price
		lastOp := lastRecord.Operation

		log.Println("[Last]", lastPrice, "[Op]", OpToString(lastOp))

		if lastOp == 1 { // 开多
			value := xorm.Redis(redisKey).V
			rangeMax := util.StringToFloat32(value)
			if rangeMax == 0 {
				rangeMax = price
			} else if rangeMax < price {
				rangeMax = price
				KeepToRedis(redisKey, saveValue)
			}

			losePrice := lastPrice - self.LOSE*price
			explain := fmt.Sprintf("[止损点] %f,[最大止损] %f", losePrice, lastSell)
			log.Println(explain)

			if price < losePrice || price < lastSell {
				Sell(plat, contract, symbol, self.Name(), base.SELL_LONG, price, t)
			}
		} else if lastOp == 2 { // 开空
			value := xorm.Redis(redisKey).V
			rangeMin := util.StringToFloat32(value)
			if rangeMin == 0 {
				rangeMin = price
			} else if rangeMin > price {
				rangeMin = price
				KeepToRedis(redisKey, saveValue)
			}

			losePrice := lastPrice + self.LOSE*price
			explain := fmt.Sprintf("[止损点] %f,[最大止损] %f", losePrice, lastSell)
			log.Println(explain)

			if price > losePrice || price > lastBuy { // 止损点
				Sell(plat, contract, symbol, self.Name(), base.SELL_SHORT, price, t)
			}
		}
	}
}

func (self Simple) Buy(plat base.PlatBase, contract base.CONTRACT_PERIOD, symbol base.SYMBOL,  op base.ORDER, price float32, explain string, t time.Time) {
	s := plat.Symbol(symbol)

	xorm := orm.NewXOrm()
	account := xorm.Account(s)
	buy := account.Buy
	balance := account.Balance
	total := account.Total

	records := xorm.LastRecord(s, self.NAME, 1)

	// 加仓单位
	canUnit := total

	size := BuySize(price, canUnit, ZDollar(s))
	canUnit = size * 100.0 / price                                                                                          // 重新计算开仓token
	if balance < (canUnit/Times) || buy+(canUnit/Times) >= total*MaxBuy || (len(records) > 0 && records[0].Position >= 5) { // 仓位已满
		log.Print("[买入] 仓位已满")
		BuyRecord(plat, symbol, self.NAME, op, price, 0, 0, explain, t)
	} else {
		if len(records) == 0 || records[0].Operation == 3 || records[0].Operation == 4 {
		} else {
			canUnit = 0.0
		}

		realUnit := canUnit / Times // 保证金
		size := BuySize(price, canUnit, ZDollar(s))
		canUnit = size * ZDollar(s) / price // 取整数，重新计算开仓token

		Buy(plat, contract, symbol, self.NAME, op, price, size, realUnit, explain, t)
	}
}
