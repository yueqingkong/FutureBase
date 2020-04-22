package trade

import (
	"gitee.com/shieldpu_futures/FutureBase/base"
	"log"
	"time"
)

type Simple struct {
}

func (self *Simple) Name() string {
	return "simple"
}
func (self *Simple) Tick(plat base.PlatBase, contract base.CONTRACT_PERIOD, symbol base.SYMBOL, price float32, t time.Time) {
	log.Printf("[%s] [price] %f", self.Name(), price)
}

func (self *Simple) Buy(plat base.PlatBase, contract base.CONTRACT_PERIOD, symbol base.SYMBOL, op base.ORDER, price float32, explain string, t time.Time) {
}
