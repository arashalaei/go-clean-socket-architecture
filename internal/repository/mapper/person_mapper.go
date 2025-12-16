package mapper

import (
	"github.com/arashalaei/go-clean-socket-architecture/internal/domain/entity"
	"github.com/arashalaei/go-clean-socket-architecture/internal/repository/sqlite/model"
)

func PersonToEntity(p *model.Person) *entity.Person {
	if p == nil {
		return nil
	}

	return &entity.Person{
		Id:     p.ID,
		Name:   p.Name,
		Role:   entity.Role(p.Role),
		School: *SchoolToEntity(&p.School),
	}
}
