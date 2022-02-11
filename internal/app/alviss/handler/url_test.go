package handler_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/amirhosein/alviss/internal/app/alviss/handler"
	"github.com/amirhosein/alviss/internal/app/alviss/model"
	"github.com/amirhosein/alviss/internal/app/alviss/response"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/suite"
)

var (
	shortURL1 = "abcdef"
	shortURL2 = "abcdeg"

	expTime = time.Now()

	urls = map[string]model.URLMapping{
		shortURL1: {
			OriginalURL: "test.com/this/is/a/fake/long/url",
			Count:       20,
			ExpTime:     expTime,
		},
	}
)

type FakeURLRepo struct {
	model.URLRepo
	urls       map[string]model.URLMapping
	shouldFail bool
}

func (f FakeURLRepo) Get(shortURL string) (model.URLMapping, error) {
	if f.shouldFail {
		return model.URLMapping{}, errors.New("repo failure")
	}

	return f.urls[shortURL], nil
}

type URLHandlerSuite struct {
	suite.Suite
	fakeURLRepo *FakeURLRepo
	urlHandler  handler.URLHandler
	engine      *echo.Echo
}

func (suite *URLHandlerSuite) SetupSuite() {
	suite.fakeURLRepo = &FakeURLRepo{
		urls: urls,
	}

	suite.urlHandler = handler.URLHandler{
		Port:    "8080",
		URLRepo: suite.fakeURLRepo,
	}

	suite.engine = echo.New()
	suite.engine.GET("/", suite.urlHandler.Home)
	suite.engine.GET("/url/:shortURL", suite.urlHandler.HandleShortURLDetail)
}

func (suite *URLHandlerSuite) TestHome() {
	expected := response.Message{
		Message: "Welcome to Alviss! Your mythical URL shortener",
	}
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	c := suite.engine.NewContext(req, rec)

	suite.NoError(suite.urlHandler.Home(c))

	suite.Equal(http.StatusOK, rec.Code)
	if suite.Equal(http.StatusOK, rec.Code) {
		var resp response.Message
		suite.NoError(json.Unmarshal(rec.Body.Bytes(), &resp))

		suite.Equal(expected, resp)
	}
}

func (suite *URLHandlerSuite) TestHandleShortURLDetail() {
	cases := []struct {
		name     string
		shortURL string
		status   int
		repoFail bool
		expected response.Detail
	}{
		{
			name:     "successful",
			shortURL: shortURL1,
			status:   http.StatusOK,
			expected: response.Detail{
				OriginalURL: urls[shortURL1].OriginalURL,
				ShortURL:    "http://localhost:8080" + "/" + shortURL1,
				UsedCount:   urls[shortURL1].Count,
				ExpDate:     expTime.Format("2006-01-02 15:04:05"),
			},
		},
		{
			name:     "failure, url not found",
			shortURL: shortURL2,
			status:   http.StatusNotFound,
		},
		{
			name:     "failure, repo failure",
			shortURL: shortURL1,
			repoFail: true,
			status:   http.StatusInternalServerError,
		},
	}

	for i := range cases {
		tc := cases[i]

		suite.Run(tc.name, func() {
			suite.fakeURLRepo.shouldFail = tc.repoFail

			w := httptest.NewRecorder()
			req, err := http.NewRequest("GET", fmt.Sprintf("/url/%s", tc.shortURL), nil)

			suite.NoError(err)

			suite.engine.ServeHTTP(w, req)
			suite.Equal(tc.status, w.Code)

			if tc.status == http.StatusOK {
				var resp response.Detail
				suite.NoError(json.Unmarshal(w.Body.Bytes(), &resp))

				suite.Equal(tc.expected, resp)
			}
		})
	}
}

func TestURLHandlerSuite(t *testing.T) {
	suite.Run(t, &URLHandlerSuite{})
}
