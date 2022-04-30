package handler_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/amirhosein/alviss/internal/app/alviss/handler"
	"github.com/amirhosein/alviss/internal/app/alviss/model"
	"github.com/amirhosein/alviss/internal/app/alviss/request"
	"github.com/amirhosein/alviss/internal/app/alviss/response"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/suite"
)

var (
	shortURL1 = "abcdef"
	shortURL2 = "abcdeg"
	shortURL3 = "abcdeh"

	expTime = time.Now()

	urls = map[string]model.URLMapping{
		shortURL1: {
			OriginalURL: "test.com/this/is/a/fake/long/url",
			Count:       20,
			ExpTime:     expTime,
		},
		shortURL3: {
			OriginalURL: "test.com/this/is/a/fake/long/url",
			Count:       20,
			ExpTime:     expTime.AddDate(0, 0, -1),
		},
	}
)

type InvalidRequest struct {
	WrongURL string `json:"WrongURL" binding:"required"`
}

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

func (f FakeURLRepo) Save(shortURL string, urlMapping model.URLMapping, expTime time.Duration) error {
	if f.shouldFail {
		return errors.New("repo failure")
	}

	urls[shortURL] = urlMapping
	shortURL1 = shortURL

	return nil
}

func (f FakeURLRepo) Update(shortURL string, urlMapping model.URLMapping) error {
	if f.shouldFail {
		return errors.New("repo failure")
	}

	return nil
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
	suite.engine.POST("/shorten", suite.urlHandler.CreateShortURL)
	suite.engine.GET("/url/:shortURL", suite.urlHandler.HandleShortURLDetail)
	suite.engine.GET("/:shortURL", suite.urlHandler.HandleShortURLRedirect)
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

func (suite *URLHandlerSuite) TestCreateShortURL() {
	cases := []struct {
		name     string
		status   int
		repoFail bool
		request  request.URLCreationRequest
		expected response.SuccessfullyCreated
	}{
		{
			name:   "successful",
			status: http.StatusOK,
			request: request.URLCreationRequest{
				LongURL: urls[shortURL1].OriginalURL,
				ExpDate: "2d",
			},
			expected: response.SuccessfullyCreated{
				Message:  "Short url created successfully",
				ShortURL: fmt.Sprintf("http://localhost:8080/%s", shortURL1),
			},
		},
		{
			name:   "invalid request",
			status: http.StatusNotAcceptable,
			request: request.URLCreationRequest{
				LongURL: "",
				ExpDate: "2d",
			},
			expected: response.SuccessfullyCreated{
				Message: "ExpTime: cannot be blank; LongURL: cannot be blank.",
			},
		},
		{
			name:     "failure, repo failure",
			repoFail: true,
			status:   http.StatusInternalServerError,
			request: request.URLCreationRequest{
				LongURL: urls[shortURL1].OriginalURL,
				ExpDate: expTime.Format("2006-01-02 15:04:05"),
			},
			expected: response.SuccessfullyCreated{
				Message: "repo failure",
			},
		},
	}

	for i := range cases {
		tc := cases[i]
		suite.Run(tc.name, func() {
			suite.fakeURLRepo.shouldFail = tc.repoFail
			w := httptest.NewRecorder()

			reqBody, err := json.Marshal(tc.request)
			if tc.name == "invalid request" {
				reqBody, err = json.Marshal(InvalidRequest{
					WrongURL: "test.com/this/is/a/fake/long/url",
				})
			}
			suite.NoError(err)

			req, err := http.NewRequest("POST", "/shorten", bytes.NewBuffer(reqBody))
			req.Header.Set("Content-Type", "application/json; charset=UTF-8")
			suite.NoError(err)

			suite.engine.ServeHTTP(w, req)
			suite.Equal(tc.status, w.Code)

			var resp response.SuccessfullyCreated
			suite.NoError(json.Unmarshal(w.Body.Bytes(), &resp))
			suite.Equal(tc.expected.Message, resp.Message)

			if tc.status == http.StatusOK {
				suite.Equal(tc.expected.Message, resp.Message)
				_, err = url.ParseRequestURI(resp.ShortURL)
				suite.NoError(err)
			}
		})
	}
}

func (suite *URLHandlerSuite) TestHandleShortURLRedirect() {
	cases := []struct {
		name     string
		shortURL string
		status   int
		repoFail bool
		expected string
	}{
		{
			name:     "successful",
			shortURL: shortURL1,
			status:   http.StatusMovedPermanently,
			expected: urls[shortURL1].OriginalURL,
		},
		{
			name:     "failure, url expired",
			shortURL: shortURL3,
			status:   http.StatusGone,
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
			req, err := http.NewRequest("GET", fmt.Sprintf("/%s", tc.shortURL), nil)

			suite.NoError(err)
			suite.engine.ServeHTTP(w, req)
			suite.Equal(tc.status, w.Code)

			if tc.status == http.StatusOK {
				suite.Equal(tc.expected, w.Header().Get("Location"))
			}
		})
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
				ExpDate:     urls[shortURL1].ExpTime.Format("2006-01-02 15:04:05"),
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
