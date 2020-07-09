package router

import (
	"github.com/gin-gonic/gin"
	"github.com/yueqingkong/FutureBase/orm"
	"github.com/yueqingkong/FutureBase/util"
)

// 账户信息
func Account(context *gin.Context) {
	symbol := context.Query("symbol")
	if symbol == "" {
		symbol = "btc"
	}

	xorm := orm.NewXOrm()
	accounts := xorm.Account(symbol)

	context.JSON(200, gin.H{
		"code":    2000,
		"message": "",
		"data":    accounts,
	})
}

// 交易记录
func Records(context *gin.Context) {
	symbol := context.Query("symbol")
	if symbol == "" {
		symbol = "btc"
	}

	xorm := orm.NewXOrm()
	records := xorm.RecordsAll(symbol)

	context.JSON(200, gin.H{
		"code":    2000,
		"message": "",
		"data":    records,
	})
}

// kline 信息
func KLine(context *gin.Context) {
	symbol := context.Query("symbol")
	period := context.Query("period")
	timestamp := context.Query("timestamp")

	var coins []orm.Coin
	if period == "15m" || period == "30m" || period == "1h" || period == "2h" || period == "4h" || period == "6h" || period == "12h" || period == "1d" {
		t := util.TimestampToTime(util.StringToInt64(timestamp))
		xorm := orm.NewXOrm()
		coins = xorm.Next(symbol, period, t)
	}

	context.JSON(200, gin.H{
		"code":    2000,
		"message": "",
		"data":    coins,
	})
}
