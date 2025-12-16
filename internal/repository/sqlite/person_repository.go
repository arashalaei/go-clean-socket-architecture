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
			return nil, fmt.Errorf("person not found: %w", err)
		}
		return nil, fmt.Errorf("failed to get person by id: %w", err)
	}
	return mapper.PersonToEntity(&person), nil

}

func (s *sqlit) GetAllPersons() (*[]entity.Person, error) {
	var persons []model.Person
	err := s.db.
		Preload("School").
		Find(&persons).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get persons: %w", err)
	}

	var personToEntities []entity.Person
	for _, p := range persons {
		personToEntities = append(personToEntities, *mapper.PersonToEntity(&p))
	}

	return &personToEntities, nil
}
