package dto

type Class struct {
	Id       uint     `json:"id,omitempty"`
	Name     string   `json:"name,omitempty"`
	SchoolId uint     `json:"school_id,omitempty"`
	Teacher  Person   `json:"teacher,omitempty"`
	Students []Person `json:"students,omitempty"`
}

type AddStudentToClassReq struct {
	StudentId uint `json:"student_id,omitempty"`
	ClassId   uint `json:"class_id,omitempty"`
}
