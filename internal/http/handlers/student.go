package handlers

import (
	"encoding/json"
	
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/SHIVA-SINGHx/Go-Project/internal/types"
	"github.com/SHIVA-SINGHx/Go-Project/internal/utils/response"

)

var students = make(map[int]types.Student)

func CreateStudent() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.NotFound(w, r)
			return
		}

		var s types.Student
		if err := json.NewDecoder(r.Body).Decode(&s); err != nil || s.Id == 0 {
			response.WriteJson(w, http.StatusBadRequest, map[string]string{"error": "Invalid input"})
			return
		}

		students[s.Id] = s
		response.WriteJson(w, http.StatusCreated, map[string]string{"message": "Student created"})
	}
}

func StudentHandlerWithID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := strings.TrimPrefix(r.URL.Path, "/api/student/")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			response.WriteJson(w, http.StatusBadRequest, map[string]string{"error": "Invalid ID"})
			return
		}

		switch r.Method {
		case http.MethodGet:
			student, ok := students[id]
			if !ok {
				response.WriteJson(w, http.StatusNotFound, map[string]string{"error": "Student not found"})
				return
			}
			response.WriteJson(w, http.StatusOK, student)

		case http.MethodPut:
			var updated types.Student
			body, err := io.ReadAll(r.Body)
			if err != nil {
				response.WriteJson(w, http.StatusBadRequest, map[string]string{"error": "Failed to read body"})
				return
			}
			json.Unmarshal(body, &updated)
			updated.Id = id
			students[id] = updated
			response.WriteJson(w, http.StatusOK, map[string]string{"message": "Student updated"})

		case http.MethodDelete:
			delete(students, id)
			response.WriteJson(w, http.StatusOK, map[string]string{"message": "Student deleted"})

		default:
			http.NotFound(w, r)
		}
	}
}
