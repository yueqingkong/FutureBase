package orm

import (
	"fmt"
	"github.com/yueqingkong/FutureBase/util"
	"sync"
)

var (
	syncMap sync.Map
)

type SyncMap struct {
}

func NewSyncMap() SyncMap {
	return SyncMap{}
}

//////////////////////////////   基本方法  /////////////////////////////////////
func (sync SyncMap) getValue(key string) string {
	k, _ := syncMap.Load(key)

	var value string
	if k == nil {
		value = ""
	} else {
		value = k.(string)
	}
	return value
}

func (sync SyncMap) setValue(key string, value string) {
	syncMap.Store(key, value)
}

//////////////////////////////   使用  /////////////////////////////////////
func (sync SyncMap) GetPrice() float32 {
	key := fmt.Sprintf("PRICE")
	value := sync.getValue(key)
	return util.StringToFloat32(value)
}

func (sync SyncMap) SetPrice(b float32) {
	sync.setValue("PRICE", util.Float32ToString(b))
}

func (sync SyncMap) GetPairRange() string {
	key := fmt.Sprintf("OKEX_PAIR_RANGE")
	value := sync.getValue(key)
	return value
}

func (sync SyncMap) SetPairRange(pair string) {
	sync.setValue("OKEX_PAIR_RANGE", pair)
}

func (sync SyncMap) GetPairTrend() string {
	key := fmt.Sprintf("OKEX_PAIR_TREND")
	value := sync.getValue(key)
	return value
}

func (sync SyncMap) SetPairTrend(pair string) {
	sync.setValue("OKEX_PAIR_TREND", pair)
}

func (sync SyncMap) GetInstrument(symbol string) string {
	key := fmt.Sprintf("OKEX_SYMBOL_%s", symbol)
	value := sync.getValue(key)
	return value
}

func (sync SyncMap) SetInstrument(symbol string, value string) {
	key := fmt.Sprintf("OKEX_SYMBOL_%s", symbol)
	sync.setValue(key, value)
}


func (sync SyncMap) Model() string {
	key := fmt.Sprintf("OKEX_MODEL")
	value := sync.getValue(key)
	return value
}

func (sync SyncMap) SetModel( value string) {
	key := fmt.Sprintf("OKEX_MODEL")
	sync.setValue(key, value)
}