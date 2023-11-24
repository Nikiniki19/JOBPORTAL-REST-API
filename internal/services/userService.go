package services

import (
	"context"
	"fmt"
	"job-portal-api/internal/models"
	"job-portal-api/internal/pkg"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func (s *Service) Signup(ctx context.Context, nu models.NewUser) (models.User, error) {
	hashedPass, err := pkg.PasswordHash(nu.Password)
	if err != nil {
		return models.User{}, err
	}
	userDetails := models.User{
		Name:         nu.Name,
		Email:        nu.Email,
		PasswordHash: hashedPass,
		Dob:          nu.Dob,
	}
	fmt.Printf("chck:: %#v", s)
	userDetails, err = s.UserRepo.CreateUser(ctx, userDetails)
	if err != nil {
		return models.User{}, err
	}
	return userDetails, nil
}

func (s *Service) Login(ctx context.Context, email, password string) (jwt.RegisteredClaims, error) {

	// We attempt to find the User record where the email
	// matches the provided email.
	var u models.User
	u, err := s.UserRepo.CheckEmail(ctx, email)
	if err != nil {
		return jwt.RegisteredClaims{}, err
	}
	// We check if the provided password matches the hashed password in the database.
	err = bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))
	if err != nil {
		return jwt.RegisteredClaims{}, err
	}

	// Successful authentication! Generate JWT claims.
	c := jwt.RegisteredClaims{
		Issuer:    "service project",
		Subject:   strconv.FormatUint(uint64(u.ID), 10),
		Audience:  jwt.ClaimStrings{"users"},
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}

	// And return those claims.
	return c, nil
}
func(s *Service)OTPGeneration(ctx context.Context,data models.ForgotPassword)(models.NewUser,error){
	check,err:=s.UserRepo.CheckEmail(ctx,data.Email )
	if err!=nil{
		
	}
}
 