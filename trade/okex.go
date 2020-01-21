package trade

import (
	"gitee.com/shieldpu_futures/FutureBase/orm"
	"gitee.com/shieldpu_futures/FutureBase/router"
	"gitee.com/shieldpu_futures/FutureBase/util"
	"github.com/yueqingkong/Okex/plat"
	"log"
	"time"
)

type OkexTrade struct {
	Mode    func() string // test|release
	DB      func() (string, string, string)
	Account func() (map[string]string)
	Symbols func() []string
	Keys    func() (string, string, string) // 交易key
	Port    func() string
}

type FutureStrategy interface {
	Name() string
	Tick(symbol string, price float32, t time.Time)
	Buy(symbol string, op int32, price float32, explain string,t time.Time)
}

/**
 * (趋势) 实时数据
 */
func (trade OkexTrade) StartTrade(strategy FutureStrategy) {
	// 交易key
	okex := plat.NewOKex()
	okex.InitKeys(trade.Keys())

	syncMap := orm.NewSyncMap()

	// 服务
	model := trade.Mode()
	syncMap.SetModel(model)

	// 数据库
	orm.ConnectSQL(trade.DB())

	// 同步账号
	for k, v := range trade.Account() {
		Account(k, v)
	}

	for _, symbol := range trade.Symbols() {
		Instrument(symbol) // 同步合约
	}

	go func() {
		// API限速规则：20次/2s
		ticker := time.NewTicker(time.Second * 2)
		for range ticker.C {
			log.Println()
			log.Println("--------------------------")
			log.Print("[create-time]", time.Now())

			for _, symbol := range trade.Symbols() {
				// 同步最新K线
				Pulls(symbol, "6h", "12h", "1d")
				trade.TickTrend(symbol, strategy)
			}
		}
	}()

	// http
	port := trade.Port()
	r := router.NewRouter()
	r.Http(port)
}

func (trade OkexTrade) TickTrend(symbol string, strategy FutureStrategy) {
	priceFloat := price(symbol)
	if priceFloat == 0.0 { // 网络异常时,价格为0
		return
	}

	start := time.Now()
	enforceStatus := enforcementSell(symbol)
	if enforceStatus {
		log.Print("[合约交割]")
		xorm := orm.NewXOrm()
		records := xorm.LastRecord(symbol, strategy.Name(), 1)
		if len(records) == 0 || records[0].Operation == 3 || records[0].Operation == 4 { // 未持有仓位

		} else {
			record := records[0]
			lastOp := record.Operation
			if lastOp == 1 {
				Sell(symbol, strategy.Name(), 3, priceFloat, start)
			} else if lastOp == 2 {
				Sell(symbol, strategy.Name(), 4, priceFloat, start)
			}
		}
	} else {
		strategy.Tick(symbol, priceFloat, start)
	}
}

/**
 * 最新价格
 */
func price(symbol string) float32 {
	syncMap := orm.NewSyncMap()
	instrumentid := syncMap.GetInstrument(symbol) // 合约id
	if instrumentid == "" {
		instrumentid = Instrument(symbol)
	}

	// 获取当前价格
	future := plat.NewOKexFuture()
	priceLimits, err := future.Trades(instrumentid, 1)
	if err != nil || len(priceLimits) == 0 { // 获取失败
		instrumentid = Instrument(symbol)
		return 0.0
	}

	priceFloat := util.StringToFloat32(priceLimits[0].Price)
	return priceFloat
}

/**
 * 交割前,一般是10-15分钟不能进行操作的
 * 交割后，一般是30分钟左右
 * 1 交割前 2 交割进行中 3 交割完成
 */
func enforcementSell(symbol string) bool {
	status := false

	xorm := orm.NewXOrm()
	instrument := xorm.Instrument(symbol)
	if instrument.Symbol != "" {
		t := time.Now()

		// 交割日期
		delTime := instrument.Delivery
		log.Println("[交割日] ", delTime.Weekday(), delTime.Hour(), delTime.Minute())

		if delTime.Year() == t.Year() && delTime.Month() == t.Month() && delTime.Day() == t.Day() {
			// 交割时间(周五 下午4点)
			if t.Weekday() == time.Friday && t.Hour() >= 15 && t.Hour() <= 16 {
				if t.Hour() == 15 && t.Minute() >= 40 {
					status = true
				} else if t.Hour() == 16 && t.Minute() < 30 {
					status = true
				} else if t.Hour() == 16 && t.Minute() == 30 {
					if t.Second() > 10 || t.Second() < 20 { // 刷新合约信息
						Instrument(symbol)
					}
				}
			}
		}
	}
	return status
}
