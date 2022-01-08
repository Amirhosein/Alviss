package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/amirhosein/alviss/internal/app/alviss/model"
	"github.com/amirhosein/alviss/internal/app/alviss/util"
	"github.com/labstack/echo/v4"
)

type URLCreationRequest struct {
	LongURL string `json:"LongURL" binding:"required"`
	ExpDate string `json:"ExpTime" binding:"required"`
}

func Home(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Welcome to Alviss! Your mythical URL shortener",
	})
}

func CreateShortURL(c echo.Context, port string, sqlURLRepo model.SQLURLRepo) error {
	urlCreationRequest := new(URLCreationRequest)
	json_map := make(map[string]interface{})
	err := json.NewDecoder(c.Request().Body).Decode(&json_map)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "Invalid request body",
		})
	} else {
		urlCreationRequest.LongURL = json_map["LongURL"].(string)
		urlCreationRequest.ExpDate = json_map["ExpTime"].(string)
	}
	urlMapping := model.URLMapping{
		Original_url: urlCreationRequest.LongURL,
		Count:        0,
		ExpTime:      time.Now().Add(util.GetExpireTime(urlCreationRequest.ExpDate)),
	}

	shortURL := util.GenerateShortLink(urlCreationRequest.LongURL)
	error := sqlURLRepo.SaveURLMapping(shortURL, urlMapping, util.GetExpireTime(urlCreationRequest.ExpDate))

	if error != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"message": error.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message":  "short url created successfully",
		"ShortURL": "http://localhost:" + port + "/" + shortURL,
	})
}

func HandleShortURLRedirect(c echo.Context, sqlURLRepo model.SQLURLRepo) error {
	shortURL := c.Param("shortURL")
	result, err := sqlURLRepo.GetURLMapping(shortURL)
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
	sqlURLRepo.UpdateURLMapping(shortURL, result)
	return c.Redirect(http.StatusMovedPermanently, result.Original_url)
}

func HandleShortURLDetail(c echo.Context, port string, sqlURLRepo model.SQLURLRepo) error {
	shortURL := c.Param("shortURL")
	result, err := sqlURLRepo.GetURLMapping(shortURL)
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
			"ShortURL":    "http://localhost:" + port + "/" + shortURL,
			"UsedCount":   result.Count,
			"ExpDate":     result.ExpTime,
		})
	}
}
