package server

import (
	"fmt"

	"github.com/labstack/echo/v4"
)

func RunServer(port string) {
	fmt.Println("Server is running on port " + port)

	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(200, "Hello, World!")
	})
	e.Logger.Fatal(e.Start(":" + port))
}
