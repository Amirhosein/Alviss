package handler

import (
	"log"
	"net/http"
	"time"

	"github.com/amirhosein/alviss/internal/app/alviss/model"
	"github.com/amirhosein/alviss/internal/app/alviss/request"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	Port       string
	SQLURLRepo model.SQLURLRepo
}

type URLCreationRequest struct {
	LongURL string `json:"LongURL" binding:"required"`
	ExpDate string `json:"ExpTime" binding:"required"`
}

func (h Handler) Home(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Welcome to Alviss! Your mythical URL shortener",
	})
}

func (h Handler) CreateShortURL(c echo.Context) error {
	return request.ShortURLCreationRequest(c, h.SQLURLRepo, h.Port)
}

func (h Handler) HandleShortURLRedirect(c echo.Context) error {
	shortURL := c.Param("shortURL")
	result, err := h.SQLURLRepo.GetURLMapping(shortURL)
	if err != nil {
		log.Println(err)
	}
	if (model.URLMapping{}) == result {
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"message": "Short url not found",
		})
	}
	if !result.ExpTime.IsZero() && result.ExpTime.Before(time.Now()) {
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"message": "Short url expired",
		})
	}
	result.Count++
	err = h.SQLURLRepo.Update(shortURL, result)
	if err != nil {
		log.Println(err)
	}
	return c.Redirect(http.StatusMovedPermanently, result.Original_url)
}

func (h Handler) HandleShortURLDetail(c echo.Context) error {
	shortURL := c.Param("shortURL")
	result, err := h.SQLURLRepo.GetURLMapping(shortURL)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"message": err.Error(),
		})
	}
	if (model.URLMapping{}) == result {
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"message": "Short url not found",
		})
	} else {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"OriginalURL": result.Original_url,
			"ShortURL":    "http://localhost:" + h.Port + "/" + shortURL,
			"UsedCount":   result.Count,
			"ExpDate":     result.ExpTime,
		})
	}
}
