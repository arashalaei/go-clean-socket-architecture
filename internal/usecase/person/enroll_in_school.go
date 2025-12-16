package person

import (
	"github.com/arashalaei/go-clean-socket-architecture/internal/domain/entity"
	"github.com/arashalaei/go-clean-socket-architecture/internal/domain/repository"
)

type EnrollInSchoolStudentUseCase struct {
	personRepo repository.PersonRepositroy
	schoolReop repository.SchoolRepository
}

func NewEnrollInSchoolStudentUseCase(
	personRepo repository.PersonRepositroy,
	schoolReop repository.SchoolRepository,
) *EnrollInSchoolStudentUseCase {
	return &EnrollInSchoolStudentUseCase{
		personRepo: personRepo,
		schoolReop: schoolReop,
	}
}

func (uc *EnrollInSchoolStudentUseCase) Execute(
	studentName,
	schoolName string,
) error {
	school, err := uc.schoolReop.GetSchoolByName(schoolName)
	if err != nil {
		return err
	}

	uc.personRepo.CreatePerson(
		&entity.Person{
			Name:   studentName,
			Role:   entity.StudentRole,
			School: *school,
		})

	return nil
}
