package router

import (
	"github.com/gin-gonic/gin"
	"github.com/yueqingkong/FutureBase/orm"
)

// 账户信息
func Accounts(context *gin.Context) {
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

// 当前价格
func Price(context *gin.Context) {
	symbol := context.Query("symbol")
	if symbol == "" {
		symbol = "btc"
	}

	syncMap := orm.NewSyncMap()
	price := syncMap.GetPrice()

	context.JSON(200, gin.H{
		"code":    2000,
		"message": "",
		"data":    price,
	})
}

// 加仓
func Buyin(context *gin.Context) {
	order := context.Query("order")

	if order == "long" {
		router.Receive("buy_long")
	} else if order == "short" {
		router.Receive("buy_short")
	}

	context.JSON(200, gin.H{
		"code":    2000,
		"message": "",
		"data":    "",
	})
}

// 平仓
func Sellout(context *gin.Context) {
	symbol := context.Query("symbol")
	if symbol == "" {
		symbol = "btc"
	}

	router.Receive("sellout")

	context.JSON(200, gin.H{
		"code":    2000,
		"message": "",
		"data":    "",
	})
}
