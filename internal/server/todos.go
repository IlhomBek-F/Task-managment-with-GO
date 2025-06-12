package server

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

type Message struct {
	Name      string `json:"name"`
	Completed bool   `json:"completed"`
	ID        int    `json:"id"`
}

func (s *Server) HelloWorldHandler(c echo.Context) error {

	rows, err := s.db.Query("SELECT * FROM  todos")

	var (
		name      string
		completed bool
		idTodo    int
	)

	todos := []Message{}

	if err != nil {
		fmt.Println("Error database ")
	} else {
		fmt.Println(rows)
	}

	for rows.Next() {
		err := rows.Scan(&name, &completed, &idTodo)

		fmt.Println(name, completed, idTodo)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "row scan failed"})
		}

		todo := Message{Name: name, Completed: completed, ID: idTodo}

		fmt.Println(todo)
		todos = append(todos, todo)
	}

	resp := map[string][]Message{
		"message": todos,
	}
	defer rows.Close()

	return c.JSON(http.StatusOK, resp)
}

// func (s *Server) HealthHandler(c echo.Context) error {
// 	return c.JSON(http.StatusOK, s.db.Health())
// }
