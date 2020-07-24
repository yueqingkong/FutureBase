package okex

import (
	"github.com/yueqingkong/FutureBase/base"
	"github.com/yueqingkong/FutureBase/orm"
	"github.com/yueqingkong/FutureBase/util"
	"log"
	"strconv"
	"time"
)

// (中间层) Api的抽象访问
type OkexFuture struct {
	*Future
}

var (
	okexFuture *OkexFuture
)

func NewOkexFuture() *OkexFuture {
	if okexFuture == nil {
		okexFuture = new(OkexFuture)
		okexFuture.Future = NewFuture()
	}
	return okexFuture
}

func (self *OkexFuture) InitKeys(keys []string) {
	self.Apikey = keys[0]
	self.SecretKey = keys[1]
	self.PhraseKey = keys[2]
}

func (self *OkexFuture) Symbol(symbol base.SYMBOL) string {
	var v string
	switch symbol {
	case base.BTCUSD, base.BTCUSDT:
		v = "btc"
		break
	case base.ETHUSD, base.ETHUSDT:
		v = "eth"
		break
	case base.LTCUSD, base.LTCUSDT:
		v = "ltc"
		break
	case base.EOSUSD, base.EOSUSDT:
		v = "eos"
		break
	}
	return v
}

func (self *OkexFuture) ORDER(order base.ORDER) int32 {
	var v int32
	switch order {
	case base.BUY_LONG:
		v = 1
		break
	case base.BUY_SHORT:
		v = 2
		break
	case base.SELL_LONG:
		v = 3
		break
	case base.SELL_SHORT:
		v = 4
		break
	}
	return v
}

func (self *OkexFuture) GetInstrument(contract base.CONTRACT_PERIOD, symbol base.SYMBOL) string {
	s := self.Symbol(symbol)

	syncMap := orm.NewSyncMap()
	instrumentid := syncMap.GetInstrument(s) // 合约id
	if instrumentid == "" {
		xorm := orm.NewXOrm()
		instrument := xorm.Instrument(s)

		if instrument.Symbol == "" {
			instrumentid, _ = self.Instrument(contract, symbol)
		} else {
			instrumentid = instrument.Key
			syncMap.SetInstrument(s, instrumentid)
		}
	}
	return instrumentid
}

// Okex合约交割时间(只有日期没有时间)
// 合约时间范围
func (self *OkexFuture) Delivery(contract base.CONTRACT_PERIOD, symbol base.SYMBOL) (time.Time, time.Time) {
	s := self.Symbol(symbol)

	var t time.Time

	xorm := orm.NewXOrm()
	instrument := xorm.Instrument(s)

	if instrument.Symbol == "" {
		_, t = self.Instrument(contract, symbol)
	} else {
		t = instrument.Delivery
	}

	begin := t.Add(time.Duration(-30) * time.Minute) // 30min 前
	end := t.Add(time.Duration(30) * time.Minute)    // 30min 后
	return begin, end
}

// 合约信息
// 合约id, 合约交割日期
// 合约交割时间(只有日期没有时间)
// 一般是周五 下午4点,上下偏差 30分钟
func (self *OkexFuture) Instrument(conrtact base.CONTRACT_PERIOD, symbol base.SYMBOL) (string, time.Time) {
	var instrumentid string
	var deliveryTime time.Time

	instruments := self.Instruments()
	for _, value := range instruments {
		lower := util.Lower(value.UnderlyingIndex)
		id := value.InstrumentID
		alias := value.Alias
		delivery := util.StringToTime(value.Delivery + " 16:00:00") // 下午四点交割

		if (conrtact == base.QUARTER && value.Alias == "quarter" && value.IsInverse == "true") || // 季度
			(conrtact == base.NEX_WEEK && value.Alias == "next_week" && value.IsInverse == "true") || // 次周
			(conrtact == base.WEEK && value.Alias == "this_week" && value.IsInverse == "true") { // 当周

			syncMap := orm.NewSyncMap()
			syncMap.SetInstrument(lower, id)

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

			if lower == util.Lower(self.Symbol(symbol)) {
				instrumentid = id
				deliveryTime = delivery
			}
		}
	}

	return instrumentid, deliveryTime
}

func (self *OkexFuture) Price(contract base.CONTRACT_PERIOD, symbol base.SYMBOL) float32 {
	instrumentid := self.GetInstrument(contract, symbol) // 合约id

	return self.Future.Ticker(instrumentid)
}

func (self *OkexFuture) KLine(contract base.CONTRACT_PERIOD, symbol base.SYMBOL, interval base.PERIOD, st time.Time) ([]orm.Coin, error) {
	instrumentid := self.GetInstrument(contract, symbol) // 合约id

	var gran int32
	if interval == base.MIN_5 {
		gran = 300
	} else if interval == base.MIN_30 {
		gran = 1800
	} else if interval == base.HOUR_1 {
		gran = 3600
	} else if interval == base.HOUR_2 {
		gran = 7200
	} else if interval == base.HOUR_4 {
		gran = 14400
	} else if interval == base.HOUR_6 {
		gran = 21600
	} else if interval == base.HOUR_12 {
		gran = 43200
	} else if interval == base.DAY_1 {
		gran = 86400
	}

	var coins []orm.Coin
	klines, err := self.Future.Candle(instrumentid, gran, st)
	if err != nil { // 重新获取 instrumentid
		self.Instrument(contract, symbol)
	} else {
		coins = klineToCoin(self.Symbol(symbol), interval, klines)
	}
	return coins, err
}

func (self *OkexFuture) Order(conrtact base.CONTRACT_PERIOD, symbol base.SYMBOL, _type base.ORDER, price float32, size int32) bool { // 下单
	instrumentid := self.GetInstrument(conrtact, symbol)

	var operation int32
	switch _type {
	case base.BUY_LONG:
		operation = 1
		break
	case base.BUY_SHORT:
		operation = 2
		break
	case base.SELL_LONG:
		operation = 3
		break
	case base.SELL_SHORT:
		operation = 4
		break
	}

	success := false // 下单成功
	result, err := self.Future.Order(instrumentid, operation, 2, price, size, 0)
	if err != nil {
		log.Println("[Buy] err: ", err)
	} else if result.Result { // 张数为0 result.Result=false 未能立马全部成交，返回的数据跟成交成功一样，不能区分
		// 未能立马全部成交，返回的数据跟成交成功一样，不能区分
		success = true
	} else {
		log.Print("[Buy] result = false, ", result)
	}
	return success
}

//  k线数据 -> orm.Coin
//  okex 的kline 是倒序的，最近的时间的在最前面
func klineToCoin(symbol string, section base.PERIOD, kline FutureCandles) []orm.Coin {
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
				Period:     base.Period(section),
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
