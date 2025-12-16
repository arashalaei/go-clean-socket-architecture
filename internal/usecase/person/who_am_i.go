package person

import (
	"github.com/arashalaei/go-clean-socket-architecture/internal/domain/entity"
	"github.com/arashalaei/go-clean-socket-architecture/internal/domain/repository"
)

type WhoAmIUseCase struct {
	personRepo repository.PersonRepositroy
}

func NewWhoAmIUseCase(
	personRepo repository.PersonRepositroy,
) *WhoAmIUseCase {
	return &WhoAmIUseCase{
		personRepo: personRepo,
	}
}

func (uc *WhoAmIUseCase) Execute(personId uint) (*entity.Person, error) {
	return uc.personRepo.GetPersonByID(personId)
}

