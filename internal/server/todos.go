package server

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (s *Server) HelloWorldHandler(c echo.Context) error {
	resp := map[string]string{
		"message": "Hello World",
	}

	return c.JSON(http.StatusOK, resp)
}

func (s *Server) HealthHandler(c echo.Context) error {
	return c.JSON(http.StatusOK, s.db.Health())
}

type Data struct {
	Value int
}

func (d Data) Hello() {
	fmt.Println(d.Value)
}
