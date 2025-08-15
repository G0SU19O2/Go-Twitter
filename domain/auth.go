package domain

import (
	"context"
	"errors"
	"fmt"

	gotwitter "github.com/G0SU19O2/Go-Twitter"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	UserRepo gotwitter.UserRepo
}

func NewAuthService(ur gotwitter.UserRepo) *AuthService {
	return &AuthService{
		UserRepo: ur,
	}
}

func (as *AuthService) Register(ctx context.Context, input gotwitter.RegisterInput) (gotwitter.AuthResponse, error) {
	input.Sanitize()
	if err := input.Validate(); err != nil {
		return gotwitter.AuthResponse{}, err
	}
	if _, err := as.UserRepo.GetByUsername(ctx, input.Username); !errors.Is(err, gotwitter.ErrNotFound) {
		return gotwitter.AuthResponse{}, gotwitter.ErrUsernameTaken
	}
	if _, err := as.UserRepo.GetByEmail(ctx, input.Email); !errors.Is(err, gotwitter.ErrNotFound) {
		return gotwitter.AuthResponse{}, gotwitter.ErrEmailTaken
	}
	user := gotwitter.User{
		Email:    input.Email,
		Username: input.Username,
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return gotwitter.AuthResponse{}, fmt.Errorf("failed to hash password: %w", err)
	}
	user.Password = string(hashedPassword)
	user, err = as.UserRepo.CreateUser(ctx, user)
	if err != nil {
		return gotwitter.AuthResponse{}, fmt.Errorf("failed to create user: %w", err)
	}
	return gotwitter.AuthResponse{AccessToken: "some-jwt-token", User: user}, nil
}
