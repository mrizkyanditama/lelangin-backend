package models

import (
	"gorm.io/gorm"
)

type OrderChat struct {
	gorm.Model
	OrderID  uint   `json:"order_id"`
	Sender   User   `json:"sender"`
	SenderID uint32 `gorm:"not null" json:"sender_id"`
	Message  string `gorm:"text;not null;" json:"message"`
}
