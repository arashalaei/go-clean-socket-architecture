package school

type SchoolUsecases struct {
	CreateUseCase *CreateSchoolUseCase
	ListUseCase   *ListSchoolsUseCase
}

func NewSchoolUseCases(
	createUseCase *CreateSchoolUseCase,
	listUseCase *ListSchoolsUseCase,
) *SchoolUsecases {
	return &SchoolUsecases{
		CreateUseCase: createUseCase,
		ListUseCase:   listUseCase,
	}
}
