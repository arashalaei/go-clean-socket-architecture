package class

type ClassUsecases struct {
	CreateUseCase            *CreateClassUseCase
	ListUseCase              *ListClassesUseCase
	AddStudentToClassUseCase *AddStudentToClassUseCase
}

func NewClassUseCases(
	createUseCase *CreateClassUseCase,
	listUseCase *ListClassesUseCase,
	addStudentToClassUseCase *AddStudentToClassUseCase,
) *ClassUsecases {
	return &ClassUsecases{
		CreateUseCase:            createUseCase,
		ListUseCase:              listUseCase,
		AddStudentToClassUseCase: addStudentToClassUseCase,
	}
}
