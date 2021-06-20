package main

import (
	"billing/DAO"
	"billing/handler"

	"github.com/gin-gonic/gin"
)

func main() {
	DAO.Init()
	r := gin.Default()
	r.GET("/add", handler.AddTX)
	r.GET("/list/all", handler.ListAll)
	r.Run(":3000") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
