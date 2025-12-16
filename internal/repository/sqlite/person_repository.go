package store

import (
	"errors"
	"fmt"

	"github.com/arashalaei/go-clean-socket-architecture/internal/domain/entity"
	"github.com/arashalaei/go-clean-socket-architecture/internal/repository/mapper"
	"github.com/arashalaei/go-clean-socket-architecture/internal/repository/sqlite/model"
	"gorm.io/gorm"
)

func (s *sqlit) CreatePerson(person *entity.Person) uint {
	p := model.Person{
		Name:     person.Name,
		Role:     model.Role(person.Role),
		SchoolID: &person.School.Id,
	}
	s.db.
		Create(&p)
	return p.ID
}

func (s *sqlit) GetPersonByID(personId uint) (*entity.Person, error) {
	var person model.Person
	err := s.db.
		Preload("School").
		First(&person, personId).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("school not found: %w", err)
		}
		return nil, fmt.Errorf("failed to get school by id: %w", err)
	}
	return mapper.PersonToEntity(&person), nil

}
