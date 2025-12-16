package store

import "github.com/arashalaei/go-clean-socket-architecture/internal/repository/sqlite/model"

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
