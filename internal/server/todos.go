package server

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
	"todo/todos"

	"github.com/labstack/echo/v4"
)

func (s *Server) Index(c echo.Context) error {

	rows, err := s.db.Query("SELECT * FROM  todos")

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	var (
		name       string
		completed  bool
		idTodo     uint
		created_at time.Time
		updated_at time.Time
	)

	data := []todos.Todo{}

	for rows.Next() {
		err := rows.Scan(&name, &completed, &idTodo, &created_at, &updated_at)

		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}

		todo := todos.Todo{ID: idTodo, Title: name, Completed: completed, CreatedAt: created_at, UpdatedAt: updated_at}
		data = append(data, todo)
	}

	resp := map[string][]todos.Todo{
		"message": data,
	}
	defer rows.Close()

	return c.JSON(http.StatusOK, resp)
}

func (s *Server) Create(c echo.Context) error {
	var todo todos.Todo

	err := json.NewDecoder(c.Request().Body).Decode(&todo)

	if err != nil {
		log.Fatal("decode error")
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
	}

	creationError := todo.Validate()

	if creationError != nil {
		return c.JSON(http.StatusUnprocessableEntity, map[string]string{"error": creationError.Error()})
	}

	_, insertError := s.db.Query(`INSERT INTO todos (title, completed) VALUES ($1, $2)`, todo.Title, todo.Completed)

	if insertError != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": insertError.Error()})
	}

	return c.JSON(http.StatusCreated, map[string]todos.Todo{"success": todo})
}

// func (s *Server) HealthHandler(c echo.Context) error {
// 	return c.JSON(http.StatusOK, s.db.Health())
// }
