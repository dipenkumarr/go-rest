package student

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/dipenkumarr/go-rest/internal/storage"
	"github.com/dipenkumarr/go-rest/internal/types"
	"github.com/dipenkumarr/go-rest/internal/utils/response"
	"github.com/go-playground/validator/v10"
)


func New(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slog.Info("Creating a student")

		var student types.Student

		err := json.NewDecoder(r.Body).Decode(&student)
		if errors.Is(err, io.EOF) {
			response.WriteJSON(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}

		if err != nil {
			response.WriteJSON(w, http.StatusBadRequest, response.GeneralError(err))
		}

		// validation
		if err := validator.New().Struct(student); err != nil {
			validateErrs := err.(validator.ValidationErrors)
			response.WriteJSON(w, http.StatusBadRequest, response.ValidationError(validateErrs))
			return
		}

		lastId, err := storage.CreateStudent(student.Name, student.Email, student.Age)

		slog.Info("USER CREATED SUCCESSFULLY", slog.String("lastId", fmt.Sprint(lastId)))

		if err != nil {
			response.WriteJSON(w, http.StatusInternalServerError, err)
			return
		}
		
		response.WriteJSON(w, http.StatusCreated, map[string]int64 {"id": lastId})
	}
}

