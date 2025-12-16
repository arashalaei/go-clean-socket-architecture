package person

// Application Layer (Application Business Rules)

import (
	"github.com/arashalaei/go-clean-socket-architecture/internal/domain/entity"
	"github.com/arashalaei/go-clean-socket-architecture/internal/domain/repository"
)

type CreatePersonUseCase struct {
	personRepo repository.PersonRepositroy
}

func NewCreateSchoolUseCase(
	personRepo repository.PersonRepositroy,
) *CreatePersonUseCase {
	return &CreatePersonUseCase{
		personRepo: personRepo,
	}
}

func (uc *CreatePersonUseCase) Execute(p entity.Person) uint {
	return uc.personRepo.CreatePerson(&p)
}
