package models

import (
	"gorm.io/gorm"
)

type Tag struct {
	gorm.Model
	Name     string     `gorm:"size:255;not null;unique" json:"name"`
	Products []*Product `gorm:"many2many:product_tag;"`
}

func (t *Tag) FindAllTags(db *gorm.DB) (*[]Tag, error) {
	var err error
	tags := []Tag{}
	err = db.Debug().Model(&Tag{}).Find(&tags).Error
	if err != nil {
		return &[]Tag{}, err
	}
	return &tags, nil
}
