package domain

import (
	"gorm.io/gorm"
	"sapphire-server/internal/dao"
)

type Sam struct {
	gorm.Model
	Onnxname string `gorm:"column:onnxname"`
}

func NewSam() *Sam {
	return &Sam{}
}

//func (f *Sam) registerSam(onnxname string) {
//	existedSAM := f.loadSAM(map[string]interface{}{"name": register.Name})
//}

func (f *Sam) loadSAM(param map[string]interface{}) *Sam {
	sam, err := dao.First[Sam]("onnx_name = ?", param["onnx_name"])
	if err != nil {
		return nil
	}
	return sam
}

func (f *Sam) loadSAMByID(id int) *Sam {
	sam, err := dao.First[Sam]("id = ?", id)
	if err != nil {
		return nil
	}
	return sam
}
