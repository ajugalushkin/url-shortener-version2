package handler

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"

	"github.com/ajugalushkin/url-shortener-version2/config"
	"github.com/ajugalushkin/url-shortener-version2/internal/cookies"
	"github.com/ajugalushkin/url-shortener-version2/internal/dto"
	"github.com/ajugalushkin/url-shortener-version2/internal/service"
	"github.com/ajugalushkin/url-shortener-version2/internal/storage/inmemory"
)

//var newConfig = config.AppConfig{
//	ServerAddress: "localhost:8080",
//	BaseURL:       "http://localhost:8080",
//}

// var ctx = config.ContextWithFlags(context.Background(), &newConfig)
var ctx = context.Background()

func TestHandler_HandleSave(t *testing.T) {
	type request struct {
		method      string
		body        string
		contentType string
	}
	type repository []dto.Shortening
	type want struct {
		code        int
		contentType string
	}
	tests := []struct {
		name       string
		request    request
		repository repository
		want       want
	}{
		{
			name: "Test StatusCreated",
			request: request{
				method:      http.MethodPost,
				body:        "https://practicum.yandex.ru/",
				contentType: echo.MIMETextPlain,
			},
			want: want{
				code:        http.StatusCreated,
				contentType: echo.MIMETextPlainCharsetUTF8,
			},
		},
		{
			name: "Test Empty URL",
			request: request{
				method:      http.MethodPost,
				body:        "",
				contentType: echo.MIMETextPlain,
			},
			want: want{
				code:        http.StatusBadRequest,
				contentType: echo.MIMETextPlainCharsetUTF8,
			},
		},
		{
			name: "Test Duplicate URL",
			request: request{
				method:      http.MethodPost,
				body:        "https://practicum.yandex.ru/",
				contentType: echo.MIMETextPlain,
			},
			repository: repository{{OriginalURL: "https://practicum.yandex.ru/"}},
			want: want{
				code:        http.StatusConflict,
				contentType: echo.MIMETextPlainCharsetUTF8,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Setup
			server := echo.New()
			req := httptest.NewRequest(test.request.method, "/", strings.NewReader(test.request.body))
			req.Header.Set(echo.HeaderContentType, test.request.contentType)
			rec := httptest.NewRecorder()
			echoCtx := server.NewContext(req, rec)

			repoMemory := inmemory.NewInMemory()
			for _, repositoryItem := range test.repository {
				_, err := repoMemory.Put(ctx, repositoryItem)
				assert.Nil(t, err)
			}
			handler := NewHandler(ctx, service.NewService(repoMemory))

			// Assertions
			if assert.NoError(t, handler.HandleSave(echoCtx)) {
				assert.Equal(t, test.want.code, rec.Code)
				assert.Equal(t, test.want.contentType, rec.Header().Get(echo.HeaderContentType))
			}
		})
	}
}

func TestHandler_HandleRedirect(t *testing.T) {
	type request struct {
		method      string
		key         string
		URL         string
		contentType string
	}
	type want struct {
		code     int
		response string
	}
	tests := []struct {
		name    string
		request request
		want    want
	}{
		{
			name: "Test OK",
			request: request{
				method:      http.MethodGet,
				key:         "rIHY5pi",
				URL:         "http://localhost:8080/rIHY5pi",
				contentType: echo.MIMETextPlain,
			},
			want: want{
				code:     http.StatusTemporaryRedirect,
				response: "https://practicum.yandex.ru/",
			},
		},
		{
			name: "Test Bad Request 1",
			request: request{
				method:      http.MethodGet,
				URL:         "http://localhost:8080/rIHY5pi",
				contentType: echo.MIMETextPlain,
			},
			want: want{
				code: http.StatusBadRequest,
			},
		},
		{
			name: "Test Bad Request 2",
			request: request{
				URL:    "http://localhost:8080/rIHY5pi",
				method: http.MethodPost,
			},
			want: want{
				code: http.StatusBadRequest,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Setup
			server := echo.New()
			req := httptest.NewRequest(test.request.method, test.request.URL, nil)
			req.Header.Set(echo.HeaderContentType, test.request.contentType)
			rec := httptest.NewRecorder()
			echoCtx := server.NewContext(req, rec)

			storageAPI := inmemory.NewInMemory()
			var err error
			if test.request.key != "" && test.want.response != "" {
				_, err = storageAPI.Put(ctx, dto.Shortening{
					ShortURL:    test.request.key,
					OriginalURL: test.want.response,
				})
			}
			if assert.NoError(t, err) {
				handler := NewHandler(ctx, service.NewService(storageAPI))

				// Assertions
				if assert.NoError(t, handler.HandleRedirect(echoCtx)) {
					assert.Equal(t, test.want.code, rec.Code)
					assert.Equal(t, test.want.response, rec.Header().Get(echo.HeaderLocation))
				}
			}
		})
	}
}

