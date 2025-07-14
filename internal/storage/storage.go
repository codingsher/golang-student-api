package storage

import "github.com/codingsher/golang-student-api/internal/types"

type Storage interface {
	CreateStudent(name string, email string, age int) (int64, error)
	GetStudentByID(id int64) (types.Student, error)
	GetStudents() ([]types.Student, error)
	UpdateStudentByID(id int64, name string, email string, age int) (int64, error)
	DeleteStudentByID(id int64) (int64, error)
}
