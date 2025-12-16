package tcp

import (
	"context"
	"encoding/json"

	"github.com/arashalaei/go-clean-socket-architecture/internal/delivery/tcp/dto"
)

func (s *server) CreateClassHandler(ctx context.Context, payload json.RawMessage) (interface{}, error) {
	var req dto.CreateClassReq

	err := json.Unmarshal(payload, &req)
	if err != nil {
		return nil, err
	}

	classUsecases := s.classUsecases
	classId := classUsecases.CreateUseCase.Execute(req.Name, req.SchoolId, req.TeacherId)

	return classId, nil
}

func (s *server) ListClassesHandler(ctx context.Context, payload json.RawMessage) (interface{}, error) {
	classUsecases := s.classUsecases
	classes, err := classUsecases.ListUseCase.Execute()
	if err != nil {
		return nil, err
	}
	return classes, nil
}

func (s *server) AddStudentToClassHandler(ctx context.Context, payload json.RawMessage) (interface{}, error) {
	var req dto.AddStudentToClassReq

	err := json.Unmarshal(payload, &req)
	if err != nil {
		return nil, err
	}

	classUsecases := s.classUsecases
	err = classUsecases.AddStudentToClassUseCase.Execute(req.ClassId, req.StudentId)
	if err != nil {
		return nil, err
	}

	return "student added to class successfully", nil
}
