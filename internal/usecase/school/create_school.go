package school

import "github.com/arashalaei/go-clean-socket-architecture/internal/domain/repository"

type CreateSchoolUseCase struct {
	schoolRepo repository.SchoolRepository
}

func NewCreateSchoolUseCase(
	schoolRepo repository.SchoolRepository,
) *CreateSchoolUseCase {
	return &CreateSchoolUseCase{
		schoolRepo: schoolRepo,
	}
}

func (uc *CreateSchoolUseCase) Execute(schoolName string) uint {
	return uc.schoolRepo.CreateSchool(schoolName)
}
