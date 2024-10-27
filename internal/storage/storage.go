package storage

type Storage interface {
	CreateStudent(name, email string, age int) (int64, error)
}
