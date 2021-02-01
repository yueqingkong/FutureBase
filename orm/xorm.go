package orm

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"log"
	"net/url"
	"time"
)

// K线数据
type Coin struct {
	Id         int64
	Symbol     string    `xorm:"varchar(255) index index(symbol,period) index(symbol,period,timestamp)"`
	Plat       string    `xorm:"varchar(255) index"`
	Period     string    `xorm:"varchar(255) index index(symbol,period) index(symbol,period,timestamp)"` // 时间间隔
	Open       float32   `xorm:"float"`
	Close      float32   `xorm:"float"`
	High       float32   `xorm:"float"`
	Low        float32   `xorm:"float"`
	Volume     float32   `xorm:"float"`
	Timestamp  int64     `xorm:"bigint index index(symbol,period) index(symbol,period,timestamp)"` // 秒
	CreateTime time.Time `xorm:"DATETIME"`
}

// 账户
type Account struct {
	Id      int64
	Symbol  string  `xorm:"varchar(255) index"` // Token
	Balance float32 `xorm:"float"`              // 可用Token余额(张)
	Buy     float32 `xorm:"float"`              // 已使用Token(张)
	Total   float32 `xorm:"float"`              // 总值
}

// 交易记录
type Record struct {
	Id         int64
	Symbol     string    `xorm:"varchar(255) index index(symbol,strategy)"` // Token
	Strategy   string    `xorm:"varchar(255) index index(symbol,strategy)"` // 策略名称
	Operation  int32     `xorm:"int"`                                       // 1: 开多 2: 开空 3: 平仓
	Position   int32     `xorm:"int"`                                       // 加仓层数
	Price      float32   `xorm:"float"`                                     // 当前价格
	AvgPrice   float32   `xorm:"float"`                                     // 均价
	Used       float32   `xorm:"float"`                                     // 已开仓Token
	Size       float32   `xorm:"float"`                                     // 开仓张数
	Total      float32   `xorm:"float"`                                     // 当前账户总值
	Explain    string    `xorm:"text"`                                      // 描述 usd->token | ust<-token
	Profit     float32   `xorm:"float"`                                     // 收益
	ProfitRate float32   `xorm:"float"`                                     // 收益率
	Timestamp  int64     `xorm:"bigint index"`
	CreateTime time.Time `xorm:"DateTime index"` // 时间
}

// key-value
type Redis struct {
	Id int64
	K  string `xorm:"varchar(255) notnull unique"` // key
	V  string `xorm:"varchar(255)"`                // value
}

// 推送设备
type Device struct {
	Id    int64
	Name  string `xorm:"varchar(255) notnull unique"` // 设备名称
	Token string `xorm:"varchar(255)"`                // 推送Token
}

// 合约信息
type Instrument struct {
	Id       int64
	Symbol   string    `xorm:"varchar(255) index"`
	Period   string    `xorm:"varchar(255)"`
	Key      string    `xorm:"varchar(255)"` // 合约id
	Delivery time.Time `xorm:"DateTime"`     // 交割时间
}

var engine *xorm.Engine

type XOrm struct {
}

// 连接数据库
func ConnectSQL(name, user, password string) {
	var err error

	// mysql配置
	// name, user, password := properties.ORM()
	sourceName := fmt.Sprintf("%s:%s@tcp(localhost:3306)/%s?charset=utf8&parseTime=true&loc=%s", user, password, name, url.QueryEscape("Asia/Shanghai"))
	log.Print(sourceName)

	engine, err = xorm.NewEngine("mysql", sourceName)
	if err != nil {
		log.Fatal("[MySql] 连接失败,", err)
	}

	engine.ShowSQL(false)
	err = engine.Sync2(new(Coin), new(Account), new(Record), new(Redis), new(Device), new(Instrument))
	if err != nil {
		log.Fatal("[MySql] 同步表失败", err)
	}
}

func NewXOrm() XOrm {
	return XOrm{}
}

///////////////////   insert ///////////////////////////////////////
func (orm XOrm) InsertInstrument(ins Instrument) {
	id, err := engine.Insert(&ins)
	if err != nil || id == 0 {
		log.Print("[InsertInstrument]", err, ins)
	} else {
		log.Print("[InsertInstrument] id = ", id)
	}
}

