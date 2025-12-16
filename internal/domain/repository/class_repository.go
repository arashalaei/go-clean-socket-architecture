package repository

import "github.com/arashalaei/go-clean-socket-architecture/internal/domain/entity"

type ClassRepository interface {
	CreateClass(name string, schoolId, teacherId uint) uint
	GetClassByID(id uint) (*entity.Class, error)
	GetAllClasses() (*[]entity.Class, error)
	AddStudentToClass(classId, studentId uint) error
}
