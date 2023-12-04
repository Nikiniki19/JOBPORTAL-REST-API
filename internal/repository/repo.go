package repository

import (
	"context"
	"errors"
	"job-portal-api/internal/models"

	"gorm.io/gorm"
)

//go:generate mockgen -source=repo.go -destination=repo_mock.go -package=repository
type Repo struct {
	DB *gorm.DB
}
type UserRepo interface {
	CreateUser(ctx context.Context, nu models.User) (models.User, error)
	CheckEmail(ctx context.Context, email string) (models.User, error)

	PostJob(nj models.Job) (models.Response, error)
	GetJobsFromCompany(comapny_id uint64) ([]models.Job, error)
	GetAllJobs() ([]models.Job, error)
	GetOneJob(id uint64) ([]models.Job, error)

	CreateCom(nc models.Company) (models.Company, error)
	GetAllTheCompanies() ([]models.Company, error)
	GetCompany(id uint64) (models.Company, error)

	FetchJobData(jid uint64) (models.Job, error)
	UpdatePwdInDb(user models.User)error
}

func NewRepository(DB *gorm.DB) (UserRepo, error) {

	if DB == nil {
		return nil, errors.New("db cannot be nil")

	}

	return &Repo{
		DB: DB,
	}, nil
}
