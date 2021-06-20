package DAO

import "gorm.io/gorm"

type Transaction struct {
	gorm.Model
	Price       float32
	Spender     string      `gorm:"type:varchar(30); not null"` // the person who pays for this transaction
	Memo        string      `gorm:"type:varchar(200); not null"`
	Status      uint8       `gorm:"default:1; not null"` // 1:unpaid 2: paid
	Liabilities []Liability // Has Many relationship with liability
}
