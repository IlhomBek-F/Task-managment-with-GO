package server

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"time"
	"todo/todos"

	"github.com/labstack/echo/v4"
)

type response struct {
	Status int                    `json:"status"`
	Res    map[string]interface{} `json:"res"`
}

func (s *Server) Index(c echo.Context) error {

	query := c.QueryParam("q")
	rows, err := s.db.Query("SELECT * FROM  todos WHERE title ILIKE $1", "%"+query+"%")

	if err != nil {
		return c.JSON(http.StatusInternalServerError, response{
			Status: http.StatusInternalServerError,
			Res:    map[string]interface{}{"error": err.Error()},
		})
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

	resp := response{
		Status: http.StatusOK,
		Res:    map[string]interface{}{"data": data},
	}

	defer rows.Close()

	return c.JSON(http.StatusOK, resp)
}

func (s *Server) Create(c echo.Context) error {
	var todo todos.Todo

	err := json.NewDecoder(c.Request().Body).Decode(&todo)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, response{Status: http.StatusInternalServerError, Res: map[string]interface{}{"error": "Internal server error"}})
	}

	creationError := todo.Validate()

	if creationError != nil {
		return c.JSON(http.StatusUnprocessableEntity, response{Status: http.StatusUnprocessableEntity, Res: map[string]interface{}{"error": creationError.Error()}})
	}

	_, insertError := s.db.Query(`INSERT INTO todos (title, completed) VALUES ($1, $2)`, todo.Title, todo.Completed)

	if insertError != nil {
		return c.JSON(http.StatusInternalServerError, response{Status: http.StatusInternalServerError, Res: map[string]interface{}{"error": insertError.Error()}})
	}

	resp := response{
		Status: http.StatusOK,
		Res:    map[string]interface{}{"data": todo},
	}

	return c.JSON(http.StatusCreated, resp)
}

func (s *Server) GetById(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))

	if !s.isExistTodo(id) {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Todo not found"})
	}

	row := s.db.QueryRow("SELECT title, completed, id, created_at, updated_at FROM todos WHERE id = $1", id)
	var todo todos.Todo
	err := row.Scan(&todo.Title, &todo.Completed, &todo.ID, &todo.CreatedAt, &todo.UpdatedAt)

	if err == sql.ErrNoRows {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Todo not found"})
	} else if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	resp := response{
		Status: http.StatusOK,
		Res:    map[string]interface{}{"data": todo},
	}

	return c.JSON(http.StatusCreated, resp)
}

func (s *Server) Update(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))

	if !s.isExistTodo(id) {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Todo not found"})
	}

	var todo todos.Todo
	_ = json.NewDecoder(c.Request().Body).Decode(&todo)

	_, err := s.db.Exec("UPDATE todos SET title = $1, completed = $2 WHERE id = $3", todo.Title, todo.Completed, id)

	if err != nil {
		return c.JSON(http.StatusBadRequest, response{Status: http.StatusBadRequest, Res: map[string]interface{}{"error": err.Error()}})
	}

	resp := response{
		Status: http.StatusOK,
		Res:    map[string]interface{}{"data": todo},
	}

	return c.JSON(http.StatusCreated, resp)
}

func (s *Server) Delete(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))

	if !s.isExistTodo(id) {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Todo not found"})
	}

	_, err := s.db.Exec("DELETE FROM todos WHERE id = $1", id)

	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	resp := response{
		Status: http.StatusOK,
	}

	return c.JSON(http.StatusOK, resp)
}

func (s *Server) isExistTodo(id int) bool {
	var tempId int
	err := s.db.QueryRow("SELECT id FROM todos WHERE id = $1", id).Scan(&tempId)

	return err == nil
}

// func (s *Server) HealthHandler(c echo.Context) error {
// 	return c.JSON(http.StatusOK, s.db.Health())
// }