func (orm XOrm) InsertRedis(redis Redis) {
	id, err := engine.Insert(&redis)
	if err != nil || id == 0 {
		log.Print("[InsertRedis]", err, redis)
	} else {
		log.Print("[InsertRedis] id = ", id)
	}
}

func (orm XOrm) InsertAccount(account Account) {
	id, err := engine.Insert(&account)
	if err != nil || id == 0 {
		log.Print("[InsertAccount]", err, account)
	} else {
		log.Print("[InsertAccount] id = ", id)
	}
}

func (orm XOrm) InsertCoins(coins []Coin) {
	id, err := engine.Insert(&coins)
	if err != nil || id == 0 {
		log.Print("[InsertCoins]", err, coins)
	} else {
		log.Print("[InsertCoins] id = ", id)
	}
}

func (orm XOrm) InsertRecord(symbol string, stratege string, op int32, position int32, price float32, average float32, used float32, size float32, total float32, profit float32, rate float32, explain string, t time.Time) {
	record := Record{
		Symbol:     symbol,
		Strategy:   stratege,
		Operation:  op,
		Position:   position,
		Price:      price,
		AvgPrice:   average,
		Used:       used,
		Size:       size,
		Total:      total,
		Profit:     profit,
		ProfitRate: rate,
		Explain:    explain,
		Timestamp:  t.Unix(),
		CreateTime: t,
	}

	id, err := engine.Insert(&record)
	if err != nil || id == 0 {
		log.Print("[InsertRecord]", err, record)
	} else {
		log.Print("[InsertRecord] id = ", id)
	}
}

///////////////////   find ///////////////////////////////////////
// 账户
func (orm XOrm) Account(symbol string) Account {
	accounts := make([]Account, 0)
	err := engine.Where("symbol = ?", symbol).
		Limit(1).
		Find(&accounts)

	if err != nil {
		log.Print("[Account]", err, "[symbol]", symbol)
	}

	var account Account
	if len(accounts) != 0 {
		account = accounts[0]
	}
	return account
}

func (orm XOrm) Accounts() []Account {
	accounts := make([]Account, 0)
	err := engine.Find(&accounts)

	if err != nil {
		log.Print("[Accounts]", err)
	}
	return accounts
}

// (正序) 指定时间前几日日线
func (orm XOrm) Before(symbol string, period string, limit int32) []Coin {
	coins := make([]Coin, 0)
	err := engine.Where("symbol = ? and period = ?", symbol, period).
		Desc("timestamp").
		Limit(int(limit)).
		Find(&coins)

	if err != nil {
		log.Print("[Before]", err)
	}

	// 重新排序
	newCoins := make([]Coin, 0)
	for index := len(coins) - 1; index >= 0; index-- {
		newCoins = append(newCoins, coins[index])
	}
	return newCoins
}

// (正序)前几日日线
func (orm XOrm) BeforeBy(symbol string, period string, start time.Time, limit int32) []Coin {
	coins := make([]Coin, 0)

	timestamp := start.Unix()
	err := engine.Where("symbol = ? and period = ? and timestamp < ?", symbol, period, timestamp).
		Desc("timestamp").
		Limit(int(limit)).
		Find(&coins)

	if err != nil {
		log.Print("[Before]", err)
	}

	// 重新排序
	newCoins := make([]Coin, 0)
	for index := len(coins) - 1; index >= 0; index-- {
		newCoins = append(newCoins, coins[index])
	}
	return newCoins
}

// 上一根kline
func (orm XOrm) Last(symbol string, period string) Coin {
	coins := make([]Coin, 0)
	err := engine.Where("symbol = ? and period = ?", symbol, period).
		Desc("timestamp").
		Limit(1).
		Find(&coins)

	if err != nil {
		log.Print("[Last]", err)
	}

	var lastCoin Coin
	if len(coins) > 0 {
		lastCoin = coins[0]
	}
	return lastCoin
}

