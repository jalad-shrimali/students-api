package student

import (
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"

	"github.com/go-playground/validator"
	"github.com/jalad-shrimali/students-api/internal/types"
	"github.com/jalad-shrimali/students-api/internal/utils/response"
)

func New() http.HandlerFunc{
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

		//validate the request
		response.WriteJson(w, http.StatusCreated, map[string]string{"success": "OK"})
	}
}