package db

import (
	"testing"
	"time"

	"github.com/amirhosein/alviss/internal/shortener"
	"github.com/amirhosein/alviss/internal/util"
	"github.com/stretchr/testify/assert"
)

var testStoreService = &DBService{}

func init() {
	testStoreService = InitializeStore()
}

func TestStoreInit(t *testing.T) {
	assert.True(t, testStoreService.redisClient != nil)
}

func TestInsertionAndRetrieval(t *testing.T) {
	urlMapping := new(UrlMapping)

	urlMapping.OriginalUrl = "https://www.guru3d.com/news-story/spotted-ryzen-threadripper-pro-3995wx-processor-with-8-channel-ddr4,2.html"
	expireTime := util.GetExpireTime("2d")
	urlMapping.ExpTime = time.Now().Add(expireTime)
	shortURL := shortener.GenerateShortLink(urlMapping.OriginalUrl)

	// Persist data mapping
	err := SaveUrlMapping(shortURL, *urlMapping, expireTime)
	assert.Nil(t, err)

	// Retrieve initial URL
	retrievedUrl := RetrieveInitialUrl(shortURL)

	assert.Equal(t, urlMapping.OriginalUrl, retrievedUrl)
}
