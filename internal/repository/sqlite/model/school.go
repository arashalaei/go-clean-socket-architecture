package model

import "gorm.io/gorm"

type School struct {
	gorm.Model
	Name    string  `gorm:"type:varchar(255);not null; unique"`
	Classes []Class `gorm:"foreignKey:SchoolID"`
}
