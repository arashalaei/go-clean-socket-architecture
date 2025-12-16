package store

import (
	"errors"
	"fmt"

	"github.com/arashalaei/go-clean-socket-architecture/internal/domain/entity"
	"github.com/arashalaei/go-clean-socket-architecture/internal/repository/mapper"
	"github.com/arashalaei/go-clean-socket-architecture/internal/repository/sqlite/model"
	"gorm.io/gorm"
)

func (s *sqlit) CreateClass(
	name string,
	schoolId, teacherId uint,
) uint {
	var class = model.Class{
		Name:      name,
		TeacherID: teacherId,
		SchoolID:  schoolId,
	}

	s.db.
		FirstOrCreate(
			&class,
			model.Class{
				Name:      class.Name,
				SchoolID:  class.SchoolID,
				TeacherID: class.TeacherID,
			})
	return class.ID
}

func (s *sqlit) GetClassByID(id uint) (*entity.Class, error) {
	var class model.Class
	err := s.db.
		Preload("School").
		Preload("Teacher").
		Preload("Students").
		First(&class, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("class not found: %w", err)
		}
		return nil, fmt.Errorf("failed to get class by id: %w", err)
	}
	return mapper.ClassToEntity(&class), nil
}

func (s *sqlit) GetAllClasses() (*[]entity.Class, error) {
	var classes []model.Class
	err := s.db.
		Preload("School").
		Preload("Teacher").
		Preload("Students").
		Find(&classes).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get classes: %w", err)
	}
	return mapper.ClassesToEntities(classes), nil
}

func (s *sqlit) AddStudentToClass(classId, studentId uint) error {
	var class model.Class
	if err := s.db.First(&class, classId).Error; err != nil {
		return fmt.Errorf("class not found: %w", err)
	}

	var student model.Person
	if err := s.db.First(&student, studentId).Error; err != nil {
		return fmt.Errorf("student not found: %w", err)
	}

	if err := s.db.Model(&class).Association("Students").Append(&student); err != nil {
		return fmt.Errorf("failed to add student to class: %w", err)
	}

	return nil
}
