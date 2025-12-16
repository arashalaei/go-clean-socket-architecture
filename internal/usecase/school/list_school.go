package school

import (
	"github.com/arashalaei/go-clean-socket-architecture/internal/domain/entity"
	"github.com/arashalaei/go-clean-socket-architecture/internal/domain/repository"
)

type ListSchoolsUseCase struct {
	schoolRepo repository.SchoolRepository
}

func NewListSchoolsUseCase(
	schoolRepo repository.SchoolRepository,
) *ListSchoolsUseCase {
	return &ListSchoolsUseCase{
		schoolRepo: schoolRepo,
	}
}

func (uc *ListSchoolsUseCase) Execute() (*[]entity.School, error) {
	return uc.schoolRepo.GetAllSchools()
}
