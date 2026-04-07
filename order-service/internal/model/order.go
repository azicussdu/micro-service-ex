package model

import "time"

type Order struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	UserID      uint      `gorm:"index;not null" json:"user_id"`
	ProductName string    `gorm:"not null" json:"product_name"`
	CreatedAt   time.Time `json:"created_at"`
}
