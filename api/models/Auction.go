package models

import (
	"time"

	"gorm.io/gorm"
)

type Auction struct {
	gorm.Model
	ProductID    uint
	IsBuyNow     bool      `gorm:"not null;default:false" json:"is_buy_now"`
	StartBid     uint32    `gorm:"not null" json:"start_bid,string,omitempy"`
	HighestBid   uint32    `json:"highest_bid"`
	BidIncrement uint32    `gorm:"not null" json:"bid_increment,string,omitempty"`
	BuyNowPrice  uint32    `json:"buy_now_price,string,omitempty"`
	IsExpired    bool      `gorm:"not null;default:false" json:"is_expired"`
	ExpiredTime  time.Time `json:"expired_time"`
	Bids         []Bid     `json:"bids"`
}

func (a *Auction) Prepare() {
	a.CreatedAt = time.Now()
	a.UpdatedAt = time.Now()
	a.HighestBid = 0
}

func (a *Auction) AddAuction(db *gorm.DB) (*Auction, error) {
	var err error
	err = db.Debug().Model(&Auction{}).Create(&a).Error
	if err != nil {
		return &Auction{}, err
	}
	return a, nil
}

func (p *Auction) FindAllAuction(db *gorm.DB) (*[]Auction, error) {
	var err error
	auctions := []Auction{}
	err = db.Debug().Model(&Auction{}).Limit(100).Order("created_at desc").Find(&auctions).Error
	if err != nil {
		return &[]Auction{}, err
	}
	// if len(auctions) > 0 {
	// 	for i, _ := range auctions {
	// 		err := db.Debug().Model(&Product{}).Where("id = ?", auctions[i].ProductID).Take(&auctions[i].Product).Error
	// 		if err != nil {
	// 			return &[]Auction{}, err
	// 		}
	// 	}
	// }
	return &auctions, nil
}

func (a *Auction) FindAuctionByID(db *gorm.DB, pid uint64) (*Auction, error) {
	var err error
	err = db.Debug().Model(&Auction{}).Where("id = ?", pid).Take(&a).Error
	if err != nil {
		return &Auction{}, err
	}
	return a, nil
}
