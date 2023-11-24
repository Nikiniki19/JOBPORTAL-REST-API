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

func Test_handler_postJob(t *testing.T) {

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
			expectedResponse:   `{"error":"Internal Server Error"}`,
		},
		{
			name: "missing jwt claims",
			setup: func() (*gin.Context, *httptest.ResponseRecorder, services.UserService) {
				rr := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(rr)
				httpRequest, _ := http.NewRequest(http.MethodGet, "http://test.com", nil)
				ctx := httpRequest.Context()
				ctx = context.WithValue(ctx, middlewares.TraceIdKey, "123")
				httpRequest = httpRequest.WithContext(ctx)
				c.Request = httpRequest

				return c, rr, nil
			},
			expectedStatusCode: http.StatusUnauthorized,
			expectedResponse:   `{"error":"Unauthorized"}`,
		},
		{
			name: "id invalid",
			setup: func() (*gin.Context, *httptest.ResponseRecorder, services.UserService) {
				rr := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(rr)
				httpRequest, _ := http.NewRequest(http.MethodGet, "http://test.com:8080", nil)
				ctx := httpRequest.Context()
				ctx = context.WithValue(ctx, middlewares.TraceIdKey, "1")
				ctx = context.WithValue(ctx, auth.Key, jwt.RegisteredClaims{})
				httpRequest = httpRequest.WithContext(ctx)
				c.Params = append(c.Params, gin.Param{Key: "cid", Value: "abc"})
				c.Request = httpRequest
				mc := gomock.NewController(t)
				ms := services.NewMockUserService(mc)
				return c, rr, ms
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   `{"error":"Bad Request"}`,
		},
		{name: "Decode failure",
			setup: func() (*gin.Context, *httptest.ResponseRecorder, services.UserService) {
				rr := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(rr)
				httpRequest, _ := http.NewRequest(http.MethodGet, "http://tests.com", strings.NewReader(`{
				"sal": "85000,
				"minNp": 3,
				"maxNp": 60,
				"budget": 85000.00,
				"jobDesc": "We are hiring a software engineer...",
				"minExp": 2.5,
				"maxExp": 5.5,
				"locationIDs": [1,2],
				"skillIDs": [1],
				"workModeIDs": [1],
				"qualificationIDs": [1],
				"shiftIDs": [1],
				"jobTypeIDs": [1]}`))
				ctx := httpRequest.Context()
				ctx = context.WithValue(ctx, middlewares.TraceIdKey, "1")
				ctx = context.WithValue(ctx, auth.Key, jwt.RegisteredClaims{})
				httpRequest = httpRequest.WithContext(ctx)
				c.Params = append(c.Params, gin.Param{Key: "cid", Value: "1"})
				c.Request = httpRequest
				mc := gomock.NewController(t)
				ms := services.NewMockUserService(mc)
				return c, rr, ms
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   `{"error":"please provide proper data"}`,
		},
		{name: "add job failure",
			setup: func() (*gin.Context, *httptest.ResponseRecorder, services.UserService) {
				rr := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(rr)
				httpRequest, _ := http.NewRequest(http.MethodGet, "http://tests.com", strings.NewReader(`{
					            "jobTitle": "asdfghj",
								"sal": "85000",
				 				"minNp": 3,
				 				"maxNp": 60,
								"budget": 85000,
				 				"jobDesc": "We are hiring a software engineer...",
				 				"minExp": 2.5,
								"maxExp": 5.5,
				 				"locationIDs": [1,2],
								"skillIDs": [1],
				 				"workModeIDs": [1],
				 				"qualificationIDs": [1],
				 				"shiftIDs": [1],
				 				"jobTypeIDs": [1]}`))
				ctx := httpRequest.Context()
				ctx = context.WithValue(ctx, middlewares.TraceIdKey, "1")
				ctx = context.WithValue(ctx, auth.Key, jwt.RegisteredClaims{})
				httpRequest = httpRequest.WithContext(ctx)
				c.Params = append(c.Params, gin.Param{Key: "cid", Value: "1"})
				c.Request = httpRequest
				mc := gomock.NewController(t)
				ms := services.NewMockUserService(mc)
				ms.EXPECT().AddJobDetails(gomock.Any(), gomock.Any(), gomock.Any()).Return(models.Response{}, errors.New("error in adding job")).AnyTimes()
				return c, rr, ms
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   `{"error":"error in adding job"}`,
		},
		{name: "add job success",
			setup: func() (*gin.Context, *httptest.ResponseRecorder, services.UserService) {
				rr := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(rr)
				httpRequest, _ := http.NewRequest(http.MethodGet, "http://tests.com", strings.NewReader(`{
					            "jobTitle": "asdfghj",
								"sal": "85000",
				 				"minNp": 3,
				 				"maxNp": 60,
								"budget": 85000.0,
				 				"jobDesc": "We are hiring a software engineer...",
				 				"minExp": 2.5,
								"maxExp": 5.5,
				 				"locationIDs": [1,2],
								"skillIDs": [1],
				 				"workModeIDs": [1],
				 				"qualificationIDs": [1],
				 				"shiftIDs": [1],
				 				"jobTypeIDs": [1]}`))
				ctx := httpRequest.Context()
				ctx = context.WithValue(ctx, middlewares.TraceIdKey, "1")
				ctx = context.WithValue(ctx, auth.Key, jwt.RegisteredClaims{})
				httpRequest = httpRequest.WithContext(ctx)
				c.Params = append(c.Params, gin.Param{Key: "cid", Value: "1"})
				c.Request = httpRequest
				mc := gomock.NewController(t)
				ms := services.NewMockUserService(mc)
				ms.EXPECT().AddJobDetails(gomock.Any(), gomock.Any(), gomock.Any()).Return(models.Response{}, nil).AnyTimes()
				return c, rr, ms
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse:   `{"ID":0}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gin.SetMode(gin.TestMode)
			c, rr, ms := tt.setup()
			h := &handler{s: ms}
			h.postJob(c)
			assert.Equal(t, tt.expectedStatusCode, rr.Code)
			assert.Equal(t, tt.expectedResponse, rr.Body.String())

		})
	}
}

func Test_handler_getJobsFromCompany(t *testing.T) {
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
		{
			name: "invalid id",
			setup: func() (*gin.Context, *httptest.ResponseRecorder, services.UserService) {
				rr := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(rr)
				httpRequest, _ := http.NewRequest(http.MethodGet, "http://test.com:8080", nil)
				ctx := httpRequest.Context()
				ctx = context.WithValue(ctx, middlewares.TraceIdKey, "10")
				httpRequest = httpRequest.WithContext(ctx)
				c.Request = httpRequest
				c.Params = append(c.Params, gin.Param{Key: "CompanyId", Value: "abc"})
				mc := gomock.NewController(t)
				ms := services.NewMockUserService(mc)
				ms.EXPECT().ViewJobFromCompany(gomock.Any()).Return([]models.Job{}, errors.New("error")).AnyTimes()
				return c, rr, ms
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedResponse:   `{"msg":"Internal Server Error"}`,
		},
		{name: "viewing  a job from company successful",
			setup: func() (*gin.Context, *httptest.ResponseRecorder, services.UserService) {
				rr := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(rr)
				httpRequest, _ := http.NewRequest(http.MethodGet, "http://tests.com", strings.NewReader(`{"job_title":"vnhvgh","sal": "189787"}`))
				ctx := httpRequest.Context()
				ctx = context.WithValue(ctx, middlewares.TraceIdKey, "10")
				httpRequest = httpRequest.WithContext(ctx)
				c.Request = httpRequest
				c.Params = append(c.Params, gin.Param{Key: "CompanyId", Value: "18"})
				mc := gomock.NewController(t)
				ms := services.NewMockUserService(mc)
				ms.EXPECT().ViewJobFromCompany(gomock.Any()).Return([]models.Job{}, nil).AnyTimes()
				return c, rr, ms
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse:   `[]`,
		},
		{name: "viewing  a job from company failed",
			setup: func() (*gin.Context, *httptest.ResponseRecorder, services.UserService) {
				rr := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(rr)
				httpRequest, _ := http.NewRequest(http.MethodGet, "http://tests.com", strings.NewReader(`{"job_title":"vnhvgh","sal": "189787"}`))
				ctx := httpRequest.Context()
				ctx = context.WithValue(ctx, middlewares.TraceIdKey, "10")
				httpRequest = httpRequest.WithContext(ctx)
				c.Request = httpRequest
				c.Params = append(c.Params, gin.Param{Key: "CompanyId", Value: "18"})
				mc := gomock.NewController(t)
				ms := services.NewMockUserService(mc)
				ms.EXPECT().ViewJobFromCompany(gomock.Any()).Return([]models.Job{}, errors.New("error")).AnyTimes()
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
			h.getJobsFromCompany(c)
			assert.Equal(t, tt.expectedStatusCode, rr.Code)
			assert.Equal(t, tt.expectedResponse, rr.Body.String())
		})
	}
}

func Test_handler_getAllJobs(t *testing.T) {
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
		{
			name: "success in viewing all jobs",
			setup: func() (*gin.Context, *httptest.ResponseRecorder, services.UserService) {
				rr := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(rr)
				httpRequest, _ := http.NewRequest(http.MethodGet, "http://test.com:8080", nil)
				ctx := httpRequest.Context()
				ctx = context.WithValue(ctx, middlewares.TraceIdKey, "123")
				httpRequest = httpRequest.WithContext(ctx)
				c.Request = httpRequest
				mc := gomock.NewController(t)
				ms := services.NewMockUserService(mc)
				ms.EXPECT().ViewAllJobs(gomock.Any()).Return([]models.Job{}, nil).AnyTimes()
				return c, rr, ms
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse:   `[]`,
		},
		{
			name: "failed in viewing all jobs",
			setup: func() (*gin.Context, *httptest.ResponseRecorder, services.UserService) {
				rr := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(rr)
				httpRequest, _ := http.NewRequest(http.MethodGet, "http://test.com:8080", nil)
				ctx := httpRequest.Context()
				ctx = context.WithValue(ctx, middlewares.TraceIdKey, "123")
				httpRequest = httpRequest.WithContext(ctx)
				c.Request = httpRequest
				mc := gomock.NewController(t)
				ms := services.NewMockUserService(mc)
				ms.EXPECT().ViewAllJobs(gomock.Any()).Return([]models.Job{}, errors.New("error")).AnyTimes()
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
			h.getAllJobs(c)
			assert.Equal(t, tt.expectedStatusCode, rr.Code)
			assert.Equal(t, tt.expectedResponse, rr.Body.String())
		})
	}
}

func Test_handler_getOneJob(t *testing.T) {
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
		{
			name: "invalid company id ",
			setup: func() (*gin.Context, *httptest.ResponseRecorder, services.UserService) {
				rr := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(rr)
				httpRequest, _ := http.NewRequest(http.MethodGet, "http://test.com:8080", nil)
				ctx := httpRequest.Context()
				ctx = context.WithValue(ctx, middlewares.TraceIdKey, "193")
				httpRequest = httpRequest.WithContext(ctx)
				c.Request = httpRequest
				c.Params = append(c.Params, gin.Param{Key: "CompanyId", Value: "abc"})

				return c, rr, nil
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   `{"msg":"Bad Request"}`,
		},
		{
			name: "success in fetching jobs",
			setup: func() (*gin.Context, *httptest.ResponseRecorder, services.UserService) {
				rr := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(rr)
				httpRequest, _ := http.NewRequest(http.MethodGet, "http://test.com:8080", nil)
				ctx := httpRequest.Context()
				ctx = context.WithValue(ctx, middlewares.TraceIdKey, "123")

				httpRequest = httpRequest.WithContext(ctx)
				c.Request = httpRequest
				c.Params = append(c.Params, gin.Param{Key: "CompanyId", Value: "9"})
				mc := gomock.NewController(t)
				ms := services.NewMockUserService(mc)
				ms.EXPECT().ViewJobById(gomock.Any(), gomock.Any()).Return([]models.Job{}, nil).AnyTimes()

				return c, rr, ms
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse:   `[]`,
		},
		{
			name: "failure in fetching jobs",
			setup: func() (*gin.Context, *httptest.ResponseRecorder, services.UserService) {
				rr := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(rr)
				httpRequest, _ := http.NewRequest(http.MethodGet, "http://test.com:8080", nil)
				ctx := httpRequest.Context()
				ctx = context.WithValue(ctx, middlewares.TraceIdKey, "123")

				httpRequest = httpRequest.WithContext(ctx)
				c.Request = httpRequest
				c.Params = append(c.Params, gin.Param{Key: "CompanyId", Value: "9"})
				mc := gomock.NewController(t)
				ms := services.NewMockUserService(mc)
				ms.EXPECT().ViewJobById(gomock.Any(), gomock.Any()).Return([]models.Job{}, errors.New("error")).AnyTimes()

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
			h.getOneJob(c)
			assert.Equal(t, tt.expectedStatusCode, rr.Code)
			assert.Equal(t, tt.expectedResponse, rr.Body.String())
		})
	}
}

func Test_handler_processApplications(t *testing.T) {
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
		{
			name: "missing jwt claims",
			setup: func() (*gin.Context, *httptest.ResponseRecorder, services.UserService) {
				rr := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(rr)
				httpRequest, _ := http.NewRequest(http.MethodGet, "http://test.com", nil)
				ctx := httpRequest.Context()
				ctx = context.WithValue(ctx, middlewares.TraceIdKey, "123")
				httpRequest = httpRequest.WithContext(ctx)
				c.Request = httpRequest

				return c, rr, nil
			},
			expectedStatusCode: http.StatusUnauthorized,
			expectedResponse:   `{"error":"Unauthorized"}`,
		},
		{name: "Decode failure",
			setup: func() (*gin.Context, *httptest.ResponseRecorder, services.UserService) {
				rr := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(rr)
				httpRequest, _ := http.NewRequest(http.MethodGet, "http://tests.com", strings.NewReader(`{
					[
						{
							"name": "niki",
							"age": "30",
							"jid": 2,
							"job_application": {
								"noticePeriod": 0,
								"location": [
									1
								],
								"technologyStack": [
									1
								],
								"experience": 2.5,
								"qualifications": [
									1
								],
								"shifts": [
									2
								],
								"work_modes": [
									1
								],
								"job_type": [
									1,
									2,
									3
								]
							}
						},
						{
							"name": "bhoomika",
							"age": "22",
							"jid": 1,
							"job_application": {
								"noticePeriod": 75,
								"location": [
									1
								],
								"technologyStack": [
									2
								],
								"experience": 2.5,
								"qualifications": [
									2
								],
								"shifts": [
									3
								],
								"work_modes": [
									2
								],
								"job_type": [
									2
								]
							}
						},
						
					]}`))
				ctx := httpRequest.Context()
				ctx = context.WithValue(ctx, middlewares.TraceIdKey, "1")
				ctx = context.WithValue(ctx, auth.Key, jwt.RegisteredClaims{})
				httpRequest = httpRequest.WithContext(ctx)
				c.Params = append(c.Params, gin.Param{Key: "cid", Value: "1"})
				c.Request = httpRequest
				return c, rr, nil
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   `{"error":"Bad Request"}`,
		},
		{name: "error in request validation",
			setup: func() (*gin.Context, *httptest.ResponseRecorder, services.UserService) {
				rr := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(rr)
				httpRequest, _ := http.NewRequest(http.MethodGet, "http://tests.com", strings.NewReader(`{
								[
									{
										"name": "niki",
										"age": "30",
										"jid": 2,
										"job_application": {
											"noticePeriod": ,
											"location": [
												1
											],
											"technologyStack": [
												1
											],
											"experience": ,
											"qualifications": [
												1
											],
											"shifts": [
												2
											],
											"work_modes": [
												1
											],
											"job_type": [
												1,
												2,
												3
											]
										}
									},
									{
										"name": "bhoomika",
										"age": "22",
										"jid": 1,
										"job_application": {
											"noticePeriod": 75,
											"location": [
												1
											],
											"technologyStack": [
												2
											],
											"experience": 2.5,
											"qualifications": [
												2
											],
											"shifts": [
												3
											],
											"work_modes": [
												2
											],
											"job_type": [
												2
											]
										}
									},
									
								]}`))
				ctx := httpRequest.Context()
				ctx = context.WithValue(ctx, middlewares.TraceIdKey, "7")
				ctx = context.WithValue(ctx, auth.Key, jwt.RegisteredClaims{})
				httpRequest = httpRequest.WithContext(ctx)
				c.Request = httpRequest
				//c.Params = append(c.Params, gin.Param{Key: "id", Value: " abc"})
				return c, rr, nil
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   `{"error":"Bad Request"}`,
		},
		// {
		// 	name:"error in validation",
		// 	setup: func() (*gin.Context, *httptest.ResponseRecorder, services.UserService) {
		// 		rr := httptest.NewRecorder()
		// 		c, _ := gin.CreateTestContext(rr)
		// 		httpRequest, _ := http.NewRequest(http.MethodGet, "http://tests.com", strings.NewReader(`{
		// 			[
		// 				{
		// 					"name": "niki",
		// 					"age": "30",
		// 					"jid": 2,
		// 					"job_application": {
		// 						"noticePeriod": ,
		// 						"location": [
		// 							1
		// 						],
		// 						"technologyStack": [
		// 							1
		// 						],
		// 						"experience": ,
		// 						"qualifications": [
		// 							1
		// 						],
		// 						"shifts": [
		// 							2
		// 						],
		// 						"work_modes": [
		// 							1
		// 						],
		// 						"job_type": [
		// 							1,
		// 							2,
		// 							3
		// 						]
		// 					}
		// 				},
		// 				{
		// 					"name": "bhoomika",
		// 					"age": "22",
		// 					"jid": 1,
		// 					"job_application": {
		// 						"noticePeriod": 75,
		// 						"location": [
		// 							1
		// 						],
		// 						"technologyStack": [
		// 							2
		// 						],
		// 						"experience": 2.5,
		// 						"qualifications": [
		// 							2
		// 						],
		// 						"shifts": [
		// 							3
		// 						],
		// 						"work_modes": [
		// 							2
		// 						],
		// 						"job_type": [
		// 							2
		// 						]
		// 					}
		// 				},
						
		// 			]}`))
		// 		ctx := httpRequest.Context()
		// 		ctx = context.WithValue(ctx, middlewares.TraceIdKey, "1")
		// 		ctx = context.WithValue(ctx, auth.Key, jwt.RegisteredClaims{})
		// 		httpRequest = httpRequest.WithContext(ctx)
		// 		c.Request = httpRequest
		// 		return c, rr, nil
		// 	},
		// 	expectedStatusCode: http.StatusBadRequest,
		// 	expectedResponse:  `{"error":"Bad Request"}`,
		// },
		{
			name: "process application success",
			setup: func() (*gin.Context, *httptest.ResponseRecorder, services.UserService) {
				rr := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(rr)
				httpRequest, _ := http.NewRequest(http.MethodGet, "http://tests.com", strings.NewReader(`
					[
						{
							"name": "some name",
							"age": "21",
							"jid": 1,
							"job_application": {
								"noticePeriod": 25,
								"location": [
									1
								],
								"technologyStack": [
									1
								],
								"experience": 2.5,
								"qualifications": [
									1
								],
								"shifts": [
									2
								],
								"work_modes": [
									1
								],
								"job_type": [
									1
								]
							}
						}
					]`))
				ctx := httpRequest.Context()
				ctx = context.WithValue(ctx, middlewares.TraceIdKey, "1")
				ctx = context.WithValue(ctx, auth.Key, jwt.RegisteredClaims{})
				httpRequest = httpRequest.WithContext(ctx)
				c.Request = httpRequest
				mc := gomock.NewController(t)
				ms := services.NewMockUserService(mc)
				ms.EXPECT().ProcessJobApplications(gomock.Any()).Return([]models.NewUserApplication{}, nil).AnyTimes()
				return c, rr, ms
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse:   `[]`,
		},
		{
			name: "process application failure",
			setup: func() (*gin.Context, *httptest.ResponseRecorder, services.UserService) {
				rr := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(rr)
				httpRequest, _ := http.NewRequest(http.MethodGet, "http://tests.com", strings.NewReader(`
					[
						{
							"name": "some name",
							"age": "21",
							"jid": 1,
							"job_application": {
								"noticePeriod": 25,
								"location": [
									1
								],
								"technologyStack": [
									1
								],
								"experience": 2.5,
								"qualifications": [
									1
								],
								"shifts": [
									2
								],
								"work_modes": [
									1
								],
								"job_type": [
									1
								]
							}
						}
					]`))
				ctx := httpRequest.Context()
				ctx = context.WithValue(ctx, middlewares.TraceIdKey, "1")
				ctx = context.WithValue(ctx, auth.Key, jwt.RegisteredClaims{})
				httpRequest = httpRequest.WithContext(ctx)
				c.Request = httpRequest
				mc := gomock.NewController(t)
				ms := services.NewMockUserService(mc)
				ms.EXPECT().ProcessJobApplications(gomock.Any()).Return([]models.NewUserApplication{}, errors.New("error")).AnyTimes()
				return c, rr, ms
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedResponse:   `{"error":"Internal Server Error"}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gin.SetMode(gin.TestMode)
			c, rr, ms := tt.setup()
			h := &handler{s: ms}
			h.processApplications(c)
			assert.Equal(t, tt.expectedStatusCode, rr.Code)
			assert.Equal(t, tt.expectedResponse, rr.Body.String())
		})
	}
}
