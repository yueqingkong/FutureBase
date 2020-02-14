package router

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
)

type Router struct {
}

func NewRouter() Router {
	return Router{

	}
}

func (r Router) Http(port string) {
	engine := gin.Default()

	r.v1(engine)
	err := engine.Run(fmt.Sprintf(":%s", port))
	if err != nil {
		log.Print("[Http] engine Run: ", err)
	}
}

/**
 * API
 */
func (r Router) v1(engine *gin.Engine) {
	app := engine.Group("/trend")
	{
		app.GET("/account", Account)
		app.GET("/records", Records)
		app.GET("/kline", KLine)
	}
}