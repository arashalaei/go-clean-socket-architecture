package repository

import (
	"github.com/arashalaei/go-clean-socket-architecture/internal/domain/entity"
)

type PersonRepositroy interface {
	CreatePerson(person *entity.Person) uint
	GetPersonByID(personId uint) (*entity.Person, error)
}
