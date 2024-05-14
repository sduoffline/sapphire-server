package domain

import (
	"gorm.io/gorm"
)

type Image struct {
	gorm.Model
	Filename  string `gorm:"column:filename"`
	Extension string `gorm:"column:extension"`
	Url       string `gorm:"column:url"`
}

func NewImage() *Image {
	return &Image{}
}

func (f *Image) UploadImage() {

}
