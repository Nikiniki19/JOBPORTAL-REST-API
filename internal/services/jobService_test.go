package services

import (
	"context"
	"errors"
	"job-portal-api/internal/auth"
	"job-portal-api/internal/caching"
	"job-portal-api/internal/models"
	"job-portal-api/internal/repository"
	"reflect"
	"testing"

	"go.uber.org/mock/gomock"
)

func TestService_ViewJobFromCompany(t *testing.T) {
	type args struct {
		cid uint64
	}
	tests := []struct {
		name string
		//s       *Service
		args             args
		want             []models.Job
		wantErr          bool
		mockRepoResponse func() ([]models.Job, error)
	}{
		{name: "success if jobs are retrieved from cid",
			args:    args{cid: 10},
			want:    []models.Job{{JobTitle: "sde", Salary: "10,000"}, {JobTitle: "qa tester", Salary: "5,000"}},
			wantErr: false,
			mockRepoResponse: func() ([]models.Job, error) {
				return []models.Job{{JobTitle: "sde", Salary: "10,000"}, {JobTitle: "qa tester", Salary: "5,000"}}, nil

			},
		},
		{name: "failure if jobs are  not retrieved from cid",
			args:    args{cid: 10},
			want:    []models.Job{},
			wantErr: true,
			mockRepoResponse: func() ([]models.Job, error) {
				return []models.Job{}, errors.New("jobs not fetched")

			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mc := gomock.NewController(t)
			MockUserRepo := repository.NewMockUserRepo(mc)
			if tt.mockRepoResponse != nil {
				MockUserRepo.EXPECT().GetJobsFromCompany(gomock.Any()).Return(tt.mockRepoResponse()).AnyTimes()
			}
			s, _ := NewService(MockUserRepo, &auth.Auth{}, &caching.Redis{})
			got, err := s.ViewJobFromCompany(tt.args.cid)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.ViewJobFromCompany() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Service.ViewJobFromCompany() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestService_ViewAllJobs(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name string
		//s       *Service
		args             args
		want             []models.Job
		wantErr          bool
		mockRepoResponse func() ([]models.Job, error)
	}{
		{name: "success if all jobs are retrieved ",
			args:    args{context.Background()},
			want:    []models.Job{{JobTitle: "sde", Salary: "10,000"}, {JobTitle: "qa tester", Salary: "5,000"}},
			wantErr: false,
			mockRepoResponse: func() ([]models.Job, error) {
				return []models.Job{{JobTitle: "sde", Salary: "10,000"}, {JobTitle: "qa tester", Salary: "5,000"}}, nil

			},
		},
		{name: "failure if jobs are  not retrieved ",
			args:    args{context.Background()},
			want:    nil,
			wantErr: true,
			mockRepoResponse: func() ([]models.Job, error) {
				return []models.Job{}, errors.New("all jobs not fetched")

			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mc := gomock.NewController(t)
			MockUserRepo := repository.NewMockUserRepo(mc)
			if tt.mockRepoResponse != nil {
				MockUserRepo.EXPECT().GetAllJobs().Return(tt.mockRepoResponse()).AnyTimes()
			}
			s, _ := NewService(MockUserRepo, &auth.Auth{}, &caching.Redis{})
			got, err := s.ViewAllJobs(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.ViewAllJobs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Service.ViewAllJobs() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestService_ViewJobById(t *testing.T) {
	type args struct {
		ctx context.Context
		jid uint64
	}
	tests := []struct {
		name             string
		s                *Service
		args             args
		want             []models.Job
		wantErr          bool
		mockRepoResponse func() ([]models.Job, error)
	}{
		{name: "success if all jobs are fetched by job id ",
			args: args{
				ctx: context.Background(),
				jid: 10,
			},
			want:    []models.Job{{JobTitle: "sde", Salary: "10,000"}, {JobTitle: "qa tester", Salary: "5,000"}},
			wantErr: false,
			mockRepoResponse: func() ([]models.Job, error) {
				return []models.Job{{JobTitle: "sde", Salary: "10,000"}, {JobTitle: "qa tester", Salary: "5,000"}}, nil

			},
		},
		{name: "failure if jobs are  not retrieved by job id ",
			args: args{ctx: context.Background(),
				jid: 10,
			},
			want:    []models.Job{},
			wantErr: false,
			mockRepoResponse: func() ([]models.Job, error) {
				return []models.Job{}, errors.New(" jobs not fetched by id ")

			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mc := gomock.NewController(t)
			MockUserRepo := repository.NewMockUserRepo(mc)
			if tt.mockRepoResponse != nil {
				MockUserRepo.EXPECT().GetOneJob(gomock.Any()).Return(tt.mockRepoResponse()).AnyTimes()
			}
			s, _ := NewService(MockUserRepo, &auth.Auth{}, &caching.Redis{})
			got, err := s.ViewJobById(tt.args.ctx, tt.args.jid)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.ViewJobById() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Service.ViewJobById() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestService_AddJobDetails(t *testing.T) {
	type args struct {
		ctx context.Context
		cj  models.NewJobRequest
		cid uint64
	}
	tests := []struct {
		name             string
		args             args
		want             models.Response
		wantErr          bool
		mockRepoResponse func() (models.Response, error)
	}{
		{
			name: "success in adding jobs",
			want: models.Response{
				ID: 1,
			},
			args: args{
				ctx: context.Background(),
				cj: models.NewJobRequest{
					JobTitle:            "job",
					Salary:              "10000",
					MinimumNoticePeriod: int(10),
					MaximumNoticePeriod: uint64(20),
					Budget:              float64(10),
					JobDescription:      "This is a golang job",
					MinExperience:       float64(5),
					MaxExperience:       float64(10),
					LocationIDs:         []uint{uint(1), uint(2)},
					SkillIDs:            []uint{uint(1), uint(2)},
					WorkModeIDs:         []uint{uint(1), uint(2)},
					QualificationIDs:    []uint{uint(1), uint(2)},
					ShiftIDs:            []uint{uint(1), uint(2)},
					JobTypeIDs:          []uint{uint(1), uint(2)},
				},
				cid: 1,
			},
			wantErr: false,
			mockRepoResponse: func() (models.Response, error) {
				return models.Response{ID: 1}, nil
			},
		},
		{
			name: "error in adding jobs",
			want: models.Response{},
			args: args{
				ctx: context.Background(),
				cj:  models.NewJobRequest{},
				cid: 1,
			},
			wantErr: true,
			mockRepoResponse: func() (models.Response, error) {
				return models.Response{}, errors.New("job adding error")
			},
		},
		{
			name: "repo layer error",
			want: models.Response{
				ID: 1,
			},
			args: args{
				ctx: context.Background(),
				cj: models.NewJobRequest{
					JobTitle:            "job ",
					Salary:              "10000",
					MinimumNoticePeriod: int(10),
					MaximumNoticePeriod: uint64(20),
					Budget:              float64(10),
					JobDescription:      "This is a golang job",
					MinExperience:       float64(5),
					MaxExperience:       float64(10),
					LocationIDs:         []uint{uint(1), uint(2)},
					SkillIDs:            []uint{uint(1), uint(2)},
					WorkModeIDs:         []uint{uint(1), uint(2)},
					QualificationIDs:    []uint{uint(1), uint(2)},
					ShiftIDs:            []uint{uint(1), uint(2)},
					JobTypeIDs:          []uint{uint(1), uint(2)},
				},
				cid: 1,
			},
			wantErr: false,
			mockRepoResponse: func() (models.Response, error) {
				return models.Response{ID: 1}, nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mc := gomock.NewController(t)
			MockUserRepo := repository.NewMockUserRepo(mc)
			if tt.mockRepoResponse != nil {
				MockUserRepo.EXPECT().PostJob(gomock.Any()).Return(tt.mockRepoResponse()).AnyTimes()
			}
			s, _ := NewService(MockUserRepo, &auth.Auth{}, &caching.Redis{})

			got, err := s.AddJobDetails(tt.args.ctx, tt.args.cj, tt.args.cid)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.AddJobDetails() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Service.AddJobDetails() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestService_ProcessJobApplications(t *testing.T) {
	type args struct {
		applications []models.NewUserApplication
	}
	tests := []struct {
		name    string
		s       *Service
		args    args
		want    []models.NewUserApplication
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.s.ProcessJobApplications(tt.args.applications)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.ProcessJobApplications() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Service.ProcessJobApplications() = %v, want %v", got, tt.want)
			}
		})
	}
}
