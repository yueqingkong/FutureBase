package trade

import (
	"github.com/yueqingkong/FutureBase/base"
	"github.com/yueqingkong/FutureBase/orm"
	"github.com/yueqingkong/FutureBase/router"
	"log"
	"time"
)

type BaseTrade struct {
	Mode           func() string // test|release
	DB             func() (string, string, string)
	InitKeys       func() []string
	Account        func() map[string]string
	SymbolPair     func() base.SYMBOL          // 合约交易对
	ContractPeriod func() base.CONTRACT_PERIOD // 交易的合约类型(当周,次周，季度)
	Port           func() string
	Plat           func() base.PlatBase // 合约交易平台
}

// 策略核心
type FutureStrategy interface {
	Name() string
	Tick(plat base.PlatBase, contract base.CONTRACT_PERIOD, symbol base.SYMBOL, price float32, t time.Time)
	Buy(plat base.PlatBase, contract base.CONTRACT_PERIOD, symbol base.SYMBOL, op base.ORDER, price float32, explain string, t time.Time)
}

// 运行策略
func (self BaseTrade) Start(strategy FutureStrategy) {
	plat := self.Plat()
	plat.InitKeys(self.InitKeys())

	syncMap := orm.NewSyncMap()

	// 服务
	model := self.Mode()
	syncMap.SetModel(model)

	// 数据库
	orm.ConnectSQL(self.DB())

	// 同步账号
	for k, v := range self.Account() {
		Account(k, v)
	}

	// (初始化)同步合约id
	contract := self.ContractPeriod()
	symbol := self.SymbolPair()
	plat.GetInstrument(contract, symbol)

	receive := make(chan string)

	go func() {
		// API限速规则：20次/2s
		ticker := time.NewTicker(time.Second * 2)

		for range ticker.C {
			log.Println()
			log.Println("--------------------------")
			log.Print("[create-time]", time.Now())

			self.Pulls(plat, contract, symbol, base.MIN_15, base.MIN_30, base.HOUR_4, base.HOUR_6, base.HOUR_12, base.DAY_1)

			priceFloat := self.price(plat, contract, symbol)
			if priceFloat == 0.0 { // 网络异常时,价格为0
				continue
			}
			syncMap.SetPrice(priceFloat)

			start := time.Now()

			select {
			case msg := <-receive:
				if msg == "buy_long" {
					log.Print("[buy_long]")
					strategy.Buy(plat, contract, symbol, base.BUY_LONG, priceFloat, "[开多]", start)
				} else if msg == "buy_short" {
					log.Print("[buy_short]")
					strategy.Buy(plat, contract, symbol, base.BUY_SHORT, priceFloat, "[开空]", start)
				} else if msg == "sellout" {
					log.Print("[sellout]")
					xorm := orm.NewXOrm()

					s := plat.Symbol(symbol)
					account := xorm.Account(s)
					account.Balance = account.Total
					account.Buy = 0.0
					xorm.UpdateAccount(account)

					xorm.ClearRecords()
				}
				break
			default:
				log.Println("[no receive]")
				break
			}

			self.Tick(plat, contract, symbol, strategy, priceFloat, start)
		}
	}()

	// http
	port := self.Port()
	r := router.NewRouter()
	r.Http(receive, port)
}

func (self BaseTrade) Tick(plat base.PlatBase, contract base.CONTRACT_PERIOD, symbol base.SYMBOL, strategy FutureStrategy, priceFloat float32, start time.Time) {
	if self.delivering(contract, symbol) { // 正在交割时间
		log.Print("[合约交割]")
		s := plat.Symbol(symbol)

		xorm := orm.NewXOrm()
		records := xorm.LastRecord(s, strategy.Name(), 1)
		if len(records) == 0 || records[0].Operation == 3 || records[0].Operation == 4 { // 未持有仓位

		} else {
			record := records[0]
			lastOp := record.Operation
			if lastOp == 1 {
				Sell(plat, contract, symbol, strategy.Name(), base.SELL_LONG, priceFloat, start)
			} else if lastOp == 2 {
				Sell(plat, contract, symbol, strategy.Name(), base.SELL_SHORT, priceFloat, start)
			}
		}
	} else {
		strategy.Tick(plat, contract, symbol, priceFloat, start)
	}
}

// 最新价格
func (self BaseTrade) price(plat base.PlatBase, contract base.CONTRACT_PERIOD, symbol base.SYMBOL) float32 {
	// 获取当前价格
	priceLimits := plat.Price(contract, symbol)
	if priceLimits == 0 { // 获取失败(可能 合约id为错误)
		plat.Instrument(contract, symbol)
		return 0.0
	}

	return priceLimits
}

// 交割前,一般是10-15分钟不能进行操作的
// 交割后，一般是30分钟左右
// 1 交割前 2 交割进行中 3 交割完成
func (self BaseTrade) delivering(period base.CONTRACT_PERIOD, symbol base.SYMBOL) bool {
	t := time.Now()
	b, begin, end := self.Plat().Delivery(period, symbol)

	return b && t.After(begin) && t.Before(end)
}

func (self BaseTrade) Pulls(plat base.PlatBase, contract base.CONTRACT_PERIOD, symbol base.SYMBOL, srctions ...base.PERIOD) {
	for _, value := range srctions {
		self.PullHistory(plat, contract, symbol, value)
	}
}

// 同步 历史kline数据
func (self BaseTrade) PullHistory(plat base.PlatBase, contract base.CONTRACT_PERIOD, symbol base.SYMBOL, section base.PERIOD) {
	s := self.Plat().Symbol(symbol)

	xorm := orm.XOrm{}
	lastCoin := xorm.Last(s, base.Period(section))

	var startTime time.Time
	if lastCoin.Symbol == "" { // 记录为空
		startTime = time.Time{}
	} else {
		startTime = lastCoin.CreateTime
	}

	// 最后一条记录是昨天的
	diffHours := time.Now().Sub(startTime).Hours()

	// 是否最新的数据
	if section == base.MIN_15 {
		if diffHours < 0.5 {
			return
		}
	} else if section == base.MIN_30 {
		if diffHours < 1 {
			return
		}
	} else if section == base.HOUR_1 {
		if diffHours < 2 {
			return
		}
	} else if section == base.HOUR_2 {
		if diffHours < 4 {
			return
		}
	} else if section == base.HOUR_4 {
		if diffHours < 8 {
			return
		}
	} else if section == base.HOUR_6 {
		if diffHours < 12 {
			return
		}
	} else if section == base.HOUR_12 {
		if diffHours < 24 {
			return
		}
	} else if section == base.DAY_1 {
		if diffHours < 48 {
			return
		}
	}

	log.Print("[PullHistory] ", symbol, "  diffHours: ", diffHours, " section: ", section)

	// 避免重复返回最后一条的k线数据，加30s
	if !startTime.IsZero() {
		startTime = startTime.Add(time.Duration(30) * time.Second)
	}

	coins, err := plat.KLine(contract, symbol, section, startTime)
	if err != nil { // 重新获取 instrumentid
		log.Println(err)
	} else {
		if len(coins) == 0 {
			log.Print("[PullHistory] len(coins) == 0 ")
		} else {
			xorm.InsertCoins(coins)
			self.PullHistory(plat, contract, symbol, section)
		}
	}
}
