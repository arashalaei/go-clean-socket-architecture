package store

import (
	"errors"
	"fmt"

	"github.com/arashalaei/go-clean-socket-architecture/internal/domain/entity"
	"github.com/arashalaei/go-clean-socket-architecture/internal/repository/mapper"
	"github.com/arashalaei/go-clean-socket-architecture/internal/repository/sqlite/model"
	"gorm.io/gorm"
)

func (s *sqlit) CreateSchool(name string) uint {
	school := &model.School{}
	s.db.
		Where(model.School{Name: name}).
		FirstOrCreate(school)
	return school.ID
}

func (s *sqlit) GetSchoolByID(schoolId uint) (*entity.School, error) {
	var school model.School
	if err := s.db.First(&school, schoolId).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("school not found: %w", err)
		}
		return nil, fmt.Errorf("failed to get school by id: %w", err)
	}
	return mapper.SchoolToEntity(&school), nil
}

func (s *sqlit) GetSchoolByName(schoolName string) (*entity.School, error) {
	var school model.School

	if err := s.db.Where("name = ?", school).First(&school).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("school not found: %w", err)
		}
		return nil, fmt.Errorf("failed to get school by id: %w", err)
	}
	return mapper.SchoolToEntity(&school), nil
}
