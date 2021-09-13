package models

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

type Bid struct {
	gorm.Model
	Value     uint32 `json:"bid_amount"`
	AuctionID uint   `json:"auction_id"`
	Bidder    User   `json:"bidder"`
	BidderID  uint32 `gorm:"not null" json:"bidder_id"`
}

func (b *Bid) Prepare() {
	b.CreatedAt = time.Now()
	b.UpdatedAt = time.Now()
}

func (b *Bid) Validate(db *gorm.DB) map[string]string {
	var errorMessages = make(map[string]string)
	var err error

	if b.Value == 0 {
		err = errors.New("Required Value")
		errorMessages["Required_value"] = err.Error()
	}
	auction := Auction{}
	err = db.Debug().Model(Auction{}).Where("id = ?", b.AuctionID).Take(&auction).Error
	if err != nil {
		err = errors.New("Error, auction not found")
		errorMessages["Auction_not_found"] = err.Error()
	}
	if b.Value < auction.HighestBid {
		err = errors.New("Bid must be higher than the highest bid")
		errorMessages["Higher_bid"] = err.Error()
	}
	return errorMessages
}

func (b *Bid) AddBid(db *gorm.DB) (*Bid, error) {
	err := db.Debug().Model(&Bid{}).Create(&b).Error
	if err != nil {
		return &Bid{}, err
	}
	if b.BidderID != 0 {
		err = db.Debug().Model(&User{}).Where("id = ?", b.BidderID).Take(&b.Bidder).Error
		if err != nil {
			return &Bid{}, err
		}
	}
	if b.AuctionID != 0 {
		err = db.Debug().Model(&Auction{}).Where("id = ?", b.AuctionID).Update("highest_bid", b.Value).Error
		if err != nil {
			return &Bid{}, err
		}

	}
	return b, nil
}

func (b *Bid) GetBids(db *gorm.DB, pid uint64) (*[]Bid, error) {

	bids := []Bid{}
	err := db.Debug().Model(&Bid{}).Where("auction_id = ?", pid).Order("created_at desc").Find(&bids).Error
	if err != nil {
		return &[]Bid{}, err
	}
	if len(bids) > 0 {
		for i, _ := range bids {
			err := db.Debug().Model(&User{}).Where("id = ?", bids[i].BidderID).Take(&bids[i].Bidder).Error
			if err != nil {
				return &[]Bid{}, err
			}
		}
	}
	return &bids, err
}
