package repository

import "github.com/arashalaei/go-clean-socket-architecture/internal/domain/entity"

type SchoolRepository interface {
	CreateSchool(name string) uint
	GetSchoolByID(id uint) (*entity.School, error)
	GetSchoolByName(schoolName string) (*entity.School, error)
}
