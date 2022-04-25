package server

import (
	"fmt"
	"time"

	db "github.com/amirhosein/alviss/internal/app/alviss/DB"
	"github.com/amirhosein/alviss/internal/app/alviss/handler"
	"github.com/amirhosein/alviss/internal/app/alviss/model"
	"github.com/labstack/echo/v4"
)

var Port string

func Run(port string) {
	fmt.Println("Server is running on port " + port)
	time.Sleep(time.Second * 3)

	sqlURLRepo := model.SQLURLRepo{DB: db.InitDB()}
	cacheURLRepo := model.CacheURLRepo{URLDB: sqlURLRepo, RedisClient: db.InitRedis()}
	h := handler.URLHandler{
		Port:    port,
		URLRepo: cacheURLRepo,
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
