package storage

import "github.com/Sairam-04/students-api/internal/types"

type Storage interface {
	CreateStudent(name, email string, age int) (int64, error)
	GetStudentByID(id int64) (types.Student, error)
	GetStudents() ([]types.Student, error)
	UpdateStudentByID(id int64, updatedStudent types.Student) (types.Student, error)
	DeleteByID(id int64) error
}
