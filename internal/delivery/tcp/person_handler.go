package tcp

import (
	"context"
	"encoding/json"
)

func (s *server) CreatePersonHandler(ctx context.Context, payload json.RawMessage) (interface{}, error) {
	panic("unimplemented")
}

func (s *server) WhoAmIHandler(ctx context.Context, payload json.RawMessage) (interface{}, error) {
	panic("unimplemented")
}
