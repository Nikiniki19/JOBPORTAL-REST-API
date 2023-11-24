package middlewares

import (
	"fmt"
	"job-portal-api/internal/auth"
)

// Mid struct
type Mid struct {
	a auth.Authentication
}

// func new mid
func NewMid(a auth.Authentication) (Mid, error) {
	if a == nil {
		return Mid{}, fmt.Errorf("auth cannot be nil")
	}
	return Mid{a: a}, nil
}
