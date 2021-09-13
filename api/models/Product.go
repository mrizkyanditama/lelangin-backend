package models

import (
	"errors"
	"html"
	"strings"
	"time"

	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	Name        string  `gorm:"size:255;not null;unique" json:"name"`
	Description string  `gorm:"text;not null;" json:"description"`
	Owner       User    `json:"owner"`
	OwnerID     uint32  `gorm:"not null" json:"owner_id"`
	CategoryID  uint    `json:"category_id"`
	Tags        []*Tag  `gorm:"many2many:product_tag;"`
	Auction     Auction `json:"auction"`
	PhotoPath   string  `gorm:"size:255;null;" json:"photo_path"`
}

func (p *Product) Prepare() {
	p.Name = html.EscapeString(strings.TrimSpace(p.Name))
	p.Description = html.EscapeString(strings.TrimSpace(p.Description))
	p.Owner = User{}
	p.CreatedAt = time.Now()
	p.UpdatedAt = time.Now()
}

func (p *Product) Validate() map[string]string {
	var err error

	var errorMessages = make(map[string]string)

	if p.Name == "" {
		err = errors.New("Required Name of Product")
		errorMessages["Required_name"] = err.Error()

	}
	if p.Description == "" {
		err = errors.New("Required Description")
		errorMessages["Required_description"] = err.Error()

	}
	if p.OwnerID < 1 {
		err = errors.New("Required Owner")
		errorMessages["Required_owner"] = err.Error()
	}
	if p.CategoryID == 0 {
		err = errors.New("Required Category")
		errorMessages["Required_category"] = err.Error()
	}
	return errorMessages

}

func (p *Product) AddProduct(db *gorm.DB) (*Product, error) {
	var err error
	err = db.Debug().Omit("Tags").Create(&p).Error
	if err != nil {
		return &Product{}, err
	}
	if p.ID != 0 {
		err = db.Debug().Model(&User{}).Where("id = ?", p.OwnerID).Take(&p.Owner).Error
		if err != nil {
			return &Product{}, err
		}
	}
	return p, nil
}

func (p *Product) FindAllProducts(db *gorm.DB) (*[]Product, error) {
	var err error
	products := []Product{}
	err = db.Debug().Model(&Product{}).Limit(100).Order("created_at desc").Find(&products).Error
	if err != nil {
		return &[]Product{}, err
	}
	if len(products) > 0 {
		for i, _ := range products {
			err := db.Debug().Model(&User{}).Where("id = ?", products[i].OwnerID).Take(&products[i].Owner).Error
			err = db.Debug().Model(&Auction{}).Where("id = ?", products[i].ID).Take(&products[i].Auction).Error
			if err != nil {
				return &[]Product{}, err
			}
		}

	}
	return &products, nil
}

func (p *Product) FindProductByID(db *gorm.DB, pid uint64) (*Product, error) {
	var err error
	err = db.Debug().Model(&Product{}).Where("id = ?", pid).Take(&p).Error
	if err != nil {
		return &Product{}, err
	}
	if p.ID != 0 {
		err = db.Debug().Model(&User{}).Where("id = ?", p.OwnerID).Take(&p.Owner).Error
		err = db.Debug().Model(&Auction{}).Where("id = ?", p.ID).Take(&p.Auction).Error
		if err != nil {
			return &Product{}, err
		}
	}
	return p, nil
}

func (p *Product) FindUserProducts(db *gorm.DB, uid uint32) (*[]Product, error) {

	var err error
	products := []Product{}
	err = db.Debug().Model(&Product{}).Where("owner_id = ?", uid).Limit(100).Order("created_at desc").Find(&products).Error
	if err != nil {
		return &[]Product{}, err
	}
	if len(products) > 0 {
		for i, _ := range products {
			err := db.Debug().Model(&User{}).Where("id = ?", products[i].OwnerID).Take(&products[i].Owner).Error
			if err != nil {
				return &[]Product{}, err
			}
		}
	}
	return &products, nil
}
