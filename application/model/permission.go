package model

import "github.com/kainonly/gin-extra/datatype"

type Permission struct {
	ID         uint64
	Key        string
	Name       datatype.JSONObject `gorm:"type:json"`
	Note       string
	Status     bool
	CreateTime uint64 `gorm:"autoCreateTime"`
	UpdateTime uint64 `gorm:"autoUpdateTime"`
}
