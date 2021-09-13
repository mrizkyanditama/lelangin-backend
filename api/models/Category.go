package models

import (
	"gorm.io/gorm"
)

type Category struct {
	gorm.Model
	Name     string `gorm:"size:255;not null;unique" json:"name"`
	Products []Product
}

func (c *Category) FindAllCategories(db *gorm.DB) (*[]Category, error) {
	var err error
	categories := []Category{}
	err = db.Debug().Model(&Category{}).Find(&categories).Error
	if err != nil {
		return &[]Category{}, err
	}
	return &categories, nil
}
