package trade

import (
	"fmt"
	"gitee.com/shieldpu_futures/FutureBase/orm"
	"gitee.com/shieldpu_futures/FutureBase/util"
	"github.com/yueqingkong/Okex/plat"
	"log"
	"strconv"
	"time"
)

var (
	fee     float32 = 0.002 // 手续费+交易滑点
	times   float32 = 10.0  // 合约倍数
	MaxLoss float32 = 400   // 最高点最大回撤
	MaxBuy  float32 = 0.8   // 最大层位
)

/**
 * 获取季度合约id
 */
func Instrument(symbol string) string {
	var instrumentid string

	future := plat.NewOKexFuture()
	instruments := future.Instruments()
	for _, value := range instruments {
		lower := util.Lower(value.UnderlyingIndex)
		id := value.InstrumentID
		alias := value.Alias
		delivery := util.StringYearToTime(value.Delivery)

		if value.Alias == "quarter" && value.IsInverse == "true" { // 季度
			// 同步 合约信息
			xorm := orm.NewXOrm()
			instrument := xorm.Instrument(lower)
			instrument.Key = id
			instrument.Delivery = delivery
			instrument.Period = alias
			if instrument.Symbol == "" {
				instrument.Symbol = lower
				xorm.InsertInstrument(instrument)
			} else {
				xorm.UpdateInstrument(instrument)
			}

			syncMap := orm.NewSyncMap()
			syncMap.SetInstrument(lower, id)

			if lower == symbol {
				instrumentid = id
			}
		}
	}

	return instrumentid
}

func Pulls(symbol string, srctions ...string) {
	for _, value := range srctions {
		PullHistory(symbol, value)
	}
}

/**
 * 同步 历史kline数据
 */
func PullHistory(symbol string, section string) {
	xorm := orm.XOrm{}
	lastCoin := xorm.Last(symbol, section)

	var startTime time.Time
	if lastCoin.Symbol == "" { // 记录为空
		startTime = time.Time{}
	} else {
		startTime = lastCoin.CreateTime
	}

	// 最后一条记录是昨天的
	diffHours := time.Now().Sub(startTime).Hours()
	// log.Print("[PullHistory] diffHours: ", diffHours)

	// 是否最新的数据
	if section == "30m" {
		if diffHours < 1 {
			return
		}
	} else if section == "1h" {
		if diffHours < 2 {
			return
		}
	} else if section == "2h" {
		if diffHours < 4 {
			return
		}
	} else if section == "4h" {
		if diffHours < 8 {
			return
		}
	} else if section == "6h" {
		if diffHours < 12 {
			return
		}
	} else if section == "12h" {
		if diffHours < 24 {
			return
		}
	} else if section == "1d" {
		if diffHours < 48 {
			return
		}
	}

	log.Print("[PullHistory] ", symbol, "  diffHours: ", diffHours, " section: ", section)
	api := plat.NewOKexFuture()

	// 避免重复返回最后一条的k线数据，加30s
	if !startTime.IsZero() {
		startTime = startTime.Add(time.Duration(30) * time.Second)
	}

	syncMap := orm.NewSyncMap()
	instrumentid := syncMap.GetInstrument(symbol)
	if instrumentid == "" {
		instrumentid = Instrument(symbol)
	}

	klines, err := api.Candle(instrumentid, section, startTime)
	if err != nil { // 重新获取 instrumentid
		instrumentid = Instrument(symbol)
		klines, _ = api.Candle(instrumentid, section, startTime)
	}

	coins := klineToCoin(symbol, section, klines)
	if len(coins) == 0 {
		log.Print("[PullHistory] len(coins) == 0 ")
	} else {
		xorm.InsertCoins(coins)
		PullHistory(symbol, section)
	}
}

/**
 * k线数据 -> orm.Coin
 * okex 的kline 是倒序的，最近的时间的在最前面
 */
