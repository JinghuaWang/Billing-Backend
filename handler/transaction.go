package handler

import (
	"billing/DAO"
	"fmt"
	"strings"
)

func AddTX(req *AddTXReq) error {
	req.Spender = strings.ToLower(req.Spender)
	req.Spender = strings.TrimSpace(req.Spender)
	if !isUser(req.Spender) {
		return &Error{400, "Invalid Spender"}
	}

	// validate all payers
	if len(req.Payers) == 0 {
		return &Error{400, "Should have at least one payer"}
	}
	for i, p := range req.Payers {
		name := strings.ToLower(p)
		name = strings.TrimSpace(name)
		fmt.Println(name)
		req.Payers[i] = name
		if !isUser(name) {
			return &Error{400, "Invalid Payer"}
		}
	}

	// computer the liability for each payer
	var amount = req.Price / float32(len(req.Payers))
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
		return &Error{500, "DB Error"}
	}

	return nil
}

func ListAll(req *ListReq) ([]*DAO.Transaction, error) {
	var txs []*DAO.Transaction
	tx := DAO.DB().Find(&txs)
	// Only show unpaid transactions
	if req.Unpaid {
		tx.Where("status = 1")
	}

	if err := tx.Find(&txs).Error; err != nil {
		return nil, &Error{500, "DB Error"}
	}
	if err := FullTransaction(txs); err != nil {
		return nil, &Error{500, "DB Error"}
	}

	return txs, nil
}

// ListBalance groups the unpaid balance by name
func ListBalance() (map[string]*BalanceResp, error) {
	balances := make(map[string]*BalanceResp)
	balances["john"] = &BalanceResp{}
	balances["kelvin"] = &BalanceResp{}
	balances["jun"] = &BalanceResp{}

	var assets []AccountResp
	err := DAO.DB().Model(&DAO.Transaction{}).
		Select("spender as name, sum(price) as amount").
		Where("status = 1").
		Group("spender").
		Find(&assets).Error
	if err != nil {
		return nil, &Error{500, "DB Error"}
	}

	var liabs []AccountResp
	err = DAO.DB().Model(&DAO.Transaction{}).
		Select("liabilities.payer as name, sum(liabilities.amount) as amount").
		Joins("left join liabilities on transactions.id = liabilities.transaction_id ").
		Where("transactions.status = 1").
		Group("liabilities.payer").
		Find(&liabs).Error
	if err != nil {
		return nil, &Error{500, "DB Error"}
	}

	for _, a := range assets {
		person := balances[a.Name]
		person.Asset = a.Amount
	}
	for _, l := range liabs {
		person := balances[l.Name]
		person.Deficit = l.Amount
	}
	for _, person := range balances {
		person.Balance = person.Asset - person.Deficit
	}

	return balances, nil
}

/*** Helper Methods ***/

func FullTransaction(txs []*DAO.Transaction) error {
	if len(txs) == 0 {
		return nil
	}
	// map transaction id to index in the array;
	var ids []uint
	var idsToIndex map[uint]int = make(map[uint]int)
	for index, tx := range txs {
		idsToIndex[tx.ID] = index
		ids = append(ids, tx.ID)
	}

	var liabs []DAO.Liability
	if err := DAO.DB().Find(&liabs).Where("transaction_id IN ?", ids).Error; err != nil {
		return err
	}

	// Map each liability to corresponding transaction
	for _, l := range liabs {
		tx := txs[idsToIndex[l.TransactionID]]
		tx.Liabilities = append(tx.Liabilities, l)
	}
	return nil
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
