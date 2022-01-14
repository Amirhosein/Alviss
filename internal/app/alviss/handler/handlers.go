package handler

import (
	"log"
	"net/http"
	"time"

	"github.com/amirhosein/alviss/internal/app/alviss/model"
	"github.com/amirhosein/alviss/internal/app/alviss/request"
	"github.com/amirhosein/alviss/internal/app/alviss/util"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	Port    string
	URLRepo model.URLRepo
}

func (h Handler) Home(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Welcome to Alviss! Your mythical URL shortener",
	})
}

func (h Handler) CreateShortURL(c echo.Context) error {
	urlCreationRequest := new(request.URLCreationRequest)

	if err := c.Bind(urlCreationRequest); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "Invalid request body",
		})
	}

	if err := urlCreationRequest.Validate(); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": err.Error(),
		})
	}

	urlMapping := model.URLMapping{
		OriginalURL: urlCreationRequest.LongURL,
		Count:       0,
		ExpTime:     time.Now().Add(util.GetExpireTime(urlCreationRequest.ExpDate)),
	}

	shortURL := util.GenerateShortLink(urlCreationRequest.LongURL)

	err := h.URLRepo.Save(shortURL, urlMapping, util.GetExpireTime(urlCreationRequest.ExpDate))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message":  "short url created successfully",
		"ShortURL": "http://localhost:" + h.Port + "/" + shortURL,
	})
}

func (h Handler) HandleShortURLRedirect(c echo.Context) error {
	shortURL := c.Param("shortURL")
	result, err := h.URLRepo.Get(shortURL)
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
	err = h.URLRepo.Update(shortURL, result)
	if err != nil {
		log.Println(err)
	}
	return c.Redirect(http.StatusMovedPermanently, result.OriginalURL)
}

func (h Handler) HandleShortURLDetail(c echo.Context) error {
	shortURL := c.Param("shortURL")
	result, err := h.URLRepo.Get(shortURL)
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
			"OriginalURL": result.OriginalURL,
			"ShortURL":    "http://localhost:" + h.Port + "/" + shortURL,
			"UsedCount":   result.Count,
			"ExpDate":     result.ExpTime,
		})
	}
}
