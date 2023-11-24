package handlers

import (
	"context"
	"errors"

	"job-portal-api/internal/middlewares"
	"job-portal-api/internal/models"
	"job-portal-api/internal/services"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
	"go.uber.org/mock/gomock"
)

func Test_handler_createCom(t *testing.T) {
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
			expectedResponse:   `{"msg1":"Internal Server Error"}`,
		},
		{
			name: "error in decoding",
			setup: func() (*gin.Context, *httptest.ResponseRecorder, services.UserService) {
				rr := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(rr)
				httpRequest, _ := http.NewRequest(http.MethodGet, "http://tests.com", strings.NewReader(`
				{"company_name":" jndsjhd",
				"address":"bangalore",
				"domain":"software}`))
				ctx := httpRequest.Context()
				ctx = context.WithValue(ctx, middlewares.TraceIdKey, "1")
				httpRequest = httpRequest.WithContext(ctx)
				c.Request = httpRequest
				return c, rr, nil
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedResponse:   `{"msg2":"Internal Server Error"}`,
		},
		{name: "error in request validation",
		setup: func() (*gin.Context, *httptest.ResponseRecorder, services.UserService) {
			rr := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rr)
			httpRequest, _ := http.NewRequest(http.MethodGet, "http://tests.com", strings.NewReader(`{"company_name":"qjjqj","address":"niki@gmail.com","password":"1234"}`))
			ctx := httpRequest.Context()
			ctx = context.WithValue(ctx, middlewares.TraceIdKey, "1")
			httpRequest = httpRequest.WithContext(ctx)
			c.Request = httpRequest
			c.Params = append(c.Params, gin.Param{Key: "id", Value: " 1"})
			return c, rr, nil
		},
		expectedStatusCode: http.StatusBadRequest,
		expectedResponse:   `{"msg3":"Bad Request"}`,
	},
	
		{name: "company creation successful",
			setup: func() (*gin.Context, *httptest.ResponseRecorder, services.UserService) {
				rr := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(rr)
				httpRequest, _ := http.NewRequest(http.MethodGet, "http://tests.com", strings.NewReader(`{"company_name":"tek","address":"bangalore","domain":"software"}`))
				ctx := httpRequest.Context()
				ctx = context.WithValue(ctx, middlewares.TraceIdKey, "1")
				httpRequest = httpRequest.WithContext(ctx)
				c.Request = httpRequest
				c.Params = append(c.Params, gin.Param{Key: "id", Value: "1"})
				mc := gomock.NewController(t)
				ms := services.NewMockUserService(mc)
				ms.EXPECT().AddCompanyDetails(gomock.Any(), gomock.Any()).Return(models.Company{}, nil).AnyTimes()
				return c, rr, ms
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse:   `{"ID":0,"CreatedAt":"0001-01-01T00:00:00Z","UpdatedAt":"0001-01-01T00:00:00Z","DeletedAt":null,"company_name":"","address":"","domain":""}`,
		},
		{name: "company creation failed",
			setup: func() (*gin.Context, *httptest.ResponseRecorder, services.UserService) {
				rr := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(rr)
				httpRequest, _ := http.NewRequest(http.MethodGet, "http://tests.com", strings.NewReader(`{"company_name":"tek","address":"bangalore","domain":"software"}`))
				ctx := httpRequest.Context()
				ctx = context.WithValue(ctx, middlewares.TraceIdKey, "1")
				httpRequest = httpRequest.WithContext(ctx)
				c.Request = httpRequest
				mc := gomock.NewController(t)
				ms := services.NewMockUserService(mc)
				ms.EXPECT().AddCompanyDetails(gomock.Any(), gomock.Any()).Return(models.Company{}, errors.New("error in company creation")).AnyTimes()
				return c, rr, ms
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:  `{"msg4":"user not found"}` ,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gin.SetMode(gin.TestMode)
			c, rr, ms := tt.setup()
			h := &handler{s: ms}
			h.createCom(c)
			assert.Equal(t, tt.expectedStatusCode, rr.Code)
			assert.Equal(t, tt.expectedResponse, rr.Body.String())
		})
	}
}

