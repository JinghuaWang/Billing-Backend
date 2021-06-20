package DAO

import "gorm.io/gorm"

type Liability struct {
	gorm.Model
	TransactionID uint   `gorm:"not null"`                   // the primary key of transaction
	Payer         string `gorm:"type:varchar(30); not null"` // the person owe the amount
	Amount        float32
}
