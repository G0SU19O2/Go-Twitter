package gotwitter

import (
	"context"
	"fmt"
	"regexp"
	"strings"
)

var (
	UsernameMinLength = 3
	PasswordMinLength = 6
	EmailRegexp       = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
)

type AuthService interface {
	Register(ctx context.Context, input RegisterInput) (AuthResponse, error)
}
type AuthResponse struct {
	AccessToken string
	User        User
}
type RegisterInput struct {
	Email           string
	Username        string
	Password        string
	ConfirmPassword string
}

func (input *RegisterInput) Sanitize() {
	input.Email = strings.TrimSpace(input.Email)
	input.Email = strings.ToLower(input.Email)
	input.Username = strings.TrimSpace(input.Username)
}

func (input *RegisterInput) Validate() error {
	if !EmailRegexp.MatchString(input.Email) {
		return fmt.Errorf("%w: invalid email format", ErrValidation)
	}
	if len(input.Username) < UsernameMinLength {
		return fmt.Errorf("%w: username must be at least %d characters long", ErrValidation, UsernameMinLength)
	}
	if len(input.Password) < PasswordMinLength {
		return fmt.Errorf("%w: password must be at least %d characters long", ErrValidation, PasswordMinLength)
	}
	if input.Password != input.ConfirmPassword {
		return fmt.Errorf("%w: passwords do not match", ErrValidation)
	}
	return nil
}

type LoginInput struct {
	Email    string
	Password string
}

func (input *LoginInput) Sanitize() {
	input.Email = strings.TrimSpace(input.Email)
	input.Email = strings.ToLower(input.Email)
}

func (input *LoginInput) Validate() error {
	if !EmailRegexp.MatchString(input.Email) {
		return fmt.Errorf("%w: invalid email format", ErrValidation)
	}
	if len(input.Password) < 1 {
		return fmt.Errorf("%w: password is required", ErrValidation)
	}
	return nil
}