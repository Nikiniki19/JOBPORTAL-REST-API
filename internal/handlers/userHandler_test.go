package handlers

import (
	"context"
	"errors"
	"job-portal-api/internal/auth"
	"job-portal-api/internal/middlewares"
	"job-portal-api/internal/models"
	"job-portal-api/internal/services"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/mock/gomock"
)

func Test_handler_Registration(t *testing.T) {
	tests := []struct {
		name               string
		setup              func() (*gin.Context, *httptest.ResponseRecorder, services.UserService)
		expectedStatusCode int
		expectedResponse   string
	}{
		{name: "traceid missing from context",
			setup: func() (*gin.Context, *httptest.ResponseRecorder, services.UserService) {
				rr := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(rr)
				httpRequest, _ := http.NewRequest(http.MethodGet, "http://tests.com", nil)
				c.Request = httpRequest
				return c, rr, nil
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedResponse:   `{"msg":"Internal Server Error"}`,
		},

		{name: "request validation success",
			setup: func() (*gin.Context, *httptest.ResponseRecorder, services.UserService) {
				rr := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(rr)
				httpRequest, _ := http.NewRequest(http.MethodGet, "http://tests.com", strings.NewReader(`{"name":"nikitha","email":"niki@gmail.com","password":"1234"}`))
				c.Request = httpRequest
				return c, rr, nil
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedResponse:   `{"msg":"Internal Server Error"}`,
		},
		{name: "error in request validation",
			setup: func() (*gin.Context, *httptest.ResponseRecorder, services.UserService) {
				rr := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(rr)
				httpRequest, _ := http.NewRequest(http.MethodGet, "http://tests.com", strings.NewReader(`{"name":"","email":"niki@gmail.com","password":"1234"}`))
				ctx := httpRequest.Context()
				ctx = context.WithValue(ctx, middlewares.TraceIdKey, "7")
				httpRequest = httpRequest.WithContext(ctx)
				c.Request = httpRequest
				c.Params = append(c.Params, gin.Param{Key: "id", Value: " abc"})
				return c, rr, nil
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   `{"msg":"Bad Request"}`,
		},

		{name: "error in decoding",
			setup: func() (*gin.Context, *httptest.ResponseRecorder, services.UserService) {
				rr := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(rr)
				httpRequest, _ := http.NewRequest(http.MethodGet, "http://tests.com", strings.NewReader(`
				{"name":" nikitha","email":"niki@gmail.com","password":"1234}`))
				ctx := httpRequest.Context()
				ctx = context.WithValue(ctx, middlewares.TraceIdKey, "1")
				httpRequest = httpRequest.WithContext(ctx)
				c.Request = httpRequest
				c.Params = append(c.Params, gin.Param{Key: "id", Value: "1"})
				return c, rr, nil
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedResponse:   `{"msg":"Internal Server Error"}`,
		},

		{name: "registration successful",
			setup: func() (*gin.Context, *httptest.ResponseRecorder, services.UserService) {
				rr := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(rr)
				httpRequest, _ := http.NewRequest(http.MethodGet, "http://tests.com", strings.NewReader(`{"name":"nikitha","email":"niki@gmail.com","password":"1234"}`))
				ctx := httpRequest.Context()
				ctx = context.WithValue(ctx, middlewares.TraceIdKey, "1")
				httpRequest = httpRequest.WithContext(ctx)
				c.Request = httpRequest
				c.Params = append(c.Params, gin.Param{Key: "id", Value: "1"})
				mc := gomock.NewController(t)
				ms := services.NewMockUserService(mc)
				ms.EXPECT().Signup(gomock.Any(), gomock.Any()).Return(models.User{}, nil).AnyTimes()
				return c, rr, ms
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse:   `{"ID":0,"CreatedAt":"0001-01-01T00:00:00Z","UpdatedAt":"0001-01-01T00:00:00Z","DeletedAt":null,"name":"","email":""}`,
		},
		{name: "registration failure",
			setup: func() (*gin.Context, *httptest.ResponseRecorder, services.UserService) {
				rr := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(rr)
				httpRequest, _ := http.NewRequest(http.MethodGet, "http://tests.com", strings.NewReader(`{"name":"nikitha","email":"niki@gmail.com","password":"1234"}`))
				ctx := httpRequest.Context()
				ctx = context.WithValue(ctx, middlewares.TraceIdKey, "1")
				httpRequest = httpRequest.WithContext(ctx)
				c.Request = httpRequest
				c.Params = append(c.Params, gin.Param{Key: "id", Value: "1"})
				mc := gomock.NewController(t)
				ms := services.NewMockUserService(mc)
				ms.EXPECT().Signup(gomock.Any(), gomock.Any()).Return(models.User{}, errors.New("error")).AnyTimes()
				return c, rr, ms
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   `{"error":"Bad Request"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gin.SetMode(gin.TestMode)
			c, rr, ms := tt.setup()
			h := &handler{s: ms}
			h.Registration(c)
			assert.Equal(t, tt.expectedStatusCode, rr.Code)
			assert.Equal(t, tt.expectedResponse, rr.Body.String())

		})
	}
}

func Test_handler_Login(t *testing.T) {

	tests := []struct {
		name               string
		setup              func() (*gin.Context, *httptest.ResponseRecorder, services.UserService, auth.Authentication)
		expectedStatusCode int
		expectedResponse   string
	}{
		{name: "traceid missing from context",
			setup: func() (*gin.Context, *httptest.ResponseRecorder, services.UserService, auth.Authentication) {
				rr := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(rr)
				httpRequest, _ := http.NewRequest(http.MethodGet, "http://tests.com", nil)
				c.Request = httpRequest
				return c, rr, nil, nil
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedResponse:   `{"msg":"Internal Server Error"}`,
		},
		{name: "error in decoding",
			setup: func() (*gin.Context, *httptest.ResponseRecorder, services.UserService, auth.Authentication) {
				rr := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(rr)
				httpRequest, _ := http.NewRequest(http.MethodGet, "http://tests.com", strings.NewReader(`
				{"name":" nikitha","email":"niki@gmail.com","password":"1234}`))
				ctx := httpRequest.Context()
				ctx = context.WithValue(ctx, middlewares.TraceIdKey, "1")
				httpRequest = httpRequest.WithContext(ctx)
				c.Request = httpRequest
				c.Params = append(c.Params, gin.Param{Key: "id", Value: "1"})
				return c, rr, nil, nil
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedResponse:   `{"msg":"Internal Server Error"}`,
		},

		{name: "request validation",
			setup: func() (*gin.Context, *httptest.ResponseRecorder, services.UserService, auth.Authentication) {
				rr := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(rr)
				httpRequest, _ := http.NewRequest(http.MethodGet, "http://tests.com", strings.NewReader(`{"email":"niki@gmail.com","password":"1234"}`))
				c.Request = httpRequest
				return c, rr, nil, nil
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedResponse:   `{"msg":"Internal Server Error"}`,
		},
		{name: "error in request validation",
			setup: func() (*gin.Context, *httptest.ResponseRecorder, services.UserService, auth.Authentication) {
				rr := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(rr)
				httpRequest, _ := http.NewRequest(http.MethodGet, "http://tests.com", strings.NewReader(`{"email":" ","password":"1234"}`))
				ctx := httpRequest.Context()
				ctx = context.WithValue(ctx, middlewares.TraceIdKey, "1")
				httpRequest = httpRequest.WithContext(ctx)
				c.Request = httpRequest
				c.Params = append(c.Params, gin.Param{Key: "id", Value: " 1"})
				return c, rr, nil, nil
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   `{"msg":"please provide Email and Password"}`,
		},

		{name: "login successful",
			setup: func() (*gin.Context, *httptest.ResponseRecorder, services.UserService, auth.Authentication) {
				rr := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(rr)
				httpRequest, _ := http.NewRequest(http.MethodGet, "http://tests.com", strings.NewReader(`{"name":"nikitha","email":"niki@gmail.com","password":"1234"}`))
				ctx := httpRequest.Context()
				ctx = context.WithValue(ctx, middlewares.TraceIdKey, "1")
				httpRequest = httpRequest.WithContext(ctx)
				c.Request = httpRequest
				mc := gomock.NewController(t)
				ms := services.NewMockUserService(mc)
				ma := auth.NewMockAuthentication(mc)

				ms.EXPECT().Login(gomock.Any(), gomock.Any(), gomock.Any()).Return(jwt.RegisteredClaims{}, nil)
				ma.EXPECT().GenerateToken(gomock.Any()).Return("", nil)

				return c, rr, ms, ma
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse:   `{"token":""}`,
		},
		{name: "login failure",
			setup: func() (*gin.Context, *httptest.ResponseRecorder, services.UserService, auth.Authentication) {
				rr := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(rr)
				httpRequest, _ := http.NewRequest(http.MethodGet, "http://tests.com", strings.NewReader(`{"email":"werty@gmail.com","password":"1234"}`))
				ctx := httpRequest.Context()
				ctx = context.WithValue(ctx, middlewares.TraceIdKey, "1")
				httpRequest = httpRequest.WithContext(ctx)
				c.Request = httpRequest
				mc := gomock.NewController(t)
				ms := services.NewMockUserService(mc)
				ms.EXPECT().Login(gomock.Any(), gomock.Any(), gomock.Any()).Return(jwt.RegisteredClaims{}, errors.New("error"))

				return c, rr, ms, nil
			},
			expectedStatusCode: http.StatusUnauthorized,
			expectedResponse:   `{"msg":"login failed"}`,
		},
		{name: " failure in generating token",
			setup: func() (*gin.Context, *httptest.ResponseRecorder, services.UserService, auth.Authentication) {
				rr := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(rr)
				httpRequest, _ := http.NewRequest(http.MethodGet, "http://tests.com", strings.NewReader(`{"email":"qwerty@gmail.com","password":"1234"}`))
				ctx := httpRequest.Context()
				ctx = context.WithValue(ctx, middlewares.TraceIdKey, "1")
				httpRequest = httpRequest.WithContext(ctx)
				c.Request = httpRequest
				mc := gomock.NewController(t)
				ms := services.NewMockUserService(mc)
				ma := auth.NewMockAuthentication(mc)
				ms.EXPECT().Login(gomock.Any(), gomock.Any(), gomock.Any()).Return(jwt.RegisteredClaims{}, nil)
				ma.EXPECT().GenerateToken(gomock.Any()).Return("", errors.New("error in generating token"))
				return c, rr, ms, ma
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedResponse:   `{"msg":"Internal Server Error"}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gin.SetMode(gin.TestMode)
			c, rr, ms, ma := tt.setup()
			h := &handler{s: ms, a: ma}
			h.Signin(c)
			assert.Equal(t, tt.expectedStatusCode, rr.Code)
			assert.Equal(t, tt.expectedResponse, rr.Body.String())

		})
	}
}
