package tcp

import (
	"context"
	"encoding/json"

	"github.com/arashalaei/go-clean-socket-architecture/internal/delivery/tcp/dto"
	"github.com/arashalaei/go-clean-socket-architecture/internal/domain/entity"
)

func (s *server) CreatePersonHandler(ctx context.Context, payload json.RawMessage) (interface{}, error) {
	var req dto.CreatePersonReq

	err := json.Unmarshal(payload, &req)
	if err != nil {
		return nil, err
	}

	personUsecases := s.personUsecases
	personId := personUsecases.CreateUseCase.Execute(entity.Person{
		Name:   req.Name,
		Role:   entity.Role(req.Role),
		School: entity.School{Id: req.SchoolId},
	})

	return personId, nil
}

func (s *server) ListPersonsHandler(ctx context.Context, payload json.RawMessage) (interface{}, error) {
	personUsecases := s.personUsecases
	persons, err := personUsecases.ListUseCase.Execute()
	if err != nil {
		return nil, err
	}
	return persons, nil
}

func (s *server) WhoAmIHandler(ctx context.Context, payload json.RawMessage) (interface{}, error) {
	var req dto.WhoAmIReq

	err := json.Unmarshal(payload, &req)
	if err != nil {
		return nil, err
	}

	personUsecases := s.personUsecases
	person, err := personUsecases.WhoAmIUseCase.Execute(req.PersonId)
	if err != nil {
		return nil, err
	}

	return person, nil
}
