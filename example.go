package main

import (
	"billing/DAO"
	"billing/handler"

	"github.com/gin-gonic/gin"
)

func main() {
	DAO.Init()
	r := gin.Default()

	r.POST("/add", handler.AddTXHandler)

	r.GET("/list", handler.ListHandler)
	r.GET("/list/balance", handler.ListBalanceHandler)

	r.GET("/pay", handler.PayHelper)

	r.POST("/delete", handler.ResetHandler)

	r.Run(":3001") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
