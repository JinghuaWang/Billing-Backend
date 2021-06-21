package handler

type AddTXReq struct {
	Price   float32  `form:"price" binding:"required"`
	Spender string   `form:"spender" binding:"required"`
	Payers  []string `form:"payers" binding:"required"`
	Memo    string   `form:"memo"`
}

type ListReq struct {
	Unpaid bool `form:"unpaid"`
}

type Response struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type BalanceResp struct {
	Asset   float32 `json:"amount"`
	Deficit float32 `json:"deficit"`
	Balance float32 `json:"balance"`
}

type AccountResp struct {
	Name   string  `json:"name"`
	Amount float32 `json:"amount"`
}

type Error struct {
	Code int
	Msg  string
}

func New(code int, msg string) *Error {
	return &Error{
		Code: code,
		Msg:  msg,
	}
}

func (e *Error) Error() string {
	return e.Msg
}
