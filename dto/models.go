package dto

import "github.com/google/uuid"

type Balance struct {
	Id            uuid.UUID `gorm:"type:uuid;default:gen_random_uuid()"`
	UserId        string
	CurrencyId    string
	Currency      Currency `gorm:"foreignKey:CurrencyId;references:Id"`
	Balance       float32
	LockedBalance float32
	UpdatedAt     int64
}

type Currency struct {
	Id   string `gorm:"primaryKey"`
	Name string
}
