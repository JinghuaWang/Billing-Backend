package handler

import (
	"billing/DAO"

	"github.com/gin-gonic/gin"
)

// AddTXHandler adds a transaction record in the database
func AddTXHandler(c *gin.Context) {
	var req AddTXReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		InvalidResp(c, "Invalid Parameters")
		return
	}

	if err = AddTX(&req); err != nil {
		ErrResp(c, err)
		return
	}
	SuccResp(c)
}

func ListHandler(c *gin.Context) {
	var req ListReq
	err := c.ShouldBindQuery(&req)
	if err != nil {
		InvalidResp(c, "Invalid Parameters")
		return
	}

	result, err := ListAll(&req)
	if err != nil {
		ErrResp(c, err)
		return
	}

	DataResp(c, result)
}

func ListBalanceHandler(c *gin.Context) {
	result, err := ListBalance()
	if err != nil {
		ErrResp(c, err)
		return
	}

	DataResp(c, result)
}

func PayHelper(c *gin.Context) {
	err := DAO.DB().Model(&DAO.Transaction{}).
		Where("status = 1").
		Update("status", 2).Error
	if err != nil {
		ErrResp(c, &Error{500, "DB Error"})
		return
	}
	SuccResp(c)
}

// ResetHandler deletes all the record in the databse
func ResetHandler(c *gin.Context) {
	err := DAO.DB().Exec("DELETE FROM transactions;").Error
	err = DAO.DB().Exec("DELETE FROM liabilities;").Error
	if err != nil {
		ErrResp(c, &Error{500, "DB Error"})
		return
	}
	SuccResp(c)
}

func InvalidResp(c *gin.Context, msg string) {
	c.JSON(400, Response{Code: 400, Message: msg})
}

func ErrResp(c *gin.Context, e error) {
	Error, ok := e.(*Error)
	if !ok {
		c.JSON(400, Response{Code: 400, Message: "Error Occurred"})
		return
	}

	c.JSON(Error.Code, Response{Code: Error.Code, Message: Error.Msg})
}

func SuccResp(c *gin.Context) {
	c.JSON(200, Response{Code: 200, Message: "Success"})
}

func DataResp(c *gin.Context, data interface{}) {
	c.JSON(200, gin.H{"data": data})
}
