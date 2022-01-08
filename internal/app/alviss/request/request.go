package request

import (
	"encoding/json"
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

func ShortURLCreationRequest(c echo.Context, SQLURLRepo model.SQLURLRepo, port string) error {
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
	error := SQLURLRepo.SaveURLMapping(shortURL, urlMapping, util.GetExpireTime(urlCreationRequest.ExpDate))

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
