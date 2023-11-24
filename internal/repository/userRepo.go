package repository

import (
	"context"
	"errors"
	"job-portal-api/internal/models"

	"github.com/rs/zerolog/log"
)

func (r *Repo) CreateUser(ctx context.Context, nu models.User) (models.User, error) {

	err := r.DB.Create(&nu).Error
	if err != nil {
		log.Info().Err(err).Send()
		return models.User{}, errors.New("could not create user")
	}

	return nu, nil
}
func (r *Repo) CheckEmail(ctx context.Context, email string) (models.User, error) {
	var userDetails models.User
	result := r.DB.Where("email = ?", email).First(&userDetails)
	if result.Error != nil {
		log.Info().Err(result.Error).Send()
		return models.User{}, errors.New("email not found")
	}
	return userDetails, nil

}
func(r * Repo)ValidateEmail(ctx context.Context,email string)(models.ForgotPassword,error){
	var validemail models.ForgotPassword
	res:=r.DB.Where("email=?",email).First(&validemail)
	if res.Error!=nil{
		log.Info().Err(res.Error).Send()
		return models.ForgotPassword{},errors.New("email not matched")
	}
	return validemail,nil
}
