package server

import (
	"fmt"
	"time"

	db "github.com/amirhosein/alviss/internal/app/alviss/DB"
	"github.com/amirhosein/alviss/internal/app/alviss/handler"
	"github.com/amirhosein/alviss/internal/app/alviss/model"
	"github.com/labstack/echo/v4"
)

var (
	ServerPort string
)

func RunServer(port string) {
	ServerPort = port
	fmt.Println("Server is running on port " + port)
	// sleep for a while to wait for the database to be ready
	time.Sleep(time.Second * 3)

	sqlUrlRepo := model.SQLUrlRepo{DB: db.InitDB()}

	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return handler.Home(c)
	})

	e.POST("/shorten", func(c echo.Context) error {
		return handler.CreateShortUrl(c, ServerPort, sqlUrlRepo)
	})

	e.GET("/:shortUrl", func(c echo.Context) error {
		return handler.HandleShortUrlRedirect(c, sqlUrlRepo)
	})

	e.GET("/url/:shortUrl", func(c echo.Context) error {
		return handler.HandleShortUrlDetail(c, ServerPort, sqlUrlRepo)
	})

	e.Logger.Fatal(e.Start(":" + port))
}
