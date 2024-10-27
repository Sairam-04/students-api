package storage

import "github.com/Sairam-04/students-api/internal/types"

type Storage interface {
	CreateStudent(name, email string, age int) (int64, error)
	GetStudentByID(id int64) (types.Student, error)
}
