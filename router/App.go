package router

import (
	"gitee.com/shieldpu_futures/FutureBase/orm"
	"github.com/gin-gonic/gin"
)

/**
 * 账户
 */
func Account(context *gin.Context) {
	xorm := orm.NewXOrm()
	accounts := xorm.Accounts()

	context.JSON(200, gin.H{
		"code":    2000,
		"message": "",
		"data":    accounts,
	})
}

/**
 * 交易记录
 */
func Records(context *gin.Context) {
	xorm := orm.NewXOrm()
	records := xorm.RecordsAll()

	context.JSON(200, gin.H{
		"code":    2000,
		"message": "",
		"data":    records,
	})
}