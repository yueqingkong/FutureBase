package router

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
)

type Router struct {
	ch chan string
}

var (
	router *Router
)

func NewRouter() *Router {
	if router == nil {
		router = new(Router)
	}
	return router
}

func (r *Router) Receive(msg string) {
	r.ch <- msg
}

func (r *Router) Http(c chan string, port string) {
	r.ch = c

	engine := gin.Default()

	r.v1(engine)
	err := engine.Run(fmt.Sprintf(":%s", port))
	if err != nil {
		log.Print("[Http] engine Run: ", err)
	}
}

// API
func (r *Router) v1(engine *gin.Engine) {
	app := engine.Group("/trend")
	{
		app.GET("/account", Accounts)
		app.GET("/records", Records)

		app.GET("/price", Price)
		app.GET("/buyin", Buyin)
		app.GET("/sellout", Sellout)
		app.GET("/kline", KLine)
	}
}
