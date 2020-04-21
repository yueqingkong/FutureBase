package trade

import (
	"fmt"
	"gitee.com/shieldpu_futures/FutureBase/base"
	"gitee.com/shieldpu_futures/FutureBase/orm"
	"gitee.com/shieldpu_futures/FutureBase/util"
	"log"
	"time"
)

var (
	fee     float32 = 0.002 // 手续费+交易滑点
	Times   float32 = 10.0  // 合约倍数
	MaxLoss float32 = 400   // 最高点最大回撤
	MaxBuy  float32 = 0.8   // 最大层位
)

/**
 * 同步账号
 */
func Account(symbol string, total string) {
	syncMap := orm.NewSyncMap()
	mode := syncMap.Model()

	if mode == "test" {
		log.Print(fmt.Sprintf("[测试账号] %s", symbol))
		xorm := orm.NewXOrm()

		account := xorm.Account(symbol)
		if account.Symbol == "" || account.Total == 0 {
			account.Symbol = symbol
			account.Balance = 1.0
			account.Buy = 0
			account.Total = 1.0
			xorm.InsertAccount(account)
		}
	} else if mode == "release" {
		xorm := orm.NewXOrm()
		local := xorm.Account(symbol)

		// 账户总额
		equity := util.StringToFloat32(total)
		if equity < 0.0 { // 金额不为0 ,如果小于0 金额不变
			return
		}

		local.Total = equity
		if local.Symbol == "" {
			local.Symbol = symbol
			local.Buy = 0.0
			local.Balance = equity
			xorm.InsertAccount(local)
		} else {
			local.Balance = equity - local.Buy
			xorm.UpdateAccount(local)
		}
	}
}

/**
 * (通道突破) 周期内最高价/最低价
 */
func Channel(symbol string, section string, backRange int32, t time.Time) (float32, float32) {
	xorm := orm.NewXOrm()

	var high float32
	var low float32

	coins := xorm.Before(symbol, section, backRange)
	for k, value := range coins {
		if k == int(backRange) {

		} else if k == 0 {
			high = value.High
			low = value.Low
		} else {
			if high < value.High {
				high = value.High
			}
			if low > value.Low {
				low = value.Low
			}
		}
	}
	return high, low
}

/**
 * 平均波动幅度
 *  1、当前交易日的最高价与最低价间的波幅
 *  2、前一交易日收盘价与当个交易日最高价间的波幅
 *  3、前一交易日收盘价与当个交易日最低价间的波幅
 */
func ATR(symbol string, section string, limit int32) float32 {
	xorm := orm.NewXOrm()
	coins := xorm.Before(symbol, section, limit)

	var totalRange float32
	for k, value := range coins {
		var dayRange float32
		if k == 0 {
			dayRange = value.High - value.Low
		} else {
			last := coins[k-1]
			lastHigh := util.Abs(last.Close - value.High)
			lastLow := util.Abs(last.Close - value.Low)
			todayRange := value.High - value.Low
			dayRange = util.Max(todayRange, lastHigh, lastLow)
		}
		totalRange += dayRange
	}

	length := len(coins)
	atr := totalRange / float32(length)
	return atr
}

/**
 * 买入单位,回撤 ATR 亏损比例
 * rate*Total = atr/price * Unit
 * rate: 总账户最大亏损(总账户 满足连续最大亏损20次)
 * atr/price 回撤(高位最大回撤 0.04)
 */
func Unit(price float32, total float32) float32 {
	unit := (total * price) / (20.0 * MaxLoss)
	return unit
}

/**
 * 操作 -->显示
 */
func OpToString(op int32) string {
	var ty string
	if op == 1 {
		ty = "[开多]"
	} else if op == 2 {
		ty = "[开空]"
	} else if op == 3 {
		ty = "[平多]"
	} else if op == 4 {
		ty = "[平空]"
	}
	return ty
}

/**
 * 简介 (交易详情)
 * 1: 开多 2: 开空 3: 平仓
 */
func Explain(t time.Time, strategy string, op int32, position int32, price float32, averageprice float32, used float32, size float32, total, profit float32, profitrate float32, totalrate float32) string {
	ty := OpToString(op)

	var explain string
	if op == 1 || op == 2 {
		explain = fmt.Sprintf("[时间]: %s,[%s], %s ,[仓位] %d,[当前价] %f,[均价] %f,[开仓] %f,[张数] %f,[总额] %f",
			util.TimeToString(t), strategy, ty, position, price, averageprice, used, size, total)
	} else if op == 3 || op == 4 {
		explain = fmt.Sprintf("[时间]: %s,[%s], %s ,[仓位] %d,[当前价] %f,[平仓] %f,[张数] %f,[总额] %f, [收益] %f, [收益率] %f,[账户收益率] %f",
			util.TimeToString(t), strategy, ty, position, price, used, size, total, profit, profitrate, totalrate)
	}

	return explain
}

