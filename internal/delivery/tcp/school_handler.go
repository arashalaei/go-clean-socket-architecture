package tcp

import (
	"context"
	"encoding/json"

	"github.com/arashalaei/go-clean-socket-architecture/internal/delivery/tcp/dto"
)

func (s *server) CreateSchoolHandler(
	ctx context.Context,
	payload json.RawMessage,
) (interface{}, error) {
	var req dto.CreateSchoolReq

	err := json.Unmarshal(payload, &req)
	if err != nil {
		return nil, err
	}

	schoolUsecases := s.schoolUsecases
	schoolId := schoolUsecases.CreateUseCase.Execute(req.Name)

	return schoolId, nil
}

func (s *server) ListSchoolsHandler(
	ctx context.Context,
	payload json.RawMessage,
) (interface{}, error) {
	schoolUsecases := s.schoolUsecases
	schools, err := schoolUsecases.ListUseCase.Execute()
	if err != nil {
		return nil, err
	}
	return schools, nil
}
