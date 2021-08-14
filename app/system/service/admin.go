package service

import (
	"lab-api/model"
)

type Admin struct {
	*Dependency
}

func NewAdmin(d Dependency) *Admin {
	return &Admin{&d}
}

func (x *Admin) FindByUsername(username string) (data model.Admin, err error) {
	if err = x.Db.
		Where("username = ?", username).
		Where("status = ?", true).
		First(&data).Error; err != nil {
		return
	}
	return
}