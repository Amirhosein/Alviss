package handler

import (
	"encoding/json"
	"net/http"
	"time"

	db "github.com/amirhosein/alviss/internal/app/alviss/DB"
	"github.com/amirhosein/alviss/internal/app/alviss/util"
	"github.com/labstack/echo/v4"
)

type UrlCreationRequest struct {
	LongUrl string `json:"LongUrl" binding:"required"`
	ExpDate string `json:"ExpTime" binding:"required"`
}

func Home(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Welcome to Alviss! Your mythical URL shortener",
	})
}

func CreateShortUrl(c echo.Context, port string) error {
	urlCreationRequest := new(UrlCreationRequest)
	json_map := make(map[string]interface{})
	err := json.NewDecoder(c.Request().Body).Decode(&json_map)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "Invalid request body",
		})
	} else {
		urlCreationRequest.LongUrl = json_map["LongUrl"].(string)
		urlCreationRequest.ExpDate = json_map["ExpTime"].(string)
	}
	urlMapping := db.UrlMapping{
		Original_url: urlCreationRequest.LongUrl,
		Count:        0,
		ExpTime:      time.Now().Add(util.GetExpireTime(urlCreationRequest.ExpDate)),
	}

	shortUrl := util.GenerateShortLink(urlCreationRequest.LongUrl)
	error := db.SaveUrlMapping(shortUrl, urlMapping, util.GetExpireTime(urlCreationRequest.ExpDate))

	if error != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"message": error.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message":  "short url created successfully",
		"ShortUrl": "http://localhost:" + port + "/" + shortUrl,
	})
}

func HandleShortUrlRedirect(c echo.Context) error {

	shortUrl := c.Param("shortUrl")
	initialUrl := db.RetrieveInitialUrl(shortUrl)
	return c.Redirect(http.StatusFound, initialUrl)

}

func HandleShortUrlDetail(c echo.Context, port string) error {

	shortUrl := c.Param("shortUrl")
	result := db.RetrieveUrlMapping(shortUrl)
	if result == nil {
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"message": "short url not found",
		})
	} else {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"OriginalUrl": result.Original_url,
			"ShortUrl":    "http://localhost:" + port + "/" + shortUrl,
			"UsedCount":   result.Count,
			"ExpDate":     result.ExpTime,
		})
	}

}