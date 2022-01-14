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

type Handler struct {
	Port    string
	URLRepo model.URLRepo
}

func (h Handler) Home(c echo.Context) error {
	return response.Welcome(c)
}

func (h Handler) CreateShortURL(c echo.Context) error {
	urlCreationRequest := new(request.URLCreationRequest)

	if err := c.Bind(urlCreationRequest); err != nil {
		return response.InvalidRequest(c)
	}

	if err := urlCreationRequest.Validate(); err != nil {
		return response.Error(c, err)
	}

	urlMapping := model.URLMapping{
		OriginalURL: urlCreationRequest.LongURL,
		Count:       0,
		ExpTime:     time.Now().Add(util.GetExpireTime(urlCreationRequest.ExpDate)),
	}

	shortURL := util.GenerateShortLink(urlCreationRequest.LongURL)

	err := h.URLRepo.Save(shortURL, urlMapping, util.GetExpireTime(urlCreationRequest.ExpDate))
	if err != nil {
		return response.Error(c, err)
	}

	return response.SuccessfullyCreated(c, h.Port, shortURL)
}

func (h Handler) HandleShortURLRedirect(c echo.Context) error {
	shortURL := c.Param("shortURL")
	result, err := h.URLRepo.Get(shortURL)
	if err != nil {
		log.Println(err)
	}
	if (model.URLMapping{}) == result {
		return response.NotFound(c)
	}
	if !result.ExpTime.IsZero() && result.ExpTime.Before(time.Now()) {
		return response.Expired(c)
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
		return response.Error(c, err)
	}
	if (model.URLMapping{}) == result {
		return response.NotFound(c)
	} else {
		return response.Detail(c, result, shortURL, h.Port)
	}
}