/**
 * 保存到Redis
 */
func KeepToRedis(key string, value string) {
	xorm := orm.NewXOrm()
	redis := xorm.Redis(key)
	if redis.K == "" {
		redis.K = key
		redis.V = value
		xorm.InsertRedis(redis)
	} else {
		redis.V = value
		xorm.UpdateRedis(redis)
	}
}

func Buy(plat base.PlatBase, cnotract base.CONTRACT_PERIOD, symbol base.SYMBOL, strategy string, operation base.ORDER, price float32, size float32, canUnit float32, show string, t time.Time) {
	syncMap := orm.NewSyncMap()
	mode := syncMap.Model()

	if mode == "test" {
		BuyRecord(plat, symbol, strategy, operation, price, size, canUnit, show, t)
	} else if mode == "release" {
		buySize := int32(size)

		if buySize == 0 {
			BuyRecord(plat, symbol, strategy, operation, price, size, canUnit, show, t)
		} else {
			var dif float32 = 50.0 // 滑点大些，在波动大的行情才能买进

			var orderPrice float32
			if operation == base.BUY_LONG {
				orderPrice = price + dif
			} else if operation == base.BUY_SHORT {
				orderPrice = price - dif
			}

			success := plat.Order(cnotract, symbol, operation, orderPrice, int32(size))
			if success { // 未能立马全部成交，返回的数据跟成交成功一样，不能区分
				BuyRecord(plat, symbol, strategy, operation, price, size, canUnit, show, t)
			}
		}
	}
}

/**
 * 开仓记录
 */
func BuyRecord(plat base.PlatBase, symbol base.SYMBOL, strategy string, operation base.ORDER, price float32, canSize float32, canUnit float32, show string, t time.Time) {
	s := plat.Symbol(symbol)

	xorm := orm.NewXOrm()

	// 账户变更
	account := xorm.Account(s)
	buyOld := account.Buy
	balanceOld := account.Balance
	accountBefore := fmt.Sprintf("[account-before] Used: %f,Balance: %f, buyUnit: %f", buyOld, balanceOld, canUnit)

	account.Buy = buyOld + canUnit
	account.Balance = balanceOld - canUnit - PayFee(price, canUnit, ZDollar(s))

	total := account.Buy + account.Balance
	account.Total = total
	xorm.UpdateAccount(account)

	// 交易记录
	var position int32
	var avgPrce float32
	var totalUsed float32
	var totalSize float32

	records := xorm.LastRecord(s, strategy, 1)
	if len(records) == 0 || records[0].Operation == 3 || records[0].Operation == 4 {
		position = int32(1)
		avgPrce = price
		totalUsed = canUnit
		totalSize = canSize
	} else {
		lastRecord := records[0]
		lastPosition := lastRecord.Position
		lastAvg := lastRecord.AvgPrice
		lastUsed := lastRecord.Used
		lastSize := lastRecord.Size

		position = lastPosition + 1

		// 余额太小 第二次开单为0
		if lastUsed+canUnit == 0 {
			avgPrce = price
		} else {
			avgPrce = (lastUsed*lastAvg + canUnit*price) / (lastUsed + canUnit)
		}

		totalUsed = lastUsed + canUnit
		totalSize = lastSize + canSize
	}

	op := plat.ORDER(operation)

	accountAfter := fmt.Sprintf("[account-after] Used: %f,Balance: %f", account.Buy, account.Balance)
	explain := accountBefore + Explain(t, strategy, op, position, price, avgPrce, totalUsed, totalSize, total, 0.0, 0.0, 0.0) + show + accountAfter

	xorm.InsertRecord(s, strategy, op, position, price, avgPrce, totalUsed, totalSize, total, 0, 0, explain, t)
}

/**
 * 平仓
 * op 1: 开多 2: 开空 3: 平多 4: 平空
 */