// Redirects to the original URL if it exists and is not deleted
func TestHandleRedirect_RedirectsToOriginalURL(t *testing.T) {
	// Arrange
	ctx := context.Background()
	serviceAPI := service.NewService(inmemory.NewInMemory())

	handler := &Handler{
		ctx:     ctx,
		servAPI: serviceAPI,
	}

	server := echo.New()

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.URL.Path = "/test"
	rec := httptest.NewRecorder()
	echoCtx := server.NewContext(req, rec)

	shortening := &dto.Shortening{
		ShortURL:    "test",
		OriginalURL: "http://example.com",
		IsDeleted:   false,
	}

	_, err := serviceAPI.Shorten(ctx, *shortening)
	if err != nil {
		assert.NoError(t, err)
	}

	_, err = serviceAPI.Redirect(ctx, "test")
	if err != nil {
		assert.NoError(t, err)
	}

	// Act
	err = handler.HandleRedirect(echoCtx)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, http.StatusTemporaryRedirect, echoCtx.Response().Status)
	assert.Equal(t, "http://example.com", echoCtx.Response().Header().Get("Location"))
}

func TestHandler_HandleShorten(t *testing.T) {
	type request struct {
		method      string
		body        string
		contentType string
	}
	type want struct {
		code        int
		contentType string
	}
	tests := []struct {
		name    string
		request request
		want    want
	}{
		{
			name: "Test StatusCreated",
			request: request{
				method:      http.MethodPost,
				body:        "{\n  \"url\": \"https://practicum.yandex.ru\"\n}",
				contentType: echo.MIMEApplicationJSON,
			},
			want: want{
				code:        http.StatusCreated,
				contentType: echo.MIMEApplicationJSON,
			},
		},
		{
			name: "Test Empty URL",
			request: request{
				method:      http.MethodPost,
				body:        "{\n  \"url\": \"\"\n}",
				contentType: echo.MIMEApplicationJSON,
			},
			want: want{
				code:        http.StatusBadRequest,
				contentType: echo.MIMETextPlainCharsetUTF8,
			},
		},
		{
			name: "Test Empty Body",
			request: request{
				method:      http.MethodPost,
				body:        "",
				contentType: echo.MIMEApplicationJSON,
			},
			want: want{
				code:        http.StatusBadRequest,
				contentType: echo.MIMETextPlainCharsetUTF8,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Setup
			server := echo.New()

			req := httptest.NewRequest(test.request.method, "/api/shorten", strings.NewReader(test.request.body))
			req.Header.Set(echo.HeaderContentType, test.request.contentType)

			rec := httptest.NewRecorder()
			newContext := server.NewContext(req, rec)

			handler := NewHandler(ctx, service.NewService(inmemory.NewInMemory()))

			// Assertions
			if assert.NoError(t, handler.HandleShorten(newContext)) {
				assert.Equal(t, test.want.code, rec.Code)
				assert.Equal(t, test.want.contentType, rec.Header().Get(echo.HeaderContentType))
			}
		})
	}
}

func TestHandler_HandleShortenBatch(t *testing.T) {
	type fields struct {
		ctx     context.Context
		servAPI *service.Service
	}
	tests := []struct {
		name             string
		fields           fields
		inputContentType string
		inputMethod      string
		inputBody        string
		expectedHeader   string
		expectedCode     int
	}{
		{
			name: "Test ОК",
			fields: fields{
				ctx:     ctx,
				servAPI: service.NewService(inmemory.NewInMemory())},
			inputContentType: echo.MIMEApplicationJSON,
			inputMethod:      http.MethodPost,
			inputBody:        "[\n    {\n        \"correlation_id\": \"1\",\n        \"original_url\": \"https://vk.com/ajugalushkin\"\n    }\n]",
			expectedHeader:   echo.MIMEApplicationJSON,
			expectedCode:     http.StatusCreated,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Setup
			req := httptest.NewRequest(
				test.inputMethod,
				"/api/shorten/batch",
				strings.NewReader(test.inputBody),
			)
			req.Header.Set(echo.HeaderContentType, test.inputContentType)
			rec := httptest.NewRecorder()

			echoCtx := echo.New().NewContext(req, rec)

			handler := Handler{ctx: test.fields.ctx, servAPI: test.fields.servAPI}

			// Assertions
			if assert.NoError(t, handler.HandleShortenBatch(echoCtx)) {
				assert.Equal(t, test.expectedHeader, rec.Header().Get(echo.HeaderContentType))
				assert.Equal(t, test.expectedCode, rec.Code)
			}
		})
	}
}