// 查询给定时间之后的k线
// start 为空，返回所有
func (orm XOrm) Next(symbol string, period string, timestamp int64) []Coin {
	coins := make([]Coin, 0)

	var err error
	err = engine.Where("symbol = ? and period = ? and timestamp > ?", symbol, period, timestamp).
		Asc("timestamp").
		Limit(200).
		Find(&coins)

	if err != nil {
		log.Print("[Next]", err)
	}
	return coins
}

// 最近交易记录
func (orm XOrm) LastRecord(symbol string, strategy string, limit int32) []Record {
	records := make([]Record, 0)
	err := engine.Where("symbol = ? and strategy = ?", symbol, strategy).
		Desc("timestamp").
		Limit(int(limit)).
		Find(&records)

	if err != nil {
		log.Print("[LastRecord]", err, "[symbol] ", symbol)
	}

	return records
}

// 交易记录
func (orm XOrm) Records(symbol string, strategy string) []Record {
	records := make([]Record, 0)
	err := engine.Where("symbol = ? and strategy = ?", symbol, strategy).
		Find(&records)

	if err != nil {
		log.Print("[Records]", err, "[symbol] ", symbol)
	}
	return records
}

func (orm XOrm) RecordsAll(symbol string) []Record {
	records := make([]Record, 0)
	err := engine.Where("symbol = ?", symbol).Desc("timestamp").Find(&records)

	if err != nil {
		log.Print("[Records]", err)
	}
	return records
}

// 合约信息
func (orm XOrm) Instrument(symbol string) Instrument {
	instruments := make([]Instrument, 0)
	err := engine.Where("symbol = ?", symbol).
		Find(&instruments)

	if err != nil {
		log.Print("[Instrument]", err, "[symbol] ", symbol)
	}

	var instrument Instrument
	if len(instruments) > 0 {
		instrument = instruments[0]
	}
	return instrument
}

// 所有推送设备
func (xorm XOrm) Deveices() []Device {
	devices := make([]Device, 0)
	err := engine.Find(&devices)
	if err != nil {
		log.Print(err)
	}
	return devices
}

// key-value
func (orm XOrm) Redis(key string) Redis {
	rediss := make([]Redis, 0)
	err := engine.Where("k = ?", key).
		Limit(1).
		Find(&rediss)

	if err != nil {
		log.Print("[Redis]", err, "[k] ", key)
	}

	var redis Redis
	if len(rediss) > 0 {
		redis = rediss[0]
	}
	return redis
}

///////////////////   update ///////////////////////////////////////
func (orm XOrm) UpdateRedis(redis Redis) {
	_, err := engine.Id(redis.Id).Cols("v").Update(redis)
	if err != nil {
		log.Print("[UpdateRedis]", err, redis)
	}
}

func (orm XOrm) UpdateAccount(account Account) {
	_, err := engine.Id(account.Id).Cols("balance", "buy", "total").Update(account)
	if err != nil {
		log.Print("[UpdateAccount]", err, account)
	}
}

func (orm XOrm) UpdateInstrument(instrument Instrument) {
	_, err := engine.Id(instrument.Id).Cols("symbol", "key", "period", "delivery").Update(instrument)
	if err != nil {
		log.Print("[UpdateInstrument]", err, instrument)
	}
}

///////////////////   remove ///////////////////////////////////////
func (orm XOrm) ClearCoin() {
	sql := "delete from coin;"
	_, err := engine.Exec(sql)
	if err != nil {
		log.Print("[ClearCoin]", err)
	}
}

func (orm XOrm) ClearAccount() {
	sql := "delete from account;"
	_, err := engine.Exec(sql)
	if err != nil {
		log.Print("[ClearAccount]", err)
	}
}

// 清空交易记录
func (orm XOrm) ClearRecords() {
	sql := "delete from record where id >= 0;"
	_, err := engine.Exec(sql)
	if err != nil {
		log.Print("[Clear Records]", err)
	}
}

func (orm XOrm) ClearBackAssess() {
	sql := "delete from back_assess"
	_, err := engine.Exec(sql)
	if err != nil {
		log.Print("[ClearBackAssess]", err)
	}
}
