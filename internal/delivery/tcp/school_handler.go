package tcp

import (
	"context"
	"encoding/json"

	"github.com/arashalaei/go-clean-socket-architecture/internal/delivery/tcp/dto"
)

func (s *server) CreateSchoolHandler(ctx context.Context, payload json.RawMessage) (interface{}, error) {
	var req dto.CreateSchoolReq

	err := json.Unmarshal(payload, &req)
	if err != nil {
		return nil, err
	}

	schoolUsecases := s.schoolUsecases
	schoolId := schoolUsecases.CreateUseCase.Execute(req.Name)

	return schoolId, nil
}
