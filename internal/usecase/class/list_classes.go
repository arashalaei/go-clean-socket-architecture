package class

import (
	"github.com/arashalaei/go-clean-socket-architecture/internal/domain/entity"
	"github.com/arashalaei/go-clean-socket-architecture/internal/domain/repository"
)

type ListClassesUseCase struct {
	classRepo repository.ClassRepository
}

func NewListClassesUseCase(
	classRepo repository.ClassRepository,
) *ListClassesUseCase {
	return &ListClassesUseCase{
		classRepo: classRepo,
	}
}

func (uc *ListClassesUseCase) Execute() (*[]entity.Class, error) {
	return uc.classRepo.GetAllClasses()
}

