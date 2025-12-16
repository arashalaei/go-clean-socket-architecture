package class

import "github.com/arashalaei/go-clean-socket-architecture/internal/domain/repository"

type CreateClassUseCase struct {
	classRepo repository.ClassRepository
}

func NewCreateClassUseCase(
	classRepo repository.ClassRepository,
) *CreateClassUseCase {
	return &CreateClassUseCase{
		classRepo: classRepo,
	}
}

func (uc *CreateClassUseCase) Execute(name string, schoolId, teacherId uint) uint {
	return uc.classRepo.CreateClass(name, schoolId, teacherId)
}

