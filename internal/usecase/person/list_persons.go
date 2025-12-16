package person

import (
	"github.com/arashalaei/go-clean-socket-architecture/internal/domain/entity"
	"github.com/arashalaei/go-clean-socket-architecture/internal/domain/repository"
)

type ListPersonsUseCase struct {
	personRepo repository.PersonRepositroy
}

func NewListPersonsUseCase(
	personRepo repository.PersonRepositroy,
) *ListPersonsUseCase {
	return &ListPersonsUseCase{
		personRepo: personRepo,
	}
}

func (uc *ListPersonsUseCase) Execute() (*[]entity.Person, error) {
	return uc.personRepo.GetAllPersons()
}

