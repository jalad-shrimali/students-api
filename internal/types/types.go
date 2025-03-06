package types

import "go.mongodb.org/mongo-driver/bson/primitive"

type Student struct{
	//create a struct to represent a student 
	Id int64 `json:"id"`
	Name string `validate:"required" json:"name"`
	Age int	`validate:"required" json:"age"`
	Email string	`validate:"required" json:"email"`
}
type MongoStudent struct {
	// create a struct to represent a student in MongoDB
	Id    primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name  string             `validate:"required" json:"name"`
	Age   int                `validate:"required" json:"age"`
	Email string             `validate:"required" json:"email"`
}