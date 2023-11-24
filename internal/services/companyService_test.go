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

func TestService_AddCompanyDetails(t *testing.T) {
	type args struct {
		ctx         context.Context
		companyData models.Company
	}
	tests := []struct {
		name string
		//	s       *Service
		args             args
		want             models.Company
		wantErr          bool
		mockRepoResponse func() (models.Company, error)
	}{
		{name: "success if companies are added",
			args: args{
				ctx:         context.Background(),
				companyData: models.Company{},
			},
			want: models.Company{
				CompanyName: "infosys",
				Address:     "bangalore",
				Domain:      "software",
			},
			wantErr: false,
			mockRepoResponse: func() (models.Company, error) {
				return models.Company{
					CompanyName: "infosys",
					Address:     "bangalore",
					Domain:      "software",
				}, nil
			},
		},

		{name: "failure if there are no companies added",
			args: args{
				ctx:         context.Background(),
				companyData: models.Company{},
			},
			want:    models.Company{},
			wantErr: true,
			mockRepoResponse: func() (models.Company, error) {
				return models.Company{}, errors.New("please provide the fields")
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mc := gomock.NewController(t)
			MockUserRepo := repository.NewMockUserRepo(mc)
			if tt.mockRepoResponse != nil {
				MockUserRepo.EXPECT().CreateCom(gomock.Any()).Return(tt.mockRepoResponse()).AnyTimes()
			}
			s, _ := NewService(MockUserRepo, &auth.Auth{},&caching.Redis{})
			got, err := s.AddCompanyDetails(tt.args.ctx, tt.args.companyData)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.AddCompanyDetails() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Service.AddCompanyDetails() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestService_ViewAllCompanies(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name string
		//s       *Service
		args             args
		want             []models.Company
		wantErr          bool
		mockRepoResponse func() ([]models.Company, error)
	}{
		{name: "companies are viewed",
			args: args{
				ctx: context.Background(),
			},
			want:    []models.Company{{CompanyName: "Tek", Address: "bangalore", Domain: "software"}, {CompanyName: "Infy", Address: "bangalore", Domain: "software"}},
			wantErr: false,
			mockRepoResponse: func() ([]models.Company, error) {
				return []models.Company{{CompanyName: "Tek", Address: "bangalore", Domain: "software"}, {CompanyName: "Infy", Address: "bangalore", Domain: "software"}}, nil
			},
		},
		{name: "companies are  not viewed",
			args: args{
				ctx: context.Background(),
			},
			want:    nil,
			wantErr: true,
			mockRepoResponse: func() ([]models.Company, error) {
				return nil, errors.New("No companies are present")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mc := gomock.NewController(t)
			MockUserRepo := repository.NewMockUserRepo(mc)
			if tt.mockRepoResponse != nil {
				MockUserRepo.EXPECT().GetAllTheCompanies().Return(tt.mockRepoResponse()).AnyTimes()
			}
			s, _ := NewService(MockUserRepo, &auth.Auth{},&caching.Redis{})
			got, err := s.ViewAllCompanies(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.ViewAllCompanies() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Service.ViewAllCompanies() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestService_ViewCompanyDetails(t *testing.T) {
	type args struct {
		ctx context.Context
		id  uint64
	}
	tests := []struct {
		name    string
		//s       *Service
		args    args
		want    models.Company
		wantErr bool
		mockRepoResponse func() (models.Company, error)

	}{
		{name: "success if companies are fetched by id",
		args: args{
			ctx:         context.Background(),
			id:1,
		},
		want: models.Company{
			CompanyName: "infosys",
			Address:     "bangalore",
			Domain:      "software",
		},
		wantErr: false,
		mockRepoResponse: func() (models.Company, error) {
			return models.Company{
				CompanyName: "infosys",
				Address:     "bangalore",
				Domain:      "software"}, nil
		},
	
       },

	   {name: "failure if companies are fetched by id",
	   args: args{
		   ctx:         context.Background(),
		  
	   },
	   want: models.Company{},
	   wantErr: true,
	   mockRepoResponse: func() (models.Company, error) {
		   return models.Company{}, errors.New("id not present")
	   },
	},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mc := gomock.NewController(t)
			MockUserRepo := repository.NewMockUserRepo(mc)
			if tt.mockRepoResponse != nil {
				MockUserRepo.EXPECT().GetCompany(gomock.Any()).Return(tt.mockRepoResponse()).AnyTimes()
			}
			s, _ := NewService(MockUserRepo, &auth.Auth{},&caching.Redis{})
			got, err := s.ViewCompanyDetails( tt.args.ctx,tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.ViewCompanyDetails() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Service.ViewCompanyDetails() = %v, want %v", got, tt.want)
			}
		})
	}
}
