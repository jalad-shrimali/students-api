package types

type Student struct{
	//create a struct to represent a student 
	Id int64
	Name string `validate:"required"`
	Age int	`validate:"required"`
	Email string	`validate:"required"`
}