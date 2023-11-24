package pkg

import (
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func PasswordHash(password string) (string, error) {

	if password == "" {
		return "", errors.New("invalid password")
	}

	hashedpass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("error in hashing the pwd:%w", err)
	}
	return string(hashedpass), nil
}