func dummyHandler(c echo.Context) error {
	return c.JSON(http.StatusOK, dto.UserURLList{})
}

func TestHandler_Authorized(t *testing.T) {
	t.Run("Test Authorized", func(t *testing.T) {
		// Setup
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/api/user/urls", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		cookie := cookies.CreateCookie(ctx, cookieName)
		req.AddCookie(cookie)

		URLSInMem := inmemory.NewInMemory()
		_, err := URLSInMem.Put(ctx, dto.Shortening{
			CorrelationID: "1",
			ShortURL:      "34ewfd",
			OriginalURL:   "http://test.com",
			UserID:        strconv.Itoa(cookies.GetUser(ctx, cookie.Value).ID)})
		if err != nil {
			return
		}

		h := Handler{
			ctx:     ctx,
			cache:   map[string]*dto.User{cookie.Value: cookies.GetUser(ctx, cookie.Value)},
			servAPI: service.NewService(URLSInMem),
		}
		c := e.NewContext(req, rec)

		handler := h.Authorized(dummyHandler)

		// Assertions
		if assert.NoError(t, handler(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
		}
	})

	t.Run("Test Authorized Empty Cookie", func(t *testing.T) {
		// Setup
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/api/user/urls", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		cookie := cookies.CreateCookie(ctx, cookieName)

		URLSInMem := inmemory.NewInMemory()
		_, err := URLSInMem.Put(ctx, dto.Shortening{
			CorrelationID: "1",
			ShortURL:      "34ewfd",
			OriginalURL:   "http://test.com",
			UserID:        strconv.Itoa(cookies.GetUser(ctx, cookie.Value).ID)})
		if err != nil {
			return
		}

		h := Handler{
			ctx:     ctx,
			cache:   make(map[string]*dto.User),
			servAPI: service.NewService(URLSInMem),
		}

		c := e.NewContext(req, rec)

		handler := h.Authorized(dummyHandler)

		// Assertions
		if assert.NoError(t, handler(c)) {
			assert.Equal(t, http.StatusUnauthorized, rec.Code)
		}
	})

	//t.Run("Test Authorized Wrong Cookie", func(t *testing.T) {
	//	// Setup
	//	e := echo.New()
	//	req := httptest.NewRequest(http.MethodPost, "/api/user/urls", nil)
	//	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	//	req.AddCookie(cookies.CreateCookie(ctx, cookieName))
	//	rec := httptest.NewRecorder()
	//
	//	cookie := cookies.CreateCookie(ctx, cookieName)
	//
	//	URLSInMem := inmemory.NewInMemory()
	//	_, err := URLSInMem.Put(ctx, dto.Shortening{
	//		CorrelationID: "1",
	//		ShortURL:      "34ewfd",
	//		OriginalURL:   "http://test.com",
	//		UserID:        strconv.Itoa(cookies.GetUser(ctx, cookie.Value).ID)})
	//	if err != nil {
	//		return
	//	}
	//
	//	h := Handler{
	//		ctx:     ctx,
	//		cache:   map[string]*dto.User{cookie.Value: cookies.GetUser(ctx, cookie.Value)},
	//		servAPI: service.NewService(URLSInMem),
	//	}
	//
	//	c := e.NewContext(req, rec)
	//
	//	handler := h.Authorized(dummyHandler)
	//
	//	// Assertions
	//	if assert.NoError(t, handler(c)) {
	//		assert.Equal(t, http.StatusUnauthorized, rec.Code)
	//	}
	//})
}

func TestHandler_HandleUserUrls(t *testing.T) {
	t.Run("Test Get User URLS Ok", func(t *testing.T) {
		// Setup
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/api/user/urls", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		cookie := cookies.CreateCookie(ctx, cookieName)
		URLSInMem := inmemory.NewInMemory()
		_, err := URLSInMem.Put(ctx, dto.Shortening{
			CorrelationID: "1",
			ShortURL:      "34ewfd",
			OriginalURL:   "http://test.com",
			UserID:        strconv.Itoa(cookies.GetUser(ctx, cookie.Value).ID)})
		if err != nil {
			return
		}

		h := Handler{
			ctx:     ctx,
			cache:   make(map[string]*dto.User),
			servAPI: service.NewService(URLSInMem),
		}
		c := &CustomContext{user: cookies.GetUser(ctx, cookie.Value), Context: e.NewContext(req, rec)}

		// Assertions
		if assert.NoError(t, h.HandleUserUrls(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
		}
	})

	t.Run("Test Get User URLS URL Not Found", func(t *testing.T) {
		// Setup
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/api/user/urls", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		h := Handler{
			ctx:     ctx,
			cache:   make(map[string]*dto.User),
			servAPI: service.NewService(inmemory.NewInMemory()),
		}

		cookie := cookies.CreateCookie(ctx, cookieName)
		c := &CustomContext{user: cookies.GetUser(ctx, cookie.Value), Context: e.NewContext(req, rec)}

		// Assertions
		if assert.NoError(t, h.HandleUserUrls(c)) {
			assert.Equal(t, http.StatusBadRequest, rec.Code)
		}
	})
}

func TestHandler_HandleUserUrlsDelete(t *testing.T) {
	t.Run("Test Delete URL for User Ok", func(t *testing.T) {
		// Setup
		e := echo.New()
		req := httptest.NewRequest(
			http.MethodDelete,
			"/api/user/urls",
			strings.NewReader("[\"6qxTVvsy\"]"),
		)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		cookie := cookies.CreateCookie(ctx, cookieName)
		URLSInMem := inmemory.NewInMemory()
		_, err := URLSInMem.Put(ctx, dto.Shortening{
			CorrelationID: "1",
			ShortURL:      "6qxTVvsy",
			OriginalURL:   "http://test.com",
			UserID:        strconv.Itoa(cookies.GetUser(ctx, cookie.Value).ID)})
		if err != nil {
			return
		}

		h := Handler{
			ctx:     ctx,
			cache:   make(map[string]*dto.User),
			servAPI: service.NewService(URLSInMem),
		}
		c := &CustomContext{user: cookies.GetUser(ctx, cookie.Value), Context: e.NewContext(req, rec)}

		// Assertions
		if assert.NoError(t, h.HandleUserUrlsDelete(c)) {
			assert.Equal(t, http.StatusAccepted, rec.Code)
		}
	})

	t.Run("Test Delete URL for User BadRequest", func(t *testing.T) {
		// Setup
		e := echo.New()
		req := httptest.NewRequest(http.MethodDelete, "/api/user/urls", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		h := Handler{
			ctx:     ctx,
			cache:   make(map[string]*dto.User),
			servAPI: service.NewService(inmemory.NewInMemory()),
		}
		c := &CustomContext{user: &dto.User{}, Context: e.NewContext(req, rec)}

		// Assertions
		if assert.NoError(t, h.HandleUserUrlsDelete(c)) {
			assert.Equal(t, http.StatusBadRequest, rec.Code)
		}
	})

	t.Run("Test Delete URL for User InternalServerError", func(t *testing.T) {
		// Setup
		e := echo.New()
		req := httptest.NewRequest(
			http.MethodDelete,
			"/api/user/urls",
			strings.NewReader("[\n    {\n        \"short_url\": \"http://...\",\n        \"original_url\": \"http://...\"\n    },\n    ...\n] "),
		)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		h := Handler{
			ctx:     ctx,
			cache:   make(map[string]*dto.User),
			servAPI: service.NewService(inmemory.NewInMemory()),
		}
		c := &CustomContext{user: &dto.User{}, Context: e.NewContext(req, rec)}

		// Assertions
		if assert.NoError(t, h.HandleUserUrlsDelete(c)) {
			assert.Equal(t, http.StatusInternalServerError, rec.Code)
		}
	})
}

// IP address within the trusted subnet proceeds to the next handler
func TestIPWithinTrustedSubnetProceeds(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	c := e.NewContext(req, rec)
	c.Request().Header.Set(echo.HeaderXRealIP, "192.168.1.1")

	config.GetConfig().TrustedSubnet = "192.168.1.0/24"
	handler := &Handler{}
	nextHandler := func(c echo.Context) error {
		return c.String(http.StatusOK, "next handler called")
	}

	err := handler.FilterIP(nextHandler)(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "next handler called", rec.Body.String())
}

// Trusted subnet is an empty string
func TestEmptyTrustedSubnet(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set(echo.HeaderXRealIP, "192.168.1.1")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	config.GetConfig().TrustedSubnet = ""

	h := &Handler{}
	next := func(c echo.Context) error {
		return c.String(http.StatusOK, "OK")
	}

	err := h.FilterIP(next)(c)

	assert.Error(t, err)
	httpError, ok := err.(*echo.HTTPError)
	assert.True(t, ok)
	assert.Equal(t, http.StatusForbidden, httpError.Code)
}