func klineToCoin(symbol string, section string, kline plat.FutureCandles) []orm.Coin {
	var coins = make([]orm.Coin, 0)

	for k, value := range kline {
		var arr = value
		var open, _ = strconv.ParseFloat(arr[1].(string), 32)
		var close, _ = strconv.ParseFloat(arr[4].(string), 32)
		var high, _ = strconv.ParseFloat(arr[2].(string), 32)
		var low, _ = strconv.ParseFloat(arr[3].(string), 32)
		var volume, _ = strconv.ParseFloat(arr[5].(string), 32)
		var createtime, _ = util.IsoToTime(arr[0].(string))

		if k != 0 { // 最近时间一条有效的K线不保存
			var coin = orm.Coin{
				Symbol:     symbol,
				Plat:       "okex",
				Period:     section,
				Open:       float32(open),
				Close:      float32(close),
				High:       float32(high),
				Low:        float32(low),
				Volume:     float32(volume),
				Timestamp:  createtime.Unix(),
				CreateTime: createtime,
			}

			coins = append(coins, coin)
		}
	}
	return coins
}

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
		futures := plat.NewOKexFuture()
		account, err := futures.Account()
		log.Print(account)

		if err != nil {
			log.Println("[Account]", err)
		} else {
			xorm := orm.NewXOrm()
			local := xorm.Account(symbol)

			// 账户总额
			equity := util.StringToFloat32(total)
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

func Buy(symbol string, strategy string, op int32, price float32, size float32, canUnit float32, t time.Time) {
	syncMap := orm.NewSyncMap()
	mode := syncMap.Model()

	instrumentid := syncMap.GetInstrument(symbol)
	if mode == "test" {
		BuyRecord(symbol, strategy, op, price, size, canUnit, t)
	} else if mode == "release" {
		buySize := int32(size)

		if buySize == 0 {
			BuyRecord(symbol, strategy, op, price, size, canUnit, t)
		} else {
			var dif float32 = 50.0 // 滑点大些，在波动大的行情才能买进

			var orderPrice float32
			if op == 1 {
				orderPrice = price + dif
			} else if op == 2 {
				orderPrice = price - dif
			}

			futures := plat.NewOKexFuture()
			result, err := futures.Order(instrumentid, op, 2, orderPrice, buySize, 0)
			if err != nil {
				log.Println("[Buy] err: ", err)
			} else if result.Result { // 张数为0 result.Result=false 未能立马全部成交，返回的数据跟成交成功一样，不能区分
				// 未能立马全部成交，返回的数据跟成交成功一样，不能区分
				BuyRecord(symbol, strategy, op, price, size, canUnit, t)
			} else {
				log.Print("[Buy] result = false, ", result)
			}
		}
	}
}

/**
 * 开仓记录
 */
func BuyRecord(symbol string, strategy string, op int32, price float32, canSize float32, canUnit float32, t time.Time) {
	xorm := orm.NewXOrm()

	// 账户变更
	account := xorm.Account(symbol)
	buyOld := account.Buy
	balanceOld := account.Balance
	accountBefore := fmt.Sprintf("[account-before] Used: %f,Balance: %f, buyUnit: %f", buyOld, balanceOld, canUnit)

	account.Buy = buyOld + canUnit
	account.Balance = balanceOld - canUnit - PayFee(price, canUnit)

	total := account.Buy + account.Balance
	account.Total = total
	xorm.UpdateAccount(account)

	// 交易记录
	var position int32
	var avgPrce float32
	var totalUsed float32
	var totalSize float32

	records := xorm.LastRecord(symbol, strategy, 1)
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

	accountAfter := fmt.Sprintf("[account-after] Used: %f,Balance: %f", account.Buy, account.Balance)
	explain := accountBefore + Explain(t, strategy, op, position, price, avgPrce, totalUsed, totalSize, total, 0.0, 0.0, 0.0) + accountAfter

	xorm.InsertRecord(symbol, strategy, op, position, price, avgPrce, totalUsed, totalSize, total, 0, 0, explain, t)
}

/**
 * 平仓
 * op 1: 开多 2: 开空 3: 平多 4: 平空
 */
func Sell(symbol string, strategy string, op int32, price float32, t time.Time) {
	xorm := orm.NewXOrm()
	account := xorm.Account(symbol)
	buyOld := account.Buy
	balanceOld := account.Balance

	records := xorm.LastRecord(symbol, strategy, 1)

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

	instrumentid := syncMap.GetInstrument(symbol)
	if mode == "test" {
		SellRecord(symbol, strategy, op, price, t)
	} else if mode == "release" {
		if lastSize == 0 { // 开的空值单
			SellRecord(symbol, strategy, op, price, t)
		} else {
			var dif float32 = 50.0 // 滑点大些，在波动大的行情才能买进

			var orderPrice float32
			if op == 3 {
				orderPrice = price - dif
			} else if op == 4 {
				orderPrice = price + dif
			}

			futures := plat.NewOKexFuture()
			result, err := futures.Order(instrumentid, op, 2, orderPrice, int32(lastSize), 0)
			if err != nil {
				log.Print("[Sell] err: ", err)
			} else if result.Result {
				SellRecord(symbol, strategy, op, price, t)
			} else {
				log.Print("[Sell] result =false, ", result)
			}
		}
	}
}

/**
 * 平仓记录
 */
func SellRecord(symbol string, strategy string, op int32, price float32, t time.Time) {
	xorm := orm.NewXOrm()

	// 账户变更
	account := xorm.Account(symbol)
	accountBefore := fmt.Sprintf("[account-before] Used: %f,Balance: %f", account.Buy, account.Balance)

	records := xorm.LastRecord(symbol, strategy, 1)
	lastRecord := records[0]
	lastPosition := lastRecord.Position
	lastAvg := lastRecord.AvgPrice
	lastUsed := lastRecord.Used
	lastSize := lastRecord.Size

	// 收益
	profit := Profit(op, price, lastAvg, lastSize)
	profitRate := ProfitRate(profit, lastUsed)
	totlaRate := ProfitRate(profit, account.Total)
	payfee := PayFee(price, lastSize)

	account.Buy = account.Buy - lastUsed
	account.Balance = account.Balance + lastUsed + profit - payfee
	total := account.Balance + account.Buy
	account.Total = total
	xorm.UpdateAccount(account)

	// 交易记录
	accountAfter := fmt.Sprintf("[account-after] Used: %f,Balance: %f", account.Buy, account.Balance)
	explain := Explain(t, strategy, op, lastPosition, price, price, lastUsed, lastSize, total, profit, profitRate, totlaRate)
	explain = accountBefore + explain + accountAfter

	xorm.InsertRecord(symbol, strategy, op, lastPosition, price, lastAvg, lastUsed, lastSize, total, profit, profitRate, explain, t)
}

/**
 * token - > 张数
 */
func BuySize(price float32, buyunit float32) float32 {
	var size float32
	if buyunit == 0.0 {
		size = 0
	} else {
		amout := price * buyunit / 100.0
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
func PayFee(price float32, size float32) float32 {
	var value float32
	if price == 0 {
		value = 0.0
	} else {
		value = 100.0 / price * size * 0.0005
	}
	return value
}

/**
 * 收益
 * op 3 平多 4 平空
 */
func Profit(op int32, price float32, lastprice float32, size float32) float32 {
	var profit float32

	if lastprice == 0 || price == 0 || size == 0 {
		profit = 0
	} else if op == 3 {
		profit = (100.0/lastprice - 100.0/price) * size
	} else if op == 4 {
		profit = (100.0/price - 100.0/lastprice) * size
	}
	return profit
}

/**
 * 收益率
 * op 3 平多 4 平空
 */
func ProfitRate(profit float32, lastused float32) float32 {
	var value float32
	if lastused == 0 {
		value = 0.0
	} else {
		value = profit / lastused * 100.0
	}
	return value
}