func Sell(plat base.PlatBase, contract base.CONTRACT_PERIOD, symbol base.SYMBOL, strategy string, operation base.ORDER, price float32, t time.Time) {
	s := plat.Symbol(symbol)

	xorm := orm.NewXOrm()
	account := xorm.Account(s)
	buyOld := account.Buy
	balanceOld := account.Balance

	records := xorm.LastRecord(s, strategy, 1)

	accountBefore := fmt.Sprintf("[account-before] Used: %f,Balance: %f", buyOld, balanceOld)
	log.Print(accountBefore)
	if len(records) == 0 || records[0].Operation == 3 || records[0].Operation == 4 {
		log.Print("[卖出] 仓位已清空")
		return
	}

	lastRecord := records[0]
	lastSize := lastRecord.Size

	// 下单
	syncMap := orm.NewSyncMap()
	mode := syncMap.Model()

	if mode == "test" {
		SellRecord(plat, symbol, strategy, operation, price, t)
	} else if mode == "release" {
		if lastSize == 0 { // 开的空值单
			SellRecord(plat, symbol, strategy, operation, price, t)
		} else {
			var dif float32 = 50.0 // 滑点大些，在波动大的行情才能买进

			var orderPrice float32
			if operation == base.SELL_LONG {
				orderPrice = price - dif
			} else if operation == base.SELL_SHORT {
				orderPrice = price + dif
			}

			success := plat.Order(contract, symbol, operation, orderPrice, int32(lastSize))
			if success {
				SellRecord(plat, symbol, strategy, operation, price, t)
			}
		}
	}
}

/**
 * 平仓记录
 */
func SellRecord(plat base.PlatBase, symbol base.SYMBOL, strategy string, operation base.ORDER, price float32, t time.Time) {
	s := plat.Symbol(symbol)
	o := plat.ORDER(operation)

	xorm := orm.NewXOrm()

	// 账户变更
	account := xorm.Account(s)
	accountBefore := fmt.Sprintf("[account-before] Used: %f,Balance: %f", account.Buy, account.Balance)

	records := xorm.LastRecord(s, strategy, 1)
	lastRecord := records[0]
	lastPosition := lastRecord.Position
	lastAvg := lastRecord.AvgPrice
	lastUsed := lastRecord.Used
	lastSize := lastRecord.Size

	// 收益
	profit := Profit(o, price, lastAvg, lastSize, ZDollar(s))
	profitRate := ProfitRate(profit, lastUsed, ZDollar(s))
	totlaRate := ProfitRate(profit, account.Total, ZDollar(s))
	payfee := PayFee(price, lastSize, ZDollar(s))

	account.Buy = account.Buy - lastUsed
	account.Balance = account.Balance + lastUsed + profit - payfee
	total := account.Balance + account.Buy
	account.Total = total
	xorm.UpdateAccount(account)

	// 交易记录
	accountAfter := fmt.Sprintf("[account-after] Used: %f,Balance: %f", account.Buy, account.Balance)
	explain := Explain(t, strategy, o, lastPosition, price, price, lastUsed, lastSize, total, profit, profitRate, totlaRate)
	explain = accountBefore + explain + accountAfter

	xorm.InsertRecord(s, strategy, o, lastPosition, price, lastAvg, lastUsed, lastSize, total, profit, profitRate, explain, t)
}

/**
 * token - > 张数
 */
func BuySize(price float32, buyunit float32, zhang float32) float32 {
	var size float32
	if buyunit == 0.0 {
		size = 0
	} else {
		amout := price * buyunit / zhang
		if amout < 1.0 {
			size = 1.0
		} else {
			size = util.Floor(amout)
		}
	}
	return size
}

/**
 * 支付手续费
 */
func PayFee(price float32, size float32, zhang float32) float32 {
	var value float32
	if price == 0 {
		value = 0.0
	} else {
		value = zhang / price * size * 0.0005
	}
	return value
}

/**
 * 收益
 * op 3 平多 4 平空
 */
func Profit(op int32, price float32, lastprice float32, size float32, zhang float32) float32 {
	var profit float32

	if lastprice == 0 || price == 0 || size == 0 {
		profit = 0
	} else if op == 3 {
		profit = (zhang/lastprice - zhang/price) * size
	} else if op == 4 {
		profit = (zhang/price - zhang/lastprice) * size
	}
	return profit
}

/**
 * 收益率
 * op 3 平多 4 平空
 */
func ProfitRate(profit float32, lastused float32, zhang float32) float32 {
	var value float32
	if lastused == 0 {
		value = 0.0
	} else {
		value = profit / lastused * zhang
	}
	return value
}

// 一张 代表的面纸
func ZDollar(symbol string) float32 {
	var v float32

	if symbol == "btc" {
		v = 100.0
	} else {
		v = 10.0
	}
	return v
}
