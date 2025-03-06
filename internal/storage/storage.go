package storage
import "github.com/jalad-shrimali/students-api/internal/types"

type Storage interface{
	// we are creating this interface to make sure that all the storage types implement these methods
	// this is a good practice to make sure that all the storage types have the same methods
	CreateStudent(name string, age int, email string) (int64, error)
	GetStudentById(id int64) (types.Student, error)
	GetAllStudents() ([]types.Student, error)
}