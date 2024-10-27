package student

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/Sairam-04/students-api/internal/types"
	"github.com/Sairam-04/students-api/internal/utils/response"
	"github.com/go-playground/validator/v10"
)

func New() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Serializing the data
		var student types.Student
		err := json.NewDecoder(r.Body).Decode(&student) // decode in student variable
		if errors.Is(err, io.EOF) {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("Empty body")))
			return
		}

		if err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}

		// validating request
		if err := validator.New().Struct(student); err != nil {
			// typecast err to validator

			validateErrs := err.(validator.ValidationErrors)
			response.WriteJson(w, http.StatusBadRequest, response.ValidationError(validateErrs))
			return
		}
		slog.Info("creating a student")
		response.WriteJson(w, http.StatusCreated, map[string]string{"success": "Ok"})
	}
}
