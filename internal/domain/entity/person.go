package entity

// Enterprise Business Rules

type Role string

const (
	StudentRole Role = "student"
	TeacherRole Role = "Teacher"
)

type Person struct {
	Id      uint
	Name    string
	Role    Role
	Classes []uint
}
