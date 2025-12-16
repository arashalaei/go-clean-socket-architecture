package person

type PersonUsecases struct {
	CreateUseCase  *CreatePersonUseCase
	ListUseCase    *ListPersonsUseCase
	WhoAmIUseCase  *WhoAmIUseCase
	EnrollUseCase  *EnrollInSchoolStudentUseCase
}

func NewPersonUseCases(
	createUseCase *CreatePersonUseCase,
	listUseCase *ListPersonsUseCase,
	whoAmIUseCase *WhoAmIUseCase,
	enrollUseCase *EnrollInSchoolStudentUseCase,
) *PersonUsecases {
	return &PersonUsecases{
		CreateUseCase: createUseCase,
		ListUseCase:   listUseCase,
		WhoAmIUseCase: whoAmIUseCase,
		EnrollUseCase: enrollUseCase,
	}
}

