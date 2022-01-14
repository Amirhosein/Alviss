package response

import (
	"net/http"

	"github.com/amirhosein/alviss/internal/app/alviss/model"
	"github.com/labstack/echo/v4"
)

func Welcome(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Welcome to Alviss! Your mythical URL shortener",
	})
}

func InvalidRequest(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Invalid request body",
	})
}

func NotFound(c echo.Context) error {
	return c.JSON(http.StatusNotFound, map[string]interface{}{
		"message": "Short url not found",
	})
}

func Expired(c echo.Context) error {
	return c.JSON(http.StatusNotFound, map[string]interface{}{
		"message": "Short url expired",
	})
}

func Error(c echo.Context, err error) error {
	return c.JSON(http.StatusBadRequest, map[string]interface{}{
		"message": err.Error(),
	})
}

func SuccessfullyCreated(c echo.Context, port string, shortURL string) error {
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message":  "short url created successfully",
		"ShortURL": "http://localhost:" + port + "/" + shortURL,
	})
}

func Detail(c echo.Context, result model.URLMapping, shortURL string, port string) error {
	return c.JSON(http.StatusOK, map[string]interface{}{
		"OriginalURL": result.OriginalURL,
		"ShortURL":    "http://localhost:" + port + "/" + shortURL,
		"UsedCount":   result.Count,
		"ExpDate":     result.ExpTime,
	})
}
