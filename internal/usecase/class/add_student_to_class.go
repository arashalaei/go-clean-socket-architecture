package class

import "github.com/arashalaei/go-clean-socket-architecture/internal/domain/repository"

type AddStudentToClassUseCase struct {
	classRepo repository.ClassRepository
}

func NewAddStudentToClassUseCase(
	classRepo repository.ClassRepository,
) *AddStudentToClassUseCase {
	return &AddStudentToClassUseCase{
		classRepo: classRepo,
	}
}

func (uc *AddStudentToClassUseCase) Execute(classId, studentId uint) error {
	return uc.classRepo.AddStudentToClass(classId, studentId)
}
