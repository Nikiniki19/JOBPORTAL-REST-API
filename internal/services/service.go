package services

import (
	"context"
	"errors"
	"job-portal-api/internal/auth"
	"job-portal-api/internal/caching"

	"job-portal-api/internal/models"
	"job-portal-api/internal/repository"

	"github.com/golang-jwt/jwt/v5"
)

//go:generate mockgen -source=service.go -destination=service_mock.go -package=services

type UserService interface {
	Signup(ctx context.Context, userData models.NewUser) (models.User, error)
	Login(ctx context.Context, email, password string) (jwt.RegisteredClaims, error)

	AddCompanyDetails(ctx context.Context, companyData models.Company) (models.Company, error)
	ViewAllCompanies(ctx context.Context) ([]models.Company, error)
	ViewCompanyDetails(ctx context.Context, id uint64) (models.Company, error)

	ViewJobFromCompany(cid uint64) ([]models.Job, error)
	AddJobDetails(ctx context.Context, jobData models.NewJobRequest, cid uint64) (models.Response, error)
	ViewAllJobs(ctx context.Context) ([]models.Job, error)
	ViewJobById(ctx context.Context, jid uint64) ([]models.Job, error)

	ProcessJobApplications(appData []models.NewUserApplication) ([]models.NewUserApplication, error)
	OTPGeneration(ctx context.Context, data models.ForgotPassword) (string, error)
	ChangePassword(ctx context.Context, otp models.OtpPassword) (string, error)
}

type Service struct {
	UserRepo repository.UserRepo
	auth     auth.Authentication
	UserService
	rdb caching.Cache
}

func NewService(userRepo repository.UserRepo, a auth.Authentication, rdb caching.Cache) (UserService, error) {
	if userRepo == nil {
		return nil, errors.New("interface cannot be nil")
	}
	return &Service{
		UserRepo: userRepo,
		auth:     a,
		rdb:      rdb,
	}, nil
}
