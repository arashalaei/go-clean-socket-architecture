package model

type Role string

const (
	StudentRole Role = "student"
	TeacherRole Role = "teacher"
)

type Person struct {
	ID   uint   `gorm:"primaryKey;autoIncrement"`
	Name string `gorm:"type:varchar(255);not null"`
	Role Role   `gorm:"not null;check: role IN ('student', 'teacher')"`

	SchoolID *uint  `gorm:"index"`
	School   School `gorm:"foreignKey:SchoolID"`
}

func (Person) TableName() string {
	return "persons"
}
