package types

type Student struct{
	//create a struct to represent a student 
	Id int64 `json:"id"`
	Name string `validate:"required" json:"name"`
	Age int	`validate:"required" json:"age"`
	Email string	`validate:"required" json:"email"`
}