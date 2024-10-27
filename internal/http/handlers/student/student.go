package student

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/Sairam-04/students-api/internal/storage"
	"github.com/Sairam-04/students-api/internal/types"
	"github.com/Sairam-04/students-api/internal/utils/response"
	"github.com/go-playground/validator/v10"
)

func New(storage storage.Storage) http.HandlerFunc {
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

		lastId, err := storage.CreateStudent(
			student.Name,
			student.Email,
			student.Age,
		)
		slog.Info("creating a student")
		slog.Info("user created successfully", slog.String("userId", fmt.Sprint(lastId)))
		if err != nil {
			fmt.Println(err)
			response.WriteJson(w, http.StatusInternalServerError, err)
			return
		}
		response.WriteJson(w, http.StatusCreated, map[string]int64{"id": lastId})
	}
}

func GetByID(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// get id
		id := r.PathValue("id")
		slog.Info("getting a student", slog.String("id", id))
		intId, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}
		student, err := storage.GetStudentByID(intId)
		if err != nil {
			slog.Error("error getting user", slog.String("id", id))
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
			return
		}

		response.WriteJson(w, http.StatusOK, student)
	}
}

func GetList(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slog.Info("getting all students")
		students, err := storage.GetStudents()
		if err != nil {
			response.WriteJson(w, http.StatusInternalServerError, err)
			return
		}
		response.WriteJson(w, http.StatusOK, students)
	}
}

func UpdateStudent(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		slog.Info("updating a student", slog.String("id", id))
		intId, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}
		var updatedStudent types.Student
		err = json.NewDecoder(r.Body).Decode(&updatedStudent)
		if err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}

		student, err := storage.UpdateStudentByID(intId, updatedStudent)
		if err != nil {
			slog.Error("error updating user", slog.String("id", id))
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
			return
		}

		response.WriteJson(w, http.StatusOK, student)
	}
}

func DeleteStudent(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		slog.Info("deleting student", slog.String("id", id))
		intID, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}
		err = storage.DeleteByID(intID)
		if err != nil {
			if strings.Contains(err.Error(), "no student found") {
				response.WriteJson(w, http.StatusNotFound, response.GeneralError(err))
				return
			}
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
			return
		}
		response.WriteJson(w, http.StatusOK, map[string]string{"message": fmt.Sprintf("successfully deleted student with id: %d", intID)})
	}
}
