package server

import (
	"fmt"
	"net/http"

	db "github.com/amirhosein/alviss/internal/DB"
	"github.com/amirhosein/alviss/internal/handler"
	"github.com/labstack/echo/v4"
)

var (
	ServerPort string
)

func RunServer(port string) {
	ServerPort = port
	fmt.Println("Server is running on port " + port)
	db.InitializeStore()

	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"message": "Welcome to Alviss! Your mythical URL shortener",
		})
	})

	e.POST("/shorten", func(c echo.Context) error {
		return handler.CreateShortUrl(c, ServerPort)
	})

	e.GET("/:shortUrl", func(c echo.Context) error {
		return handler.HandleShortUrlRedirect(c)
	})

	e.GET("/url/:shortUrl", func(c echo.Context) error {
		return handler.HandleShortUrlDetail(c, ServerPort)
	})

	e.Logger.Fatal(e.Start(":" + port))
}
