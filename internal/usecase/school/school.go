package school

type SchoolUsecases struct {
	CreateUseCase *CreateSchoolUseCase
}

func NewSchoolUseCases(
	createUseCase *CreateSchoolUseCase,
) *SchoolUsecases {
	return &SchoolUsecases{
		CreateUseCase: createUseCase,
	}
}
