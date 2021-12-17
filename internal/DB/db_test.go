package db

import (
	"testing"

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
	initialLink := "https://www.guru3d.com/news-story/spotted-ryzen-threadripper-pro-3995wx-processor-with-8-channel-ddr4,2.html"
	userUUId := "e0dba740-fc4b-4977-872c-d360239e6b1a"
	shortURL := "Jsz4k57oAX"
	expireTime := util.GetExpireTime("2d")

	// Persist data mapping
	SaveUrlMapping(shortURL, initialLink, userUUId, expireTime)

	// Retrieve initial URL
	retrievedUrl := RetrieveInitialUrl(shortURL)

	assert.Equal(t, initialLink, retrievedUrl)
}
