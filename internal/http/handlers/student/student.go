package student

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/go-playground/validator"
	"github.com/jalad-shrimali/students-api/internal/storage"
	"github.com/jalad-shrimali/students-api/internal/types"
	"github.com/jalad-shrimali/students-api/internal/utils/response"
)

func New(storage storage.Storage) http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request){
		slog.Info("creating a student")

		var student types.Student

		err := json.NewDecoder(r.Body).Decode(&student) //decode the request body into student struct
		if errors.Is(err, io.EOF){
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}
		
		if err != nil{
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}

		// request for validation
		if err := validator.New().Struct(student); err != nil{
			validateErrs := err.(validator.ValidationErrors)
			response.WriteJson(w, http.StatusBadRequest, response.ValidationError(validateErrs))
			return
		}

		lastId, err := storage.CreateStudent(
			student.Name,
			student.Age,
			student.Email,
		)

		slog.Info("student created", slog.String("userId", fmt.Sprint(lastId)))
		if err!=nil{
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
			return
		}

		//validate the request
		response.WriteJson(w, http.StatusCreated, map[string]int64{"id": lastId})
	}
}