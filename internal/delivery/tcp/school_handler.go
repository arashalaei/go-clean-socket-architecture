package tcp

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/arashalaei/go-clean-socket-architecture/internal/delivery/tcp/dto"
)

func (s *server) CreateSchoolHandler(ctx context.Context, payload json.RawMessage) (interface{}, error) {
	var req dto.School

	err := json.Unmarshal(payload, &req)
	if err != nil {
		return nil, err
	}

	fmt.Printf("%+v", req)
	return nil, nil
}
