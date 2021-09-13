package models

import (
	"gorm.io/gorm"
)

type OrderStatus uint

const (
	AuctionWon OrderStatus = iota + 1
	DealPending
	DealSuccess
	DealFailed
)

type Order struct {
	gorm.Model
	Status     OrderStatus
	Value      uint32      `gorm:"not null;default:false" json:"value"`
	Auction    Auction     `json:"auction"`
	OrderChats []OrderChat `json:"order_chats"`
}
