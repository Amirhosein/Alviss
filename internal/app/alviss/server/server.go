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
	fmt.Println("Server is running on port " + port)
	// sleep for a while to wait for the database to be ready
	time.Sleep(time.Second * 3)

	sqlURLRepo := model.SQLURLRepo{DB: db.InitDB()}
	h := handler.Handler{
		Port:       port,
		SQLURLRepo: sqlURLRepo,
	}

	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return h.Home(c)
	})

	e.POST("/shorten", func(c echo.Context) error {
		return h.CreateShortURL(c)
	})

	e.GET("/:shortURL", func(c echo.Context) error {
		return h.HandleShortURLRedirect(c)
	})

	e.GET("/url/:shortURL", func(c echo.Context) error {
		return h.HandleShortURLDetail(c)
	})

	e.Logger.Fatal(e.Start(":" + port))
}
