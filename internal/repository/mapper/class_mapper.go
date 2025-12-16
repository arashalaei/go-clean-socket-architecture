package mapper

import (
	"github.com/arashalaei/go-clean-socket-architecture/internal/domain/entity"
	"github.com/arashalaei/go-clean-socket-architecture/internal/repository/sqlite/model"
)

func ClassToEntity(c *model.Class) *entity.Class {
	if c == nil {
		return nil
	}

	var students []entity.Person
	for _, s := range c.Students {
		students = append(students, *PersonToEntity(&s))
	}

	return &entity.Class{
		Id:       c.ID,
		Name:     c.Name,
		SchoolId: c.SchoolID,
		Teacher:  *PersonToEntity(&c.Teacher),
		Students: students,
	}
}

func ClassesToEntities(classes []model.Class) *[]entity.Class {
	var classToEntities []entity.Class

	for _, c := range classes {
		classToEntities = append(classToEntities, *ClassToEntity(&c))
	}

	return &classToEntities
}
