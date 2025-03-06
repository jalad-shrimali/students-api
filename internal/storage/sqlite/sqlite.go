package sqlite

import (
	"database/sql"
	"fmt"

	// "github.com/jalad-shrimali/students-api/internal/http/handlers/student"
	"github.com/jalad-shrimali/students-api/internal/types"

	"github.com/jalad-shrimali/students-api/internal/config"
	// import the sqlite3 driver
	_ "github.com/mattn/go-sqlite3"
)
type Sqlite struct{
	Db *sql.DB 
}

func New(cfg *config.Config)(*Sqlite, error){
	db, err := sql.Open("sqlite3", cfg.StoragePath)
	if err != nil{
		return nil, err
	}
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS students(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT,
		age INTEGER,
		email TEXT 
	)`)
	if err != nil{
		return nil, err
	}
	return &Sqlite{
		Db: db,
	}, nil
}

// now to attach the interface to the storage type we just need to use the method that we have used in interface in our struct

func(s *Sqlite) CreateStudent(name string, age int, email string) (int64, error){ 
	//this is how our struct is implementing the interface
	stmt, err := s.Db.Prepare("INSERT INTO students(name, age, email) VALUES(?,?,?)")
	if err != nil{
		return 0, err
	}
	defer stmt.Close()
	result, err := stmt.Exec(name, age, email)
	if err != nil{
		return 0, err
	}
	lastId, err := result.LastInsertId()
	if err != nil{
		return 0, err
	}
	return lastId, nil //return lastId and nil
}

func (s *Sqlite) GetStudentById(id int64) (types.Student, error){
	stmt, err := s.Db.Prepare("SELECT * FROM students WHERE id = ? LIMIT 1")
	if err != nil{
		return types.Student{}, err
	}
	defer stmt.Close()
	var student types.Student
	err = stmt.QueryRow(id).Scan(&student.Id, &student.Name, &student.Age, &student.Email)
	if err != nil{
		if err == sql.ErrNoRows{
			return types.Student{}, fmt.Errorf("no student found with id %d", id)
		}
		return types.Student{}, err
	}
	return student, nil
}

func (s *Sqlite) GetAllStudents() ([]types.Student, error){
	stmt, err := s.Db.Prepare("SELECT * FROM students")
	if err != nil{
		return nil, err
	}
	defer stmt.Close()
	rows, err := stmt.Query()
	if err != nil{
		return nil, err
	}
	defer rows.Close()
	var students []types.Student
	for rows.Next(){
		var student types.Student
		err = rows.Scan(&student.Id, &student.Name, &student.Age, &student.Email)
		if err != nil{
			return nil, err
		}
		students = append(students, student)
	}
	return students, nil
}