package handler

import (
	"billing/DAO"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
)

func AddTX(c *gin.Context) {
	// Parse the request parameters
	var req AddTXReq
	err := c.ShouldBind(&req)
	if err != nil {
		ErrResp(c, 400, "Invalid Parameters")
		return
	}

	if !isUser(req.Spender) {
		ErrResp(c, 400, "Invalid Spender")
		return
	}

	// validate all payers
	if len(req.Payers) == 0 {
		ErrResp(c, 400, "Should have at least one payer")
		return
	}
	for i, p := range req.Payers {
		name := strings.ToLower(p)
		name = strings.TrimSpace(name)
		fmt.Println(name)
		req.Payers[i] = name
		if !isUser(name) {
			ErrResp(c, 400, "Invalid Payer")
			return
		}
	}

	// computer the liability for each payer
	var amount = req.Price / 3
	var liabilities []DAO.Liability
	for _, p := range req.Payers {
		liabilities = append(liabilities, DAO.Liability{
			Payer:  p,
			Amount: amount,
		})
	}

	// save to database
	error := DAO.DB().Create(&DAO.Transaction{
		Price:       req.Price,
		Spender:     req.Spender,
		Memo:        req.Memo,
		Liabilities: liabilities,
	}).Error
	if error != nil {
		ErrResp(c, 500, "DB Error")
		return
	}
	SuccResp(c)
}

func ListAll(c *gin.Context) {
	var txs []DAO.Transaction
	if err := DAO.DB().Find(&txs).Error; err != nil {
		ErrResp(c, 500, "DB Error")
		return
	}
	c.JSON(200, gin.H{"data": txs})
}

func isUser(name string) bool {
	users := [3]string{"john", "kelvin", "jun"}
	for _, user := range users {
		if user == name {
			return true
		}
	}
	return false
}

func ErrResp(c *gin.Context, code int, msg string) {
	c.JSON(code, Response{Code: code, Message: msg})
}

func SuccResp(c *gin.Context) {
	c.JSON(200, Response{Code: 200, Message: "Success"})
}
