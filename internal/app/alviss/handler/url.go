package handler

import (
	"log"
	"net/http"
	"time"

	"github.com/amirhosein/alviss/internal/app/alviss/model"
	"github.com/amirhosein/alviss/internal/app/alviss/request"
	"github.com/amirhosein/alviss/internal/app/alviss/response"
	"github.com/amirhosein/alviss/internal/app/alviss/util"

	"github.com/labstack/echo/v4"
)

type URLHandler struct {
	Port    string
	URLRepo model.URLRepo
}

func (h URLHandler) Home(c echo.Context) error {
	message := response.Message{
		Message: "Welcome to Alviss! Your mythical URL shortener",
	}

	return c.JSON(http.StatusOK, message)
}

func (h URLHandler) CreateShortURL(c echo.Context) error {
	urlCreationRequest := new(request.URLCreationRequest)

	err := c.Bind(urlCreationRequest)
	if err != nil {
		log.Print(err)
	}

	if err := urlCreationRequest.Validate(); err != nil {
		message := response.Message{
			Message: err.Error(),
		}

		return c.JSON(http.StatusNotAcceptable, message)
	}

	urlMapping := model.URLMapping{
		OriginalURL: urlCreationRequest.LongURL,
		Count:       0,
		ExpTime:     time.Now().Add(util.GetExpireTime(urlCreationRequest.ExpDate)),
	}

	shortURL := util.GenerateShortLink(urlCreationRequest.LongURL)

	err = h.URLRepo.Save(shortURL, urlMapping, util.GetExpireTime(urlCreationRequest.ExpDate))
	if err != nil {
		message := response.Message{
			Message: err.Error(),
		}

		return c.JSON(http.StatusInternalServerError, message)
	}

	successfullyCreated := response.SuccessfullyCreated{
		Message:  "Short url created successfully",
		ShortURL: "http://localhost:" + h.Port + "/" + shortURL,
	}

	return c.JSON(http.StatusOK, successfullyCreated)
}

func (h URLHandler) HandleShortURLRedirect(c echo.Context) error {
	shortURL := c.Param("shortURL")

	result, err := h.URLRepo.Get(shortURL)
	if err != nil {
		log.Println(err)

		message := response.Message{
			Message: "Short url not found",
		}

		return c.JSON(http.StatusInternalServerError, message)
	}

	if (model.URLMapping{}) == result {
		message := response.Message{
			Message: "Short url not found",
		}

		return c.JSON(http.StatusNotFound, message)
	}

	if !result.ExpTime.IsZero() && result.ExpTime.Before(time.Now()) {
		message := response.Message{
			Message: "Short url expired",
		}

		return c.JSON(http.StatusGone, message)
	}
	result.Count++

	err = h.URLRepo.Update(shortURL, result)
	if err != nil {
		log.Println(err)
	}

	return c.Redirect(http.StatusMovedPermanently, result.OriginalURL)
}

func (h URLHandler) HandleShortURLDetail(c echo.Context) error {
	shortURL := c.Param("shortURL")

	result, err := h.URLRepo.Get(shortURL)
	if err != nil {
		message := response.Message{
			Message: err.Error(),
		}

		return c.JSON(http.StatusInternalServerError, message)
	}

	if (model.URLMapping{}) == result {
		message := response.Message{
			Message: "Short url not found",
		}

		return c.JSON(http.StatusNotFound, message)
	} else {
		detail := response.Detail{
			OriginalURL: result.OriginalURL,
			ShortURL:    "http://localhost:" + h.Port + "/" + shortURL,
			UsedCount:   result.Count,
			ExpDate:     result.ExpTime.Format("2006-01-02 15:04:05"),
		}
		return c.JSON(http.StatusOK, detail)
	}
}
