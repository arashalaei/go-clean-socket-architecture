package store

import (
	"github.com/arashalaei/go-clean-socket-architecture/internal/domain/entity"
	"github.com/arashalaei/go-clean-socket-architecture/internal/repository/sqlite/model"
)

func (s *sqlit) CreatePerson(person *entity.Person) uint {
	student := model.Person{
		Name: person.Name,
		Role: model.StudentRole,
	}
	s.db.
		FirstOrCreate(
			&student,
			model.Person{
				Name: student.Name,
				Role: student.Role,
			})
	return student.ID
}
