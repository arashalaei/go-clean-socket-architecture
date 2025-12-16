package mapper

import (
	"github.com/arashalaei/go-clean-socket-architecture/internal/domain/entity"
	"github.com/arashalaei/go-clean-socket-architecture/internal/repository/sqlite/model"
)

func SchoolToEntity(s *model.School) *entity.School {
	if s == nil {
		return nil
	}

	return &entity.School{
		Id:   s.ID,
		Name: s.Name,
	}
}
