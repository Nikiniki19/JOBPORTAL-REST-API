package services

import (
	"context"
	"errors"
	"fmt"
	"job-portal-api/internal/models"
	"job-portal-api/internal/pkg"

	"math/rand"
	"net/smtp"

	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)
var otp string
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

func (s *Service) OTPGeneration(ctx context.Context, data models.ForgotPassword) (bool, string, error) {
	check, err := s.UserRepo.CheckEmail(ctx, data.Email)
	if err != nil {
		return false, "", errors.New("email not found in db")
	}

	// Sender's email address and password
	from := "niki16809@gmail.com"
	password := "xagg yzsu kkvd jbrr"

	// Recipient's email address
	to := check

	// SMTP server details
	smtpServer := "smtp.gmail.com"
	smtpPort := 587
	otp = generateOTP(4)
	// Message content
	message := fmt.Sprintf("Subject: Test Email\n\nThis is a test email body.", otp)
	//adding otp to cache
	// err = s.rdb.AddEmailToCache(ctx, check.Email, otp)
	// if err != nil {
	// 	return false, "", err
	// }
	// Authentication information
	auth := smtp.PlainAuth("", from, password, smtpServer)

	// SMTP connection
	smtpAddr := fmt.Sprintf("%s:%d", smtpServer, smtpPort)
	err = smtp.SendMail(smtpAddr, auth, from, []string{to.Email}, []byte(message))
	if err != nil {
		//fmt.Println("Error sending email:", err)
		return false, "", errors.New("error sending email:")
	}

	fmt.Println("Email sent successfully!")
	return true, otp, nil

}

func generateOTP(length int) string {
	// Seed the random number generator with the current timestamp
	rand.Seed(time.Now().UnixNano())

	// Define the characters allowed in the OTP
	otpChars := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

	// Generate the OTP
	otp := make([]byte, length)
	for i := range otp {
		otp[i] = otpChars[rand.Intn(len(otpChars))]
	}

	return string(otp)
}

func (s *Service) ChangePassword(ctx context.Context, cj models.OtpPassword) (string, error) {
	val, _ := s.rdb.GetEmailFromCache(ctx, cj.Otp)
	// if err != nil {
	// 	return "", err
	// }

	if val != otp {
		if cj.Password == cj.ConfirmPassword {
			newuserotp, err := s.UserRepo.CheckEmail(ctx, cj.Email)
			if err != nil {
				return "", errors.New("email not matching")
			}
			hashedPass, err := pkg.PasswordHash(cj.ConfirmPassword)
			if err != nil {
				return "", errors.New("error in pwd hash")
			}
			details := models.User{
				Name:         newuserotp.Name,
				PasswordHash: hashedPass,
			}
			_, err = s.UserRepo.UpdatePwdInDb(ctx, cj.Email, details)
			if err != nil {
				return "", errors.New("password not matching")
			}

		}
	}
	return "pwd set", nil
}