func Test_handler_getAllTheCompanies(t *testing.T) {

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
		{name: "viewing  all companies successful",
			setup: func() (*gin.Context, *httptest.ResponseRecorder, services.UserService) {
				rr := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(rr)
				httpRequest, _ := http.NewRequest(http.MethodGet, "http://tests.com", strings.NewReader(`{"company_name":"tek","address":"bangalore","domain":"software"}`))
				ctx := httpRequest.Context()
				ctx = context.WithValue(ctx, middlewares.TraceIdKey, "1")
				httpRequest = httpRequest.WithContext(ctx)
				c.Request = httpRequest
				c.Params = append(c.Params, gin.Param{Key: "id", Value: "1"})
				mc := gomock.NewController(t)
				ms := services.NewMockUserService(mc)
				ms.EXPECT().ViewAllCompanies(gomock.Any()).Return([]models.Company{}, nil).AnyTimes()
				return c, rr, ms
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse:   "[]",
		},
		{name: "viewing  all companies failure",
			setup: func() (*gin.Context, *httptest.ResponseRecorder, services.UserService) {
				rr := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(rr)
				httpRequest, _ := http.NewRequest(http.MethodGet, "http://tests.com", strings.NewReader(`{"company_name":"tek","address":"bangalore","domain":"software"}`))
				ctx := httpRequest.Context()
				ctx = context.WithValue(ctx, middlewares.TraceIdKey, "1")
				httpRequest = httpRequest.WithContext(ctx)
				c.Request = httpRequest
				c.Params = append(c.Params, gin.Param{Key: "id", Value: "1"})
				mc := gomock.NewController(t)
				ms := services.NewMockUserService(mc)
				ms.EXPECT().ViewAllCompanies(gomock.Any()).Return([]models.Company{}, errors.New("error in viewing company")).AnyTimes()
				return c, rr, ms
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedResponse:   `{"msg":"Internal Server Error"}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gin.SetMode(gin.TestMode)
			c, rr, ms := tt.setup()
			h := &handler{s: ms}
			h.getAllTheCompanies(c)
			assert.Equal(t, tt.expectedStatusCode, rr.Code)
			assert.Equal(t, tt.expectedResponse, rr.Body.String())
		})
	}
}

func Test_handler_viewCompany(t *testing.T) {
	
	tests := []struct {
		name string
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
			expectedResponse:   `{"error":"Internal Server Error"}`,
		},
		{
			name: "id invalid",
			setup: func() (*gin.Context, *httptest.ResponseRecorder, services.UserService) {
				rr := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(rr)
				httpRequest, _ := http.NewRequest(http.MethodGet, "http://test.com:8080", nil)
				ctx := httpRequest.Context()
				ctx = context.WithValue(ctx, middlewares.TraceIdKey, "1")
				httpRequest = httpRequest.WithContext(ctx)
				c.Params = append(c.Params, gin.Param{Key: "cid", Value: "abc"})
				c.Request = httpRequest
				mc := gomock.NewController(t)
				ms := services.NewMockUserService(mc)
				return c, rr, ms
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   `{"msg":"Bad Request"}`,
		},
		{name: "viewing  a company successful",
			setup: func() (*gin.Context, *httptest.ResponseRecorder, services.UserService) {
				rr := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(rr)
				httpRequest, _ := http.NewRequest(http.MethodGet, "http://tests.com", strings.NewReader(`{"company_name":"tek","address":"bangalore","domain":"software"}`))
				ctx := httpRequest.Context()
				ctx = context.WithValue(ctx, middlewares.TraceIdKey, "1")
				httpRequest = httpRequest.WithContext(ctx)
				c.Request = httpRequest
				c.Params = append(c.Params, gin.Param{Key: "cid", Value: "1"})
				mc := gomock.NewController(t)
				ms := services.NewMockUserService(mc)
				ms.EXPECT().ViewCompanyDetails(gomock.Any(),gomock.Any()).Return(models.Company{}, nil).AnyTimes()
				return c, rr, ms
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse:    `{"ID":0,"CreatedAt":"0001-01-01T00:00:00Z","UpdatedAt":"0001-01-01T00:00:00Z","DeletedAt":null,"company_name":"","address":"","domain":""}`,
		},
		{name: "viewing  a company failure",
			setup: func() (*gin.Context, *httptest.ResponseRecorder, services.UserService) {
				rr := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(rr)
				httpRequest, _ := http.NewRequest(http.MethodGet, "http://tests.com", strings.NewReader(`{"company_name":"tek","address":"bangalore","domain":"software"}`))
				ctx := httpRequest.Context()
				ctx = context.WithValue(ctx, middlewares.TraceIdKey, "1")
				httpRequest = httpRequest.WithContext(ctx)
				c.Request = httpRequest
				c.Params = append(c.Params, gin.Param{Key: "cid", Value: "1"})
				mc := gomock.NewController(t)
				ms := services.NewMockUserService(mc)
				ms.EXPECT().ViewCompanyDetails(gomock.Any(),gomock.Any()).Return(models.Company{}, errors.New("errors")).AnyTimes()
				return c, rr, ms
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:    `{"msg":"Bad Request"}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gin.SetMode(gin.TestMode)
			c, rr, ms := tt.setup()
			h := &handler{s: ms}
			h.viewCompany(c)
			assert.Equal(t, tt.expectedStatusCode, rr.Code)
			assert.Equal(t, tt.expectedResponse, rr.Body.String())
		})
	}
}
