package model

import "gorm.io/gorm"

type Class struct {
	gorm.Model
	Name string `gorm:"not null"`

	SchoolID uint
	School   School `gorm:"foreignKey:SchoolID"`

	TeacherID uint
	Teacher   Person `gorm:"foreignKey:TeacherID"`

	Students []Person `gorm:"many2many:class_students"`
}
