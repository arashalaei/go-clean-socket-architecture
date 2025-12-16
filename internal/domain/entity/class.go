package entity

type Class struct {
	Id       uint
	Name     string
	SchoolId uint
	Teacher  Person
	Students []Person
}
