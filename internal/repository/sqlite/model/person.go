package model

import "gorm.io/gorm"

type Role string

const (
	StudentRole Role = "student"
	TeacherRole Role = "teacher"
)

type Person struct {
	gorm.Model
	Name string `gorm:"type:varchar(255);not null"`
	Role Role   `gorm:"not null;check: role IN ('student', 'teacher')"`

	SchoolID *uint  `gorm:"index"`
	School   School `gorm:"foreignKey:SchoolID"`
}
