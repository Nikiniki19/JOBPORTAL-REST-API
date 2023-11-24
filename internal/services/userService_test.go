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
	"time"

	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/mock/gomock"
)

func TestService_UserSignup(t *testing.T) {
	type args struct {
		ctx context.Context
		nu  models.NewUser
	}
	tests := []struct {
		name string
		//	s                *Service
		args             args
		want             models.User
		wantErr          bool
		mockRepoResponse func() (models.User, error)
	}{
		{name: "error in converting password to hashpassword during signup",
			args: args{nu: models.NewUser{
				Name:     "Nikitha",
				Email:    "niki123@gmail.com",
				Password: ""},
				ctx: context.Background()},
			want:    models.User{},
			wantErr: true,
			mockRepoResponse: func() (models.User, error) {
				return models.User{}, errors.New("error in password hash")
			},
		},
		{name: "success if user signed up",
			args: args{ctx: context.Background(),
				nu: models.NewUser{
					Name:     "Nikitha",
					Email:    "niki123@gmail.com",
					Password: "1234"},
			},
			want:    models.User{Name: "Nikitha", Email: "niki123@gmail.com"},
			wantErr: false,
			mockRepoResponse: func() (models.User, error) {
				return models.User{Name: "Nikitha", Email: "niki123@gmail.com"}, nil
			},
		},
		{name: "failure if user not signed up",
			args: args{nu: models.NewUser{
				Name:     "Nikitha",
				Email:    "niki123@gmail.com",
				Password: "1234"},
				ctx: context.Background()},
			want:    models.User{},
			wantErr: true,
			mockRepoResponse: func() (models.User, error) {
				return models.User{}, errors.New("user signup failed")
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mc := gomock.NewController(t)
			MockUserRepo := repository.NewMockUserRepo(mc)
			if tt.mockRepoResponse != nil {
				MockUserRepo.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Return(tt.mockRepoResponse()).AnyTimes()
			}
			s, _ := NewService(MockUserRepo, &auth.Auth{}, &caching.Redis{})
			got, err := s.Signup(tt.args.ctx, tt.args.nu)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.UserSignup() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Service.UserSignup() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestService_Login(t *testing.T) {
	type args struct {
		ctx      context.Context
		email    string
		password string
	}
	tests := []struct {
		name             string
		s                *Service
		args             args
		want             jwt.RegisteredClaims
		wantErr          bool
		mockRepoResponse func() (models.User, error)
	}{
		{name: "failure case for login",
			args: args{email: "niki123@gmail.com",
				password: ""},
			want:    jwt.RegisteredClaims{},
			wantErr: true,
			mockRepoResponse: func() (models.User, error) {
				return models.User{Email: "niki1232gmail.com", PasswordHash: "$2a$10$vtON7w6i6G.OZT3zKpR00elHrB7P8e3IknFgOfhvfXXHFIk6ytDQC"}, nil
			},
		},
		{name: "success case for login",
			args:    args{email: "niki1232gmail.com", password: "abcdefg"},
			want:    jwt.RegisteredClaims{Issuer: "service project", Subject: "0", Audience: jwt.ClaimStrings{"users"}, ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)), IssuedAt: jwt.NewNumericDate(time.Now())},
			wantErr: false,
			mockRepoResponse: func() (models.User, error) {
				return models.User{Email: "niki1232gmail.com", PasswordHash: "$2a$10$vtON7w6i6G.OZT3zKpR00elHrB7P8e3IknFgOfhvfXXHFIk6ytDQC"}, nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mc := gomock.NewController(t)
			MockUserRepo := repository.NewMockUserRepo(mc)
			if tt.mockRepoResponse != nil {
				MockUserRepo.EXPECT().CheckEmail(gomock.Any(), gomock.Any()).Return(tt.mockRepoResponse()).AnyTimes()
			}
			s, _ := NewService(MockUserRepo, &auth.Auth{}, &caching.Redis{})
			got, err := s.Login(tt.args.ctx, tt.args.email, tt.args.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.Login() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Service.Login() = %v, want %v", got, tt.want)
			}
		})
	}
}
