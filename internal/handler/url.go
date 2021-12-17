package handler

import (
	"encoding/json"
	"log"
	"net/http"

	db "github.com/amirhosein/alviss/internal/DB"
	"github.com/amirhosein/alviss/internal/shortener"
	"github.com/amirhosein/alviss/internal/util"
	"github.com/labstack/echo/v4"
)

type UrlCreationRequest struct {
	LongUrl string `json:"long_url" binding:"required"`
	ExpDate string `json:"exp_date" binding:"required"`
	UserId  string `json:"user_id" binding:"required"`
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
		urlCreationRequest.LongUrl = json_map["long_url"].(string)
		urlCreationRequest.ExpDate = json_map["exp_date"].(string)
		urlCreationRequest.UserId = json_map["user_id"].(string)
	}
	log.Println("expDate: ", urlCreationRequest.ExpDate)
	log.Println("longUrl: ", util.GetExpireTime(urlCreationRequest.ExpDate))

	shortUrl := shortener.GenerateShortLink(urlCreationRequest.LongUrl, urlCreationRequest.UserId)
	error := db.SaveUrlMapping(shortUrl, urlCreationRequest.LongUrl, urlCreationRequest.UserId, util.GetExpireTime(urlCreationRequest.ExpDate))

	if error != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"message": error.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message":   "short url created successfully",
		"short_url": "http://localhost:" + port + "/" + shortUrl,
	})
}

func HandleShortUrlRedirect(c echo.Context) error {

	shortUrl := c.Param("shortUrl")
	initialUrl := db.RetrieveInitialUrl(shortUrl)
	return c.Redirect(302, initialUrl)

}
