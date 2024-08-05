package postgres

import (
	"balance-service/config"
	"balance-service/dto"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"time"
)

type Client struct {
	db *gorm.DB
}

func New(cfg *config.PostgresConfig) (*Client, error) {
	db, err := gorm.Open(postgres.Open(cfg.Dsn), &gorm.Config{})

	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(&dto.Balance{}, &dto.Currency{})
	if err != nil {
		return nil, err
	}

	return &Client{
		db: db,
	}, nil
}

func (c *Client) GetBalanceByUserId(userId string) []dto.Balance {
	balances := make([]dto.Balance, 1)
	c.db.Preload("Currency").Where("user_id = ?", userId).Find(&balances)

	return balances
}

func (c *Client) UpdateUserBalance(userId, currency string, amount, locked_balance float64) error {
	err := c.db.Exec("INSERT INTO balances (user_id, currency_id, balance, locked_balance, updated_at) VALUES ($1,(SELECT id FROM currencies WHERE name = $2), $3,$4,$5)\nON CONFLICT (user_id, currency_id) DO UPDATE\nSET balance = balances.balance + $3, locked_balance = balances.locked_balance + $4, updated_at = $5;",
		userId, currency, amount, locked_balance, time.Now().Unix()).Error

	return err
}
