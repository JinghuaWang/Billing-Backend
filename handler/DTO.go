package handler

type AddTXReq struct {
	Price   float32  `form:"price" binding:"required"`
	Spender string   `form:"spender" binding:"required"`
	Payers  []string `form:"payers" binding:"required"`
	Memo    string   `form:"memo"`
}

type Response struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
