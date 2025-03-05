package response

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"github.com/go-playground/validator"
)

type Response struct{ //create a struct to represent a response
	Status string
	Error string
}

const( //constants for status so that we don't have to remember the exact string
	StatusError = "Error"
	StatusOK = "OK"
)

func WriteJson(w http.ResponseWriter, status int, data interface{}) error{
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(data)
}

func GeneralError(err error) Response{
	return Response{
		Status: "error",
		Error: err.Error(),
	} 
}
func ValidationError(errs validator.ValidationErrors) Response{
	var errMsg []string
	for _, err := range errs{
			switch err.Tag(){ //check the type of error
			case "required": //if the error is required 
				errMsg = append(errMsg, fmt.Sprintf("field %s is required", err.Field()))
			default: //default case
				errMsg = append(errMsg, fmt.Sprintf("field %s is invalid", err.Field()))
		}
	}
	return Response{ //return the response
		Status: StatusError,
		Error: strings.Join(errMsg, ", "), //join the error messages
	}
}