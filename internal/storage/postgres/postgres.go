package postgres

import (
	"database/sql"
	"fmt"

	"github.com/Sairam-04/students-api/internal/config"
	"github.com/Sairam-04/students-api/internal/types"
	_ "github.com/lib/pq" // PostgreSQL driver
)

type Postgres struct {
	Db *sql.DB
}

func New(cfg *config.Config) (*Postgres, error) {
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS student(
		id SERIAL PRIMARY KEY,
		name TEXT NOT NULL,
		email TEXT NOT NULL UNIQUE,
		age INTEGER NOT NULL
	)`)

	if err != nil {
		return nil, err
	}

	return &Postgres{
		Db: db,
	}, nil
}

func (p *Postgres) CreateStudent(name, email string, age int) (int64, error) {
	stmt, err := p.Db.Prepare("INSERT INTO student (name, email, age) VALUES ($1, $2, $3) RETURNING id")
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	var lastId int64
	err = stmt.QueryRow(name, email, age).Scan(&lastId)
	if err != nil {
		return 0, err
	}

	return lastId, nil
}

func (p *Postgres) GetStudentByID(id int64) (types.Student, error) {
	stmt, err := p.Db.Prepare("SELECT id, name, email, age FROM student WHERE id = $1 LIMIT 1")
	if err != nil {
		return types.Student{}, err
	}
	defer stmt.Close()

	var student types.Student
	err = stmt.QueryRow(id).Scan(&student.Id, &student.Name, &student.Email, &student.Age)
	if err != nil {
		if err == sql.ErrNoRows {
			return types.Student{}, fmt.Errorf("no student found with id %d", id)
		}
		return types.Student{}, fmt.Errorf("query error: %w", err)
	}
	return student, nil
}

func (p *Postgres) GetStudents() ([]types.Student, error) {
	rows, err := p.Db.Query("SELECT id, name, email, age FROM student")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var students []types.Student
	for rows.Next() {
		var student types.Student
		err := rows.Scan(&student.Id, &student.Name, &student.Email, &student.Age)
		if err != nil {
			return nil, err
		}
		students = append(students, student)
	}
	return students, nil
}

func (p *Postgres) UpdateStudentByID(id int64, updatedStudent types.Student) (types.Student, error) {
	stmt, err := p.Db.Prepare("UPDATE student SET name = $1, email = $2, age = $3 WHERE id = $4")
	if err != nil {
		return types.Student{}, fmt.Errorf("failed to prepare update statement: %w", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(updatedStudent.Name, updatedStudent.Email, updatedStudent.Age, id)
	if err != nil {
		return types.Student{}, fmt.Errorf("failed to update student: %w", err)
	}
	return p.GetStudentByID(id)
}

func (p *Postgres) DeleteByID(id int64) error {
	stmt, err := p.Db.Prepare("DELETE FROM student WHERE id = $1")
	if err != nil {
		return fmt.Errorf("failed to prepare delete statement: %w", err)
	}
	defer stmt.Close()

	result, err := stmt.Exec(id)
	if err != nil {
		return fmt.Errorf("failed to delete student: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to retrieve affected rows: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("no student found with id: %d", id)
	}
	return nil
}
